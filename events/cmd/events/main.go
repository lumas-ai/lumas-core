package main

import (
  "fmt"
  "flag"
  "net"
  "log"
  //"context"

  "google.golang.org/grpc"

  api "github.com/lumas-ai/lumas-core/protos/golang"
  event "github.com/lumas-ai/lumas-core/events"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	iface      = flag.String("host", "0.0.0.0", "The interface to listen on")
	port       = flag.Int("port", 5388, "The port to listen on")
)

func main() {
  flag.Parse()

  lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *iface, *port))
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  //ctx := context.Background()
  //ctxWithCancel, cancel := context.WithCancel(ctx)
  //defer cancel()

  s := event.EventsServer{}

  grpcServer := grpc.NewServer()
  api.RegisterEventsServer(grpcServer, &s)
  grpcServer.Serve(lis)
}
