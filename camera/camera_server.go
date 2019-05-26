package camera

import (
  "fmt"
  "context"
  "errors"
  cameraAPI "github.com/lumas-ai/lumas-core/protos/golang/camera"
)

type CameraServer struct {
  cameras []*Camera
}

func (s *CameraServer) FindCamera(id string) (*Camera, int) {
  for i, cam := range s.cameras {
    if cam.Id == id {
      return cam, i
    }
  }

  return nil, -1
}

func result(errorKind string, errorMessage string) *cameraAPI.Result {
  if errorKind == "" && errorMessage == "" {
    return &cameraAPI.Result{
      Successful: true,
    }
  }
  return &cameraAPI.Result{
    Successful: false,
    ErrorKind:  errorKind,
    Message:    errorMessage,
  }
}

func (s *CameraServer) Add(ctx context.Context, config *cameraAPI.CameraConfig) (*cameraAPI.Result, error) {
  //If a camera with the same ID exists, replace it
  //in case there's new configurations
  cam, i := s.FindCamera(config.Camera.Id)
  if cam != nil {
    s.cameras = append(s.cameras[:i], s.cameras[i+1:]...)
  }

  camera, err := NewCamera(config)
  if err != nil {
    r := cameraAPI.Result{Successful: false, ErrorKind: "CameraCreateFailure", Message: err.Error()}
    return &r, err
  }
  s.cameras = append(s.cameras, camera)

  //Return a successful result
  r := cameraAPI.Result{Successful: true}
  return &r, nil
}

func (s *CameraServer) Stream(ctx context.Context, cameraID *cameraAPI.CameraID) (*cameraAPI.Result, error) {
  r := cameraAPI.Result{Successful: true}
  return &r, nil
}

func (s *CameraServer) Describe(context context.Context, camera *cameraAPI.CameraID) (*cameraAPI.CameraInfo, error) {
  cam, _ := s.FindCamera(camera.Id)
  if cam == nil {
    msg := fmt.Sprintf("Could not find camera with ID %d", camera.Id)
    r := &cameraAPI.CameraInfo{}
    return r, errors.New(msg)
  }

  c := &cameraAPI.CameraInfo{
    Name: cam.Name,
    VideoFormat: cam.VideoFormat,
    AudioFormat: cam.AudioFormat,
    PixFormat: cam.PixFormat,
    FrameRate: int32(cam.FrameRate),
    Provider: cam.provider,
    VideoSDP: cam.VideoSDP,
    AudioSDP: cam.AudioSDP,
  }

  return c, nil
}

func (s *CameraServer) List(req *cameraAPI.ListRequest, stream cameraAPI.Camera_ListServer) error {
  for _, cam := range s.cameras {
    c := &cameraAPI.CameraInfo{
      Id: cam.Id,
    }

    r := &cameraAPI.CameraConfig{
      Camera: c,
      Name: cam.Name,
      Provider: cam.provider,
      ProviderAddress: cam.providerAddress,
      ProviderConfig: cam.providerConfig,
    }

    stream.Send(r)
  }

  return nil
}

func (s *CameraServer) Remove(context context.Context, camera *cameraAPI.CameraID) (*cameraAPI.Result, error) {
  cam, i := s.FindCamera(camera.Id)
  if cam == nil {
    msg := fmt.Sprintf("Could not find camera with ID %d", camera.Id)
    r := &cameraAPI.Result{
      Successful: false,
      ErrorKind: "CameraNotFound",
      Message: msg,
    }

    return r, errors.New(msg)
  }

  err := cam.StopFeed()
  if err != nil {
    msg := fmt.Sprintf("Could not close camera feed with ID %s - %v", camera.Id, err)
    r := &cameraAPI.Result{
      Successful: false,
      ErrorKind: "CameraLeftOpen",
      Message: msg,
    }

    return r, errors.New(msg)
  }

  s.cameras = append(s.cameras[:i], s.cameras[i+1:]...)

  r := &cameraAPI.Result{ Successful: true}
  return r, nil
}
