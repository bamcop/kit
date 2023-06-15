package main

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/bamcop/kit/grpc/http_over_grpc"
	"github.com/hashicorp/yamux"
	"github.com/parnurzeal/gorequest"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

func StartLocalGrpc(addr string) {
	for i := 0; ; i++ {
		if i != 0 {
			time.Sleep(time.Second * 3)
		}

		conn, err := net.DialTimeout("tcp", addr, time.Second*5)
		if err != nil {
			slog.Error("error dialing", slog.Any("err", err))
			continue
		}

		srvConn, err := yamux.Server(conn, yamux.DefaultConfig())
		if err != nil {
			slog.Error("couldn't create yamux server", slog.Any("err", err))
			continue
		}

		// create a server instance
		s := server{}

		// create a gRPC server object
		grpcServer := grpc.NewServer()

		// attach the Ping service to the server
		pb.RegisterHTTPServer(grpcServer, &s)

		// start the gRPC erver
		slog.Info("launching gRPC server over TCP connection...")
		if err := grpcServer.Serve(srvConn); err != nil {
			slog.Error("failed to serve", slog.Any("err", err))
			continue
		}
	}
}

// server is used to implement EchoServer.
type server struct {
	pb.UnimplementedHTTPServer
	client pb.HTTPClient
	cc     *grpc.ClientConn
}

func (s server) Handle(ctx context.Context, in *pb.HTTPRequest) (*pb.HTTPResponse, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d%s", HTTPLocalPort, in.Url)

	//fmt.Println(URL)
	//fmt.Println(string(data.MarshalIndent(in).Unwrap()))

	req := gorequest.New()
	req.Url = URL
	req.Method = in.Method
	for _, header := range in.Headers {
		req.Header.Del(header.Key)
		for _, value := range header.Values {
			req.Header.Set(header.Key, value)
		}
	}
	r, b, errs := req.Send(in.Body).EndBytes()
	if errs != nil {
		return nil, errs[0]
	}

	resp := &pb.HTTPResponse{
		Code:    int32(r.StatusCode),
		Headers: []*pb.Header{},
		Body:    b,
	}
	for key, values := range r.Header {
		values := values
		resp.Headers = append(resp.Headers, &pb.Header{
			Key:    key,
			Values: values,
		})
	}
	return resp, nil
}
