package processor

import (
  "os"
  "bufio"
  "io"
  "fmt"
  "log"
  "errors"
  "sync"
  "crypto/md5"
  . "github.com/3d0c/gmf"
  api "github.com/lumas-ai/lumas-core/protos/golang/processor"
  _ "image/jpeg"
)

var allocatedRTPPorts []int

type RTPClient struct {
  sdp string
  port int32
  open bool
  inputCtx *FmtCtx
}

type Feed struct {
  Id string
  RTPConfig *api.RTPConfig
  audioRTPClient *RTPClient
  audioCloseChan chan bool
  videoRTPClient *RTPClient
  videoCloseChan chan bool
}

//func NewCamera() (*Camera, error) {
//  c := &Camera{}
//
//  return c, nil
//}

func (r *RTPClient) isOpen() bool {
  return r.open
}

func (r *RTPClient) Close() error {
  r.open = false

  //Give the ports back for new RTP clients
  err := unallocatePort(int(r.port))
  if err != nil {
    return err
  }

  return nil
}

//func (s *Camera) providerClient() (api.CameraClient, error) {
//  var opts []grpc.DialOption
//
//  if len(s.providerAddress) == 0 {
//    return nil, errors.New("Camera does not have a configured provider")
//  }
//
//  opts = append(opts, grpc.WithInsecure())
//  conn, err := grpc.Dial(s.providerAddress, opts...)
//  if err != nil {
//    return nil, err
//  }
//  client := api.NewCameraClient(conn)
//
//  return client, nil
//}

func allocatePorts(num int) ([]int, error) {
  var rslice []int
  port := 9000

  for i := 0; i < num; i++ {
    for j := 0; j < len(allocatedRTPPorts); j++ {
      if allocatedRTPPorts[j] == port {
        port++
        continue
      }
    }

    allocatedRTPPorts = append(allocatedRTPPorts, port)
    rslice = append(rslice, port)
    port++
  }

  return rslice, nil
}

func unallocatePort(port int) error {
  pslice := allocatedRTPPorts

  for i := 0; i < len(pslice); i++ {
    if pslice[i] == port {
      pslice[len(pslice) - 1], pslice[port] = pslice[port], pslice[len(pslice) -1]
      pslice = pslice[:len(pslice)-1]
      break
    }
  }

  return nil
}

func (s *Feed) NewRTPClient(sdp string) (*RTPClient, error) {
  var rtpOptions []*Option

  //gmf.OpenInput can only take a file name for a SDP file, so
  //we have to save one to disk
  h := md5.New()
  io.WriteString(h, sdp)
  filename := fmt.Sprintf("/tmp/%x.sdp", h.Sum(nil))

  f, err := os.Create(filename)
  defer f.Close()
  defer os.Remove(filename)
  if err != nil {
    return nil, err
  }

  w := bufio.NewWriter(f)
  _, err = w.WriteString(sdp)
  if err != nil {
    return nil, err
  }
  w.Flush()

  //Now we're ready to create the RTP client
  rtpOptions = append([]*Option{ {Key: "protocol_whitelist", Val: "udp,file,rtp_mpegts,rtp,crypto"} })

  inputCtx := NewCtx()
  inputCtx.SetOptions(rtpOptions)

  inputCtx.OpenInput(filename)
  inputCtx.Dump()

  err = os.Remove(filename)
  if err != nil {
    //We shouldn't fail on this error. Just log it
    s.Log(fmt.Sprintf("Could not remove %s from filesystem", filename))
  }

  rtpClient := &RTPClient{
    sdp: sdp,
    inputCtx: inputCtx,
    open: true,
  }
  return rtpClient, nil
}

//func (s *Camera) GetSnapshot() (image.Image, error) {
//  client, err := s.providerClient()
//  if err != nil {
//    return nil, err
//  }
//
//  c := api.CameraConfig{
//    Config: s.providerConfig,
//  }
//
//  img, err := client.Snapshot(context.Background(), &c)
//  if err != nil {
//    return nil, err
//  }
//
//  reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(img.Base64Image))
//  decodedImg, _, err := image.Decode(reader)
//
//  return decodedImg, nil
//}

func (s *Feed) Log(msg string) {
  log.Println(fmt.Sprintf("[%d] %s", s.Id, msg))
}

func FmtError(err *api.Result) string {
  if err == nil || err.ErrorKind == "" || err.Message == "" {
    return "Unknown Error"
  }

  return fmt.Sprintf("ErrorKind: %s - %s", err.ErrorKind, err.Message)
}

