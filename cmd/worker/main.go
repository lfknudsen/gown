package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"uuid"

	. "gown/pkg"
	"gown/proto/convoy"
	work "gown/proto/worker"

	"google.golang.org/grpc"
)

func main() {
	logFile := InitialiseLogging()
	defer CloseLogFile(logFile)
	ParseFlags()
	SetupDirs()

	if *GenerateLists {
		GenerateGalleryList()
		GenerateVideoList()
		WARN("Finished generating lists of supported domains.")
	}

	client, _, wait := NewWorker()
	port := uint32(*Port)
	subscriber := convoy.Subscriber{
		Uuid: uuid.New().String(),
		Port: port,
	}
	_, err := client.Await(context.Background(), &subscriber)
	if err != nil {
		log.Fatalf("Worker error: %s", err)
	}

	// Listen endlessly for message queue to send something to it
	for {
		select {
		case <-wait:
			os.Exit(0)
		}
	}
}

type WorkServer struct {
	work.UnimplementedWorkerServer
	port int
}

func (w WorkServer) Serve(server *grpc.Server, lis net.Listener, wait chan struct{}) {
	WARN("Listening on " + lis.Addr().String())
	err := server.Serve(lis)
	if err != nil {
		wait <- struct{}{}
		log.Fatal(err)
	}
}

func (w WorkServer) Download(link *convoy.Link, stream grpc.ServerStreamingServer[convoy.Progress]) error {
	err := stream.Send(&convoy.Progress{
		LinkUuid:   link.Id,
		Percentage: float32(math.NaN()),
		DownStatus: convoy.Status_added,
	})
	if err != nil {
		fmt.Println(errors.New("Error sending progress message: " + err.Error()))
		return err
	}
	err = Download(link.Url, nil)
	if err != nil {
		fmt.Println(errors.New("Error downloading file: " + err.Error()))
		_ = stream.Send(&convoy.Progress{
			LinkUuid:   link.Id,
			Percentage: float32(math.NaN()),
			DownStatus: convoy.Status_failed,
		})
		return err
	}
	err = stream.Send(&convoy.Progress{
		LinkUuid:   link.Id,
		Percentage: 100.0,
		DownStatus: convoy.Status_finished,
	})
	if err != nil {
		fmt.Println(errors.New("Error sending progress message: " + err.Error()))
	}
	return err
}

func NewWorker() (client convoy.BroadcastClient, server WorkServer, wait chan struct{}) {
	var clientOpts []grpc.DialOption
	connString := fmt.Sprintf("localhost:%d", *Port)
	conn, err := grpc.NewClient(connString, clientOpts...)
	FATAL(err)
	client = convoy.NewBroadcastClient(conn)

	var srvOpts []grpc.ServerOption
	server = WorkServer{
		port: *Port,
	}
	srv := grpc.NewServer(srvOpts...)
	work.RegisterWorkerServer(srv, &server)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", server.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	wait = make(chan struct{})
	go server.Serve(srv, lis, wait)
	return client, server, wait
}
