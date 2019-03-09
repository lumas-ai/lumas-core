package processor

import (
  "os"
  "bufio"
  "time"
  "net"
  "io"
  "fmt"
  "log"
  "strings"
  "encoding/base64"
  "errors"
  "image"
  "context"
  "crypto/md5"
  "google.golang.org/grpc"
  . "github.com/3d0c/gmf"
  _struct "github.com/golang/protobuf/ptypes/struct"
  api "github.com/lumas-ai/lumas-core/protos/golang/provider"
  _ "image/jpeg"
)

var allocatedRTPPorts []int

type RTPClient struct {
  sdp string
  port int32
}

type Camera struct {
  Id int64
  Name string
  audioRTPClient *RTPClient
  videoRTPClient *RTPClient
  provider string
  providerAddress string
  providerConfig *_struct.Struct //Arbitrary struct defined in the protobuf
  inputCtx *FmtCtx
}

func NewCamera() (*Camera, error) {
  c := &Camera{}
  return c, nil
}

func (r *RTPClient) Close() error {
  err := unallocatePort(int(r.port))
  if err != nil {
    return err
  }

  return nil
}

func (s *Camera) providerClient() (api.CameraClient, error) {
  var opts []grpc.DialOption

  if len(s.providerAddress) == 0 {
    return nil, errors.New("Camera does not have a configured provider")
  }

  opts = append(opts, grpc.WithInsecure())
  conn, err := grpc.Dial(s.providerAddress, opts...)
  if err != nil {
    return nil, err
  }
  client := api.NewCameraClient(conn)

  return client, nil
}

func getLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

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

  fmt.Println(fmt.Sprintf("Allocating Ports: %v", rslice))
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

func findFreePorts(num int) ([]int32, error) {
  ports := make([]int32, num)
  ip := getLocalIP()

  log.Println(fmt.Sprintf("Creating RTP clients on IP %s", ip))
  //Find an open UDP port
  for i := 0; i < num; i++ {
    for port := 9000; port < 10000; port++ {
      addr := net.UDPAddr{
        Port: port,
        IP:   net.ParseIP(ip),
      }
      lp , err := net.ListenUDP("udp", &addr)
      if err != nil {
        continue
      }
      defer lp.Close()

      ports[i] = int32(port)
      break
    }
  }

  return ports, nil
}

func (s *Camera) NewRTPClient(sdp string) (*FmtCtx, error) {
  var rtpOptions []*Option

  //gmf.OpenInput can only take a file name for a SDP file, so
  //we have to save one to disk
  h := md5.New()
  io.WriteString(h, sdp)
  filename := fmt.Sprintf("/tmp/%x.sdp", h.Sum(nil))

  f, err := os.Create(filename)
  defer f.Close()
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
    fmt.Println("Could not remove %s from filesystem", filename)
  }

  return inputCtx, nil
}

func (s *Camera) GetSnapshot() (image.Image, error) {
  client, err := s.providerClient()
  if err != nil {
    return nil, err
  }

  c := api.CameraConfig{
    Config: s.providerConfig,
  }

  img, err := client.Snapshot(context.Background(), &c)
  if err != nil {
    return nil, err
  }

  reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(img.Base64Image))
  decodedImg, _, err := image.Decode(reader)

  return decodedImg, nil
}

func (s *Camera) Log(msg string) {
  fmt.Println(fmt.Sprintf("[%d] %s", s.Id, msg))
}

