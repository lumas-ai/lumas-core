package main

import (
  "fmt"
  "log"
  "io"
  "io/ioutil"
  "flag"
  "context"
  "sync"

  yaml "gopkg.in/yaml.v2"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
  api "github.com/lumas-ai/lumas-core/protos/golang"
  _struct "github.com/golang/protobuf/ptypes/struct"
)

var (
  tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
  caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
  serverHostOverride = flag.String("server_host_override", "x.domain.com", "The server name use to verify the hostname returned by TLS handshake")
  configFile         = flag.String("config_file", "./config.yml", "The yaml configuration file")
  controller         = flag.String("controller", "127.0.0.1:5390", "The address of the camera processor in the format host:port")
  purge              = flag.Bool("purge", true, "Purge any cameras that no longer exist in the config file")
)

type config struct {
  Global struct {
    Timezone string
    Loglevel string
  }
  Cameras []struct {
    Id int
    Name string
    Provider struct {
      Name string
      Address string
      Config *_struct.Struct
    }
    Players []struct {
      Name string
      Address string
      Config *_struct.Struct
    }
  }
}

func loadConfig() (*config, error) {
  yml, err := ioutil.ReadFile(*configFile)
  if err != nil {
    fmt.Println(fmt.Sprintf("Could not read config file %s: %s", configFile, err.Error()))
    return nil, err
  }

  c := &config{}
  err = yaml.Unmarshal(yml, c)
  if err != nil {
    fmt.Println(fmt.Sprintf("Could not parse yaml: %s", err.Error()))
    return nil, err
  }

  return c, nil
}

func fmtCamera(camera *api.CameraConfig) string {
  for field := range camera.ProviderConfig.GetFields() {
    for key, value := range field {
      fmt.Println(fmt.Sprintf("%v", key))
      fmt.Println(fmt.Sprintf("%v", value))
    }
  }

  s := `
%s:
    ID: %d
    Provider: %s
    ProviderConfig: %v
`

  return fmt.Sprintf(s, camera.Name, camera.Id, camera.Provider, camera.ProviderConfig.String())
}

func main() {
  flag.Parse()

  //Load the configurations from config file
  config, err := loadConfig()
  if err != nil {
    log.Fatal("Could not load configuration file: " + err.Error())
  }

  //Set up the GRPC connection
  var opts []grpc.DialOption
  if *tls {
    creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
    if err != nil {
      log.Fatalf("Failed to create TLS credentials %v", err)
    }
    opts = append(opts, grpc.WithTransportCredentials(creds))
  } else {
    opts = append(opts, grpc.WithInsecure())
  }
  conn, err := grpc.Dial(*controller, opts...)
  if err != nil {
    log.Fatalf("fail to dial: %v", err)
  }
  defer conn.Close()
  controller := api.NewCameraClient(conn)

  if *purge {
    //Get a list of existing cameras on the controller
    var wg sync.WaitGroup

    req := &api.ListRequest{}
    client, _ := controller.List(context.Background(), req)
    for {
      camera, err := client.Recv()
      if err == io.EOF {
        break
      }
      if err != nil {
        fmt.Println(err.Error())
        continue
      }

      wg.Add(1)
      go func() {
        //Loop through the configured cameras and see if the camera returned from the controller
        //exists in the config. If not, purge it from the controller
        for _, cam := range config.Cameras {
          if int(camera.Id) == cam.Id {
            wg.Done()
            return
          }
        }

        controller.Remove(context.Background(), camera)

        wg.Done()
      }()

      wg.Wait()
    }
  }

  //For each camera, call Add()
  for _, camera := range config.Cameras {
    controller.Add(context.Background(), &api.CameraConfig{ Id: int64(camera.Id),
      Name: camera.Name,
      Provider: camera.Provider.Name,
      ProviderAddress: camera.Provider.Address,
      ProviderConfig: camera.Provider.Config,
    })
  }

  //List the cameras on the controller
  req := &api.ListRequest{}
  client, _ := controller.List(context.Background(), req)
  for {
    camera, err := client.Recv()
    if err == io.EOF {
      break
    }
    if err != nil {
      fmt.Println(err.Error())
      continue
    }

    fmt.Println(fmtCamera(camera))
  }
}
