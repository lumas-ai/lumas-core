package main

import (
  "fmt"
  "flag"

  "google.golang.org/grpc"

  . "github.com/lumas-ai/lumas-core/processor"
  api "github.com/lumas-ai/lumas-core/protos/golang"
  _struct "github.com/golang/protobuf/ptypes/struct"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	iface      = flag.String("host", "0.0.0.0", "The interface to listen on")
	port       = flag.Int("port", 10000, "The server port")
)

func main() {
  flag.Parse()

  lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", iface, port))
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }
  s := CameraServer{}
  grpcServer := grpc.NewServer()
  api.RegisterCameraServer(grpcServer, &s)
  grpcServer.Serve(lis)
}
