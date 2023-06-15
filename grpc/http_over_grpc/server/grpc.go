package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/bamcop/kit/grpc/http_over_grpc"
	"github.com/hashicorp/yamux"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

const (
	DefaultClientIdentifier = 1
)

func StartGrpc(addr string) {
	log.Println("launching tcp server...")

	// start tcp listener on all interfaces
	// note that each connection consumes a file descriptor
	// you may need to increase your fd limits if you have many concurrent clients
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("could not listen", slog.Any("err", err))
		panic(err)
	}
	defer ln.Close()

	for {
		log.Println("waiting for incoming TCP connections...")
		// Accept blocks until there is an incoming TCP connection
		incoming, err := ln.Accept()
		if err != nil {
			log.Fatalf("couldn't accept %s", err)
		}

		incomingConn, err := yamux.Client(incoming, yamux.DefaultConfig())
		if err != nil {
			log.Fatalf("couldn't create yamux %s", err)
		}

		log.Println("starting a gRPC server over incoming TCP connection")

		var conn *grpc.ClientConn
		// gRPC dial over incoming net.Conn
		conn, err = grpc.Dial(":7777", grpc.WithInsecure(),
			grpc.WithDialer(func(target string, timeout time.Duration) (net.Conn, error) {
				return incomingConn.Open()
			}),
		)

		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}

		// handle connection in goroutine so we can accept new TCP connections
		go handleConn(conn)
	}
}

func handleConn(conn *grpc.ClientConn) {
	c := pb.NewHTTPClient(conn)

	// TODO: 此处可以调用客户端的方法获取客户端信息
	slog.Info("客户端注册")

	LocalServer.Lock()
	defer LocalServer.Unlock()
	LocalServer.conns[DefaultClientIdentifier] = c
}

type localServer struct {
	sync.Mutex
	conns map[int64]pb.HTTPClient
}

var (
	LocalServer = &localServer{
		conns: map[int64]pb.HTTPClient{},
	}
)

func (s *localServer) conn(identifier int64) (pb.HTTPClient, error) {
	s.Lock()
	defer s.Unlock()

	conn, ok := s.conns[identifier]
	if !ok {
		slog.Error("客户端不在线")
		return nil, fmt.Errorf("客户端不在线")
	}
	return conn, nil
}
