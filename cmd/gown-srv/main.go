package main

import (
	"context"
	"flag"
	"fmt"
	. "gown/pkg"
	. "gown/proto"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedListenerServer
	port int
}

func (s *Server) Serve(server *grpc.Server, lis net.Listener, wait chan struct{}) {
	WARN("Listening on " + lis.Addr().String())
	err := server.Serve(lis)
	if err != nil {
		wait <- struct{}{}
		log.Fatal(err)
	}

}

func (s *Server) Download(_ context.Context, link *Link) (*Reply, error) {
	Download(link.GetUrl(), make(chan int))
	finished := true
	return &Reply{
		Finished: &finished,
	}, nil
}

func main() {
	logFile := InitialiseLogging()
	defer CloseLogFile(logFile)
	flag.Usage = PrintHelp
	flag.Parse()
	DEBUG("Parsed command-line flags.")
	if *Help {
		PrintHelp()
		os.Exit(0)
	}

	SetupDirs()

	if *GenerateLists {
		GenerateGalleryList()
		GenerateVideoList()
		WARN("Finished generating lists of supported domains.")
	}

	Srv := Server{
		port: 8080,
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
	Interactive()

	for {
		select {
		case <-wait:
			return
		}
	}

	DEBUG("Done")
}
