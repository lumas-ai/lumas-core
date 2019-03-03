package processor

import (
  "fmt"
  "context"
  "strconv"
  "errors"
  api "github.com/lumas-ai/lumas-core/protos/golang"
)

type CameraServer struct {
  cameras []*Camera
}

func (s *CameraServer) FindCamera(id int64) (*Camera, int) {
  for i, cam := range s.cameras {
    if cam.Id == id {
      return cam, i
    }
  }

  return nil, -1
}

func (s *CameraServer) Add(ctx context.Context, config *api.CameraConfig) (*api.Result, error) {
  //If a camera with the same ID exists, replace it
  //in case there's new configurations
  cam, _ := s.FindCamera(config.Id)
  if cam != nil {
    s.Remove(ctx, config)
  }

  camera, err := NewCamera()
  if err != nil {
    r := api.Result{Successful: false, ErrorKind: "CameraCreateFailure", Message: err.Error()}
    return &r, err
  }
  s.cameras = append(s.cameras, camera)

  camera.SetId(config.Id).SetName(config.Name).SetProvider(config.Provider).SetProviderAddress(config.ProviderAddress).SetProviderConfig(config.ProviderConfig)

  r := api.Result{Successful: true}
  return &r, nil
}

func (s *CameraServer) Process(ctx context.Context, cameraID *api.CameraID) (*api.Result, error) {
  camera, _ := s.FindCamera(cameraID.Id)
  if camera == nil {
    m := "Could not find camera with ID " + strconv.FormatInt(cameraID.Id, 10)
    r := api.Result{Successful: false, ErrorKind: "CameraNotFound", Message: m}
    return &r, errors.New(m)
  }

  err := camera.ProcessFeed()
  if err != nil {
    r := api.Result{Successful: false,
      ErrorKind: "CouldNotStartFeed",
      Message: err.Error(),
    }
    return &r, err
  }

  r := api.Result{Successful: true}
  return &r, nil
}

func (s *CameraServer) List(req *api.ListRequest, stream api.Camera_ListServer) error {
  for _, cam := range s.cameras {
    r := &api.CameraConfig{
      Id: cam.Id,
      Name: cam.Name,
      Provider: cam.provider,
      ProviderAddress: cam.providerAddress,
      ProviderConfig: cam.providerConfig,
    }

    stream.Send(r)
  }

  return nil
}

func (s *CameraServer) Remove(context context.Context, camera *api.CameraConfig) (*api.Result, error) {
  cam, i := s.FindCamera(camera.Id)
  if cam == nil {
    msg := fmt.Sprintf("Could not find camera with ID %d", camera.Id)
    r := &api.Result{
      Successful: false,
      ErrorKind: "CameraNotFound",
      Message: msg,
    }

    return r, errors.New(msg)
  }

  err := cam.StopFeed()
  if err != nil {
    msg := fmt.Sprintf("Could not close camera feed with ID %d", camera.Id)
    r := &api.Result{
      Successful: false,
      ErrorKind: "CameraLeftOpen",
      Message: msg,
    }

    return r, errors.New(msg)
  }

  s.cameras = append(s.cameras[:i], s.cameras[i+1:]...)

  r := &api.Result{ Successful: true}
  return r, nil
}
