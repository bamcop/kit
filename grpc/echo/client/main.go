package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/bamcop/kit/grpc/echo"
	"github.com/hashicorp/yamux"
	"google.golang.org/grpc"
)

// server is used to implement EchoServer.
type server struct {
	pb.UnimplementedPingServer
	client pb.PingClient
	cc     *grpc.ClientConn
}

func (s server) SayHello(ctx context.Context, in *pb.PingMessage) (*pb.PingMessage, error) {
	return &pb.PingMessage{Greeting: "hello"}, nil
}

// TCP client and GRPC server
func main() {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:8081", time.Second*5)
	if err != nil {
		log.Fatalf("error dialing: %s", err)
	}

	srvConn, err := yamux.Server(conn, yamux.DefaultConfig())
	if err != nil {
		log.Fatalf("couldn't create yamux server: %s", err)
	}

	// create a server instance
	s := server{}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the Ping service to the server
	pb.RegisterPingServer(grpcServer, &s)

	// start the gRPC erver
	log.Println("launching gRPC server over TCP connection...")
	if err := grpcServer.Serve(srvConn); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
