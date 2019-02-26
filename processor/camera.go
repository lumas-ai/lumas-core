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
  "google.golang.org/grpc"
  . "github.com/3d0c/gmf"
  _struct "github.com/golang/protobuf/ptypes/struct"
  api "github.com/lumas-ai/lumas-core/protos/golang/provider"
  _ "image/jpeg"
)

var cameras []*Camera

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

func FindCamera(id int64) (*Camera) {
  for _, cam := range cameras {
    if cam.Id == id {
      return cam
    }
  }

  return nil
}

func NewCamera() (*Camera, error) {
  c := &Camera{}
  cameras = append(cameras, c)
  return c, nil
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
  //defer conn.Close()
  client := api.NewCameraClient(conn)

  return client, nil
}

func findFreePort() (int32, error) {
  //Find an open UDP port
  for port := 9000; port < 10000; port++ {
    addr := net.UDPAddr{
      Port: port,
      IP:   net.ParseIP("172.17.0.2"),
    }
    lp , err := net.ListenUDP("udp", &addr)
    defer lp.Close()
    if err != nil {
      continue
    }

    return int32(port), nil
  }

  return 0, errors.New("Could not find available udp port between 9000 and 10000")
}

func (s *Camera) NewRTPClient(sdp string) (*FmtCtx, error) {
  var rtpOptions []*Option

  //gmf.OpenInput can only take a file name for a SDP file, so
  //we have to save one to disk
  filename := fmt.Sprintf("/tmp/%d.sdp", s.Id)
  rtpOptions = append([]*Option{ {Key: "protocol_whitelist", Val: "udp,file,rtp,crypto"} })

  inputCtx := NewCtx()
  inputCtx.SetOptions(rtpOptions)

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

  //Look for motion in every 5th frame to conserve resources
  go func() {
    i := 1
    for frame := range frames {
      if i == 5 {
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
    }

    if packet.StreamIndex() != videoStream.Index() {
      //It's an audio packet
      continue
    }

    ist := assert(inputCtx.GetStream(packet.StreamIndex())).(*Stream)

    frame, err := packet.Frames(ist.CodecCtx())
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

  audioPort, err := findFreePort()
  if err != nil {
    return errors.New("Could not allocate a free UDP port for audio RTP stream")
  }

  videoPort, err := findFreePort()
  if err != nil {
    return errors.New("Could not allocate a UDP port for video RTP stream")
  }

  rtpConfig := api.RTPConfig{ RtpAddress: "192.168.2.207",
    AudioRTPPort: audioPort,
    VideoRTPPort: videoPort,
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
  if err == io.EOF {
    fmt.Println("Stream ended")
    return nil
  }

  go func() {
    fmt.Println(fmt.Sprintf("Video SDP: %s", status.Sdp.Video))
    s.processVideo(status.Sdp.Video)
  }()

  go func() {
    fmt.Println(fmt.Sprintf("Audio SDP: %s", status.Sdp.Audio))
    //s.processAudio(status.Sdp.Audio)
  }()

  //Listen for the status updates
  for {
    status, err := stream.Recv()
    if err == io.EOF {
      fmt.Println("Stream ended")
      break
    }

    s.Log(fmt.Sprintf("Sent Frames: %d", status.SentFrames))
    s.Log(fmt.Sprintf("Dropped Frames: %d", status.DroppedFrames))
    time.Sleep(5 * time.Second)
  }

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
