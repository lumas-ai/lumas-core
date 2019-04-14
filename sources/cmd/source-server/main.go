package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	api "github.com/lumas-ai/lumas-core/protos/golang/source"
	"github.com/lumas-ai/lumas-core/sources"
)

var (
	tls         = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile    = flag.String("cert_file", "", "The TLS cert file")
	keyFile     = flag.String("key_file", "", "The TLS key file")
	iface       = flag.String("host", "0.0.0.0", "The interface to listen on")
	redisServer = flag.String("redis_server", "redis:6379", "The Redis server address in format ip:port")
	redisPass   = flag.String("redis_pass", "", "The Redis password")
	port        = flag.Int("port", 5392, "The port to listen on")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *iface, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s, _ := sources.Init(*redisServer, *redisPass)

	grpcServer := grpc.NewServer()
	api.RegisterSourceServer(grpcServer, s)
	grpcServer.Serve(lis)
}