func (s *Camera) processVideo(sdp string) error {
  inputCtx, err := s.NewRTPClient(sdp)
  if err != nil {
    return err
  }

  videoStream, err := inputCtx.GetStream(0)
  if err != nil {
    s.Log(fmt.Sprintf("Unable to open RTP stream: %s", err.Error()))
  }

  frames  := make(chan *Frame, 100)
  defer close(frames)
  doneFrames := make(chan *Frame, 100)
  defer close(doneFrames)
  motionFrames  := make(chan *Frame, 100)
  defer close(motionFrames)
  motions := make(chan *Motion, 100)
  defer close(motions)
  packets := make(chan *Packet, 100)
  defer close(packets)
  donePackets := make(chan *Packet, 100)
  defer close(donePackets)

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
  go func() {
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
  go s.WriteFile(packets, donePackets, inputCtx)

  go DetectMotion(motionFrames, doneFrames, motions, videoStream.CodecCtx(), videoStream.TimeBase().AVR())

  go func() {
    for motion := range motions {
      if motion.MotionDetected {
        fmt.Println("found motion in frame ")
      }
    }
  }()

  for {
    packet, err := inputCtx.GetNextPacket()
    if err == io.EOF {
      return err
    } else if err != nil {
      fmt.Println("Could not get packet. Skipping")
      continue
    }

    frame, err := packet.Frames(videoStream.CodecCtx())
    if err != nil {
      log.Println("Missed packet at " + string(packet.Pts()) + ". " + err.Error())
      continue
    }

    packets <- packet
    frames <- frame
  }
}

func (s *Camera) processAudio(sdp string) error {
  //Allocate RTP listeners for audio/video streams
  _, err := s.NewRTPClient(sdp)
  if err != nil {
    return err
  }

  return nil
}

func (s *Camera) StopFeed() error {
  provider, err := s.providerClient()
  if err != nil {
    s.Log(err.Error())
    return err
  }

  //Close the RTP clients
  if s.videoRTPClient != nil {
    s.videoRTPClient.Close()
  }
  if s.audioRTPClient != nil {
    s.audioRTPClient.Close()
  }

  c := api.CameraConfig{
    Config: s.providerConfig,
  }

  rtpConfig := api.RTPConfig{
    CameraConfig: &c,
  }

  provider.StopRTPStream(context.Background(), &rtpConfig)

  return nil
}

func (s *Camera) ProcessFeed() error {
  provider, err := s.providerClient()
  if err != nil {
    s.Log(err.Error())
    return err
  }

  c := api.CameraConfig{
    Config: s.providerConfig,
  }

  ports, err := allocatePorts(2)
  if err != nil {
    return errors.New("Could not allocate free UDP ports for RTP streams")
  }

  audioPort := ports[0]
  videoPort := ports[1]

  rtpConfig := api.RTPConfig{ RtpAddress: "192.168.2.207",
    AudioRTPPort: int32(audioPort),
    VideoRTPPort: int32(videoPort),
    CameraConfig: &c,
  }

  //Tell the camera provider to start streaming
  stream, err := provider.StreamRTP(context.Background(), &rtpConfig)
  if err != nil {
    s.Log(fmt.Sprintf("Provider could not stream: %s", err.Error()))
    return err
  }
  //The first response contains the SDP information
  status, err := stream.Recv()
  if err == io.EOF || status == nil {
    s.Log("Did not receive SDP information from provider")
    return err
  }

  if status.Sdp.Video != "" {
    go func() {
      s.processVideo(status.Sdp.Video)
    }()
  } else {
    s.Log("Provider did not provide any video streaming information. Skipping processing video")
  }

  if status.Sdp.Audio != "" {
    go func() {
      //s.processAudio(status.Sdp.Audio)
    }()
  } else {
    s.Log("Provider did not provide any audio streaming information. Skipping processing audio")
  }

  //Listen for the status updates
  go func() {
    for {
      if s.audioRTPClient != nil {
        defer s.audioRTPClient.Close()
      }

      if s. videoRTPClient != nil {
        defer s.videoRTPClient.Close()
      }

      status, err := stream.Recv()
      if err == io.EOF {
        s.Log("Stream ended")
        break
      }

      s.Log(fmt.Sprintf("Sent Frames: %d", status.SentFrames))
      s.Log(fmt.Sprintf("Dropped Frames: %d", status.DroppedFrames))
      time.Sleep(5 * time.Second)
    }
  }()

  return nil
}

func (s *Camera) SetId(id int64) (*Camera) {
  s.Id = id
  return s
}

func (s *Camera) SetName(name string) (*Camera) {
  s.Name = name
  return s
}

func (s *Camera) SetProvider(provider string) (*Camera) {
  s.provider = provider
  return s
}

func (s *Camera) SetProviderAddress(address string) (*Camera) {
  s.providerAddress = address
  return s
}

func (s *Camera) SetProviderConfig(config *_struct.Struct) (*Camera) {
  s.providerConfig = config
  return s
}
