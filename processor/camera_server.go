package processor

import (
  "context"
  "strconv"
  "errors"
  api "github.com/lumas-ai/lumas-core/protos/golang"
)

type CameraServer struct { }

func (s *CameraServer) AddCamera(ctx context.Context, config *api.CameraConfig) (*api.Result, error) {
  camera, err := NewCamera()
  if err != nil {
    r := api.Result{Successful: false, ErrorKind: "CameraCreateFailure", Message: err.Error()}
    return &r, err
  }

  camera.SetId(config.Id).SetProvider(config.Name).SetProvider(config.Provider).SetProviderAddress(config.ProviderAddress).SetProviderConfig(config.ProviderConfig)

  r := api.Result{Successful: true}
  return &r, nil
}

func (s *CameraServer) ProcessFeed(ctx context.Context, cameraID *api.CameraID) (*api.Result, error) {
  camera := FindCamera(cameraID.Id)
  if camera == nil {
    m := "Could not find camera with ID " + strconv.FormatInt(cameraID.Id, 10)
    r := api.Result{Successful: false, ErrorKind: "CameraNotFound", Message: m}
    return &r, errors.New(m)
  }

  go camera.ProcessFeed()

  r := api.Result{Successful: true}
  return &r, nil
}