func (s *Feed) processVideo(sdp string) error {
  var wg sync.WaitGroup
  var motionWaitGroup sync.WaitGroup
  var motionResultsWaitGroup sync.WaitGroup
  var err error

  s.videoRTPClient, err = s.NewRTPClient(sdp)
  if err != nil {
    return err
  }

  inputCtx := s.videoRTPClient.inputCtx
  defer inputCtx.Free()

  videoStream, err := inputCtx.GetStream(0)
  if err != nil {
    s.Log(fmt.Sprintf("Unable to open RTP stream: %s", err.Error()))
  }

  frames       := make(chan *Frame, 100)
  doneFrames   := make(chan *Frame, 100)
  motionFrames := make(chan *Frame, 100)
  motions      := make(chan *Motion, 100)
  packets      := make(chan *Packet, 100)
  donePackets  := make(chan *Packet, 100)

  // When frames are done being processed, they're sent to the
  // to the doneFrames channel so they can be freed from memory
  go func() {
    for frame := range doneFrames {
      frame.Free()
    }
  }()

  //Free up the packets from memory when they're done being written to disk
  go func() {
    for pkt := range donePackets {
      pkt.Free()
    }
  }()

  //Look for motion in every 10th frame to conserve resources
  wg.Add(1)
  go func() {
    defer wg.Done()

    i := 1
    for frame := range frames {
      if i == 10 {
        i = 1
        motionFrames <- frame
      } else {
        doneFrames <- frame
        i++
      }
    }
  }()

  //Write the packets to disk concurrently
  wg.Add(1)
  go func() {
    defer wg.Done()
    s.WriteFile(packets, donePackets, inputCtx)
  }()

  motionWaitGroup.Add(1)
  go func() {
    defer motionWaitGroup.Done()
    DetectMotion(motionFrames, doneFrames, motions, videoStream.CodecCtx(), videoStream.TimeBase().AVR())
  }()

  motionResultsWaitGroup.Add(1)
  go func() {
    defer motionResultsWaitGroup.Done()
    for motion := range motions {
      if motion.MotionDetected {
        fmt.Println("found motion in frame ")
      }
    }
  }()

  getPackets:
  for {
    select {
    case _ = <-s.videoCloseChan:
      break getPackets
    default:
    }

    packet, err := inputCtx.GetNextPacket()
    if err != nil {
      s.Log("Could not get packet. Skipping")
      continue
    }

    frame, err := packet.Frames(videoStream.CodecCtx())
    if err != nil {
      s.Log(fmt.Sprintf("Missed packet at %d: %s", packet.Pts(), err.Error()))
      continue
    }

    packets <- packet
    frames <- frame
  }

  //Close up our channels and wait for the goroutines to finish
  //Note that we can close up the channels without throwing away
  //any packets or frames that are currently waiting to be processed
  close(frames)
  close(packets)

  //Wait for the frames, packets, and motion goroutines to finish
  wg.Wait()

  //Give the motion goroutine a chance to finish reading
  //from the now closed frames channel
  close(motionFrames)
  motionWaitGroup.Wait()

  close(motions)
  motionResultsWaitGroup.Wait()

  //Close these AFTER the other goroutines have finished
  close(doneFrames)
  close(donePackets)

  //Give the final all clear
  s.videoCloseChan <-true

  s.Log("Video Processing has been closed")
  return nil
}

func (s *Feed) processAudio(sdp string) error {
  var err error

  //Allocate RTP listeners for audio/video streams
  s.audioRTPClient, err = s.NewRTPClient(sdp)
  if err != nil {
    return err
  }

  return nil
}

func (s *Feed) StopFeed() error {
  //Send the close signal to the processors and wait for them to finish
  if s.audioRTPClient != nil && s.audioRTPClient.isOpen() {
    s.audioCloseChan <- true
  }
  if s.videoRTPClient != nil && s.videoRTPClient.isOpen() {
    s.videoCloseChan <- true
  }
  if s.audioRTPClient != nil && s.audioRTPClient.isOpen() {
    _ = <-s.audioCloseChan
  }
  if s.videoRTPClient != nil && s.videoRTPClient.isOpen() {
    _ = <-s.videoCloseChan
  }

  //Close the RTP clients if they're still open
  if s.videoRTPClient != nil && s.videoRTPClient.isOpen() {
    s.videoRTPClient.Close()
  }
  if s.audioRTPClient != nil && s.audioRTPClient.isOpen() {
    s.audioRTPClient.Close()
  }

  return nil
}

func NewFeed(id string, videoSDP string, audioSDP string) (*Feed, error) {
  f := &Feed{ Id: id}

  //Allocate 2 ports - one for video RTP and one for audio RTP
  ports, err := allocatePorts(2)
  if err != nil {
    return nil, errors.New("Could not allocate free UDP ports for RTP streams")
  }

  audioPort := ports[0]
  videoPort := ports[1]

  //TODO: The RtpAddress returned needs to be dynamically determined by 
  //the address of each particular instance of the service
  rtpConfig := api.RTPConfig{ RtpAddress: "processor",
    AudioRTPPort: int32(audioPort),
    VideoRTPPort: int32(videoPort),
  }

  f.RTPConfig = &rtpConfig

  //Create a channel to pass messages to the video and audio
  //goroutines that they need to stop
  if videoSDP != "" {
    go func() {
      f.videoCloseChan = make(chan bool)
      f.processVideo(videoSDP)
    }()
  } else {
    f.Log("No video streaming information was provided. Skipping processing video")
  }

  if audioSDP != "" {
    go func() {
      f.audioCloseChan = make(chan bool)
      //f.processAudio(audioSDP, closeChan) //TODO: Processing audio breaks things for some reason
    }()
  } else {
    f.Log("Provider did not provide any audio streaming information. Skipping processing audio")
  }

  return f, nil
}
