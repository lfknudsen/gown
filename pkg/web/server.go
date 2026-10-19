package web

import (
	"context"
	"fmt"
	"gown/pkg"
	. "gown/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedListenerServer
	port int
}

func (s *Server) Serve(server *grpc.Server, lis net.Listener, wait chan struct{}) {
	pkg.WARN("Listening on " + lis.Addr().String())
	err := server.Serve(lis)
	if err != nil {
		wait <- struct{}{}
		log.Fatal(err)
	}
}

func (s *Server) Download(_ context.Context, link *Link) (*Reply, error) {
	err := pkg.Download(link.GetUrl(), make(chan int))
	pkg.ERROR(err)
	return &Reply{
		Finished: err == nil,
	}, err
}

func (s *Server) Queue(_ context.Context, link *Link) (*Confirmation, error) {
	err := pkg.Queue(link.GetUrl())
	pkg.ERROR(err)
	return &Confirmation{
		Added: err == nil,
	}, err
}

func StartServer(port int) chan struct{} {
	Srv := Server{
		port: port,
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", Srv.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	RegisterListenerServer(grpcServer, &Srv)

	wait := make(chan struct{})
	go Srv.Serve(grpcServer, lis, wait)
	return wait
}
