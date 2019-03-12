package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

  "github.com/mitchellh/go-homedir"
	"github.com/golang/protobuf/jsonpb"
	cli "gopkg.in/urfave/cli.v1"
	altsrc "gopkg.in/urfave/cli.v1/altsrc"
	"gopkg.in/yaml.v2"

	_struct "github.com/golang/protobuf/ptypes/struct"
	api "github.com/lumas-ai/lumas-core/protos/golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const version = "0.1.0"

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
	serverHostOverride = flag.String("server_host_override", "x.domain.com", "The server name use to verify the hostname returned by TLS handshake")
	controller         = flag.String("controller", "127.0.0.1:5390", "The address of the camera processor in the format host:port")
	purge              = flag.Bool("purge", true, "Purge any cameras that no longer exist in the config file")
)

type client struct {
	grpcConn     *grpc.ClientConn
	cameraClient *api.CameraClient
}

type config struct {
	Global struct {
		Timezone string
		Loglevel string
	}
	Cameras []struct {
		Id       int
		Name     string
		Provider struct {
			Name    string
			Address string
			Config  map[string]interface{}
		}
		Players []struct {
			Name    string
			Address string
			Config  *_struct.Struct `protobuf:"bytes,304,opt,name=config,yaml=config" yaml:"config,omitempty"`
		}
	}
}

func loadConfig(file string) (*config, error) {
	yml, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not read config file %s: %s", file, err.Error()))
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
	s := `
%s:
    ID: %d
    Provider: %s
    ProviderConfig: %v
`

	return fmt.Sprintf(s, camera.Name, camera.Id, camera.Provider, camera.ProviderConfig.String())
}

func newClient() (api.CameraClient, error) {
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

	grpcConn, err := grpc.Dial(*controller, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	camc := api.NewCameraClient(grpcConn)
	return camc, nil
}

func apply(file string, purge bool) error {
	//Load the configurations from config file
	config, err := loadConfig(file)
	if err != nil {
		log.Fatal("Could not load configuration file: " + err.Error())
	}

	client, err := newClient()
	if err != nil {
		fmt.Println("Could not create client connection: " + err.Error())
	}

	if purge {
		//Get a list of existing cameras on the controller
		var wg sync.WaitGroup

		req := &api.ListRequest{}
		stream, _ := client.List(context.Background(), req)
		for {
			camera, err := stream.Recv()
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

				id := &api.CameraID{Id: camera.Id}
				client.Remove(context.Background(), id)

				wg.Done()
			}()

		}
		wg.Wait()
	}

	//For each camera, call Add()
	for _, camera := range config.Cameras {

		c := &api.CameraConfig{Id: int64(camera.Id),
			Name:            camera.Name,
			Provider:        camera.Provider.Name,
			ProviderAddress: camera.Provider.Address,
			ProviderConfig:  &_struct.Struct{},
		}

		/*This is ugly. In order to package up arbitrary maps into a GRPC call,
		  we must first marshal the data object to JSON, then use jsonpb to Marshal
		  the json into a _struct.Struct object as defined in the camera protobuf
		*/
		j, _ := json.Marshal(camera.Provider.Config)
		m := jsonpb.Unmarshaler{}
		err := m.Unmarshal(strings.NewReader(string(j)), c.ProviderConfig)
		if err != nil {
			fmt.Println(err)
			return err
		}

		addResults, err := client.Add(context.Background(), c)
		if err != nil || addResults == nil {
			fmt.Println(fmt.Sprintf("Could not add camera %s", camera.Name))
			continue
		}
		if addResults.Successful != true {
			fmt.Println(fmt.Sprintf("Could not add camera %s. Error: %s - %s ", camera.Name, addResults.ErrorKind, addResults.Message))
			continue
		}

		processResults, err := client.Process(context.Background(), &api.CameraID{Id: int64(camera.Id)})
		if err != nil || processResults == nil {
			fmt.Println(fmt.Sprintf("Could not process camera %s", camera.Name))
			continue
		}
		if processResults.Successful != true {
			fmt.Println(fmt.Sprintf("Could not process feed on camera %s. Error: %s - %s ", camera.Name, processResults.ErrorKind, processResults.Message))
			continue
		}
	}

	return nil
}

func list() error {
	client, err := newClient()
	if err != nil {
		fmt.Println("Could not create client: " + err.Error())
	}

	req := &api.ListRequest{}
	stream, err := client.List(context.Background(), req)
	if err != nil {
		return err
	}

	for {
		camera, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println(fmtCamera(camera))
	}

	return nil

}

