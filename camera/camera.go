package camera

import (
  "fmt"
  "log"
  "strings"
  "encoding/base64"
  "errors"
  "context"
  "image"
  "google.golang.org/grpc"
  _struct "github.com/golang/protobuf/ptypes/struct"
  providerAPI "github.com/lumas-ai/lumas-core/protos/golang/provider"
  cameraAPI "github.com/lumas-ai/lumas-core/protos/golang/camera"
  streamAPI "github.com/lumas-ai/lumas-core/protos/golang/stream"
  _ "image/jpeg"
)

type Camera struct {
  Id string
  Name string
  VideoFormat string
  AudioFormat string
  PixFormat string
  FrameRate int
  VideoSDP string
  AudioSDP string
  HasAudio bool
  HasVideo bool
  HasPan bool
  HasTilt bool
  HasZoom bool
  AcceptsAudio bool
  provider string
  providerAddress string
  providerConfig *_struct.Struct //Arbitrary struct defined in the protobuf
}

func NewCamera(config *cameraAPI.CameraConfig) (*Camera, error) {
  c := &Camera{}

  c.SetId(config.Camera.Id)
  c.SetName(config.Name)
  c.SetProvider(config.Provider)
  c.SetProviderAddress(config.ProviderAddress)
  c.SetProviderConfig(config.ProviderConfig)

  //Ask the provider about the Camera's capabilities
  client, err := c.providerClient()
  if err != nil {
    return nil, err
  }
  cameraInfo, err := client.Describe(context.Background(), &providerAPI.CameraConfig{Config: c.providerConfig})
  if err != nil {
    return nil, err
  }

  c.VideoFormat  = cameraInfo.VideoFormat
  c.AudioFormat  = cameraInfo.AudioFormat
  c.PixFormat    = cameraInfo.PixFormat
  c.FrameRate    = int(cameraInfo.FrameRate)
  c.VideoSDP     = cameraInfo.VideoSDP
  c.AudioSDP     = cameraInfo.AudioSDP
  c.HasAudio     = cameraInfo.HasAudio
  c.HasVideo     = cameraInfo.HasVideo
  c.HasPan       = cameraInfo.HasPan
  c.HasTilt      = cameraInfo.HasTilt
  c.HasZoom      = cameraInfo.HasZoom
  c.AcceptsAudio = cameraInfo.AcceptsAudio

  return c, nil
}

func (s *Camera) providerClient() (providerAPI.CameraClient, error) {
  var opts []grpc.DialOption

  if len(s.providerAddress) == 0 {
    return nil, errors.New("Camera does not have a configured provider")
  }

  opts = append(opts, grpc.WithInsecure())
  conn, err := grpc.Dial(s.providerAddress, opts...)
  if err != nil {
    return nil, err
  }
  client := providerAPI.NewCameraClient(conn)

  return client, nil
}

func (s *Camera) streamClient() (streamAPI.StreamClient, error) {
  var opts []grpc.DialOption


  opts = append(opts, grpc.WithInsecure())
  conn, err := grpc.Dial("stream", opts...)
  if err != nil {
    return nil, err
  }
  client := streamAPI.NewStreamClient(conn)

  return client, nil
}

func (s *Camera) GetSnapshot() (image.Image, error) {
  client, err := s.providerClient()
  if err != nil {
    return nil, err
  }

  c := providerAPI.CameraConfig{
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
  log.Println(fmt.Sprintf("[%d] %s", s.Id, msg))
}

func FmtError(err *providerAPI.Result) string {
  if err == nil || err.ErrorKind == "" || err.Message == "" {
    return "Unknown Error"
  }

  return fmt.Sprintf("ErrorKind: %s - %s", err.ErrorKind, err.Message)
}

func (s *Camera) StopFeed() error {
  //Create a gRPC client to the camera's provider service
  provider, err := s.providerClient()
  if err != nil {
    s.Log(err.Error())
    return err
  }

  //Create a gRPC client to the stream service so we can get a list 
  //of active sessions for this camera
  stream, err := s.streamClient()
  if err != nil {
    s.Log(err.Error())
    return err
  }

  //Get a list of active sessions for this camera
  result, _ := stream.Describe(context.Background(), &streamAPI.Source{ Id: s.Id })
  if result == nil {
    m := fmt.Sprintf("Could not get a list of session for camera with ID %v", s.Id)
    return errors.New(m)
  }

  //Iterate over each session and stop the stream
  //Tell the provider to stop streaming
  for _ , session := range result.Sessions {
    go func() {
      providerSession := &providerAPI.Session{
        Id: session.Id,
      }

      result, _ := provider.StopRTPStream(context.Background(), providerSession)
      if result == nil || result.Successful != true {
        m := FmtError(result)
        s.Log(m)
      }
    }()
  }

  return nil
}

func (s *Camera) SetId(id string) (*Camera) {
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