func remove(cameraID int) error {
	client, err := newClient()
	if err != nil {
		fmt.Println("Could not create client: " + err.Error())
	}

	req := &api.CameraID{Id: int64(cameraID)}
	result, err := client.Remove(context.Background(), req)
	if err != nil || result == nil {
		msg := fmt.Sprintf("Could not process camera %d", cameraID)
		fmt.Println(msg)
		return err
	}
	if result.Successful != true {
		msg := fmt.Sprintf("Could not remove camera %d\nErrorKind: %s\nErrorMsg: %s", cameraID, result.ErrorKind, result.Message)
		fmt.Println(msg)
		return errors.New(msg)
	}

	return nil
}

func stop(cameraID int) error {
	client, err := newClient()
	if err != nil {
		fmt.Println("could not create client: " + err.Error())
	}

	req := &api.CameraID{Id: int64(cameraID)}
	result, err := client.Stop(context.Background(), req)
	if result.Successful != true {
		fmt.Println(fmt.Sprintf("could not stop camera %d\nerrorkind: %s\nerrormsg: %s", cameraID, result.ErrorKind, result.Message))
	}

	return nil
}

func process(cameraID int) error {
	client, err := newClient()
	if err != nil {
		fmt.Println("could not create client: " + err.Error())
	}

	req := &api.CameraID{Id: int64(cameraID)}
	result, err := client.Process(context.Background(), req)
	if result.Successful != true {
		fmt.Println(fmt.Sprintf("could not stop camera %d\nerrorkind: %s\nerrormsg: %s", cameraID, result.ErrorKind, result.Message))
	}

	return nil
}

func main() {
	home, _ := homedir.Dir()

	cfgFile := filepath.Join(home, ".lumasctl.cfg")

	app := cli.NewApp()
	app.Name = "lumasctl"
	app.Usage = "Control Lumas cameras and services"
	app.Version = version

	flags := []cli.Flag{
		altsrc.NewStringFlag(
			cli.StringFlag{
				Name:        "controller",
				Usage:       "The address and port of the Lumas controller '1.2.3.4:5389'",
				Destination: controller,
			},
		),

		cli.StringFlag{
			Name:  "client_config, c",
			Value: cfgFile,
			Usage: "The lumasctl global arguments configuration file",
		},
	}

	app.Flags = flags
	app.Before = altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewYamlSourceFromFlagFunc("client_config"))
	app.Commands = []cli.Command{
		{
			Name:      "apply",
			Usage:     "Apply a configuration to the controller",
			UsageText: "apply <options> <config file>",
			Action: func(c *cli.Context) error {
				err := apply(c.Args().Get(0), *purge)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "purge",
					Usage:       "Purge cameras from controller that are not defined in the configuration file",
					Destination: purge,
				},
			},
		},
		{
			Name:  "camera",
			Usage: "Manage cameras",
			Subcommands: cli.Commands{
				cli.Command{
					Name:  "list",
					Usage: "List all cameras on the controller",
					Action: func(c *cli.Context) error {
						err := list()
						if err != nil {
							return cli.NewExitError(err.Error(), 1)
						}
						return nil
					},
				},
				cli.Command{
					Name:      "remove",
					Usage:     "Remove the camera and stop processing its feed",
					UsageText: "remove <camera ID>",
					Action: func(c *cli.Context) error {
						id, _ := strconv.Atoi(c.Args().Get(0))
						err := remove(id)
						if err != nil {
							return cli.NewExitError(err.Error(), 1)
						}
						return nil
					},
				},
				cli.Command{
					Name:      "stop",
					Usage:     "Stop processing the camera's feed",
					UsageText: "stop <camera ID>",
					Action: func(c *cli.Context) error {
						id, _ := strconv.Atoi(c.Args().Get(0))
						err := stop(id)
						if err != nil {
							return cli.NewExitError(err.Error(), 1)
						}
						return nil
					},
				},
				cli.Command{
					Name:      "process",
					Usage:     "Start processing a camera's feed",
					UsageText: "process <camera ID>",
					Action: func(c *cli.Context) error {
						id, _ := strconv.Atoi(c.Args().Get(0))
						err := process(id)
						if err != nil {
							return cli.NewExitError(err.Error(), 1)
						}
						return nil
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
