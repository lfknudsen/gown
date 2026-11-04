package web

import (
	"context"
	"fmt"
	"uuid"

	"gown/pkg"
	"gown/proto/listen"

	"google.golang.org/grpc"
)

type GownClient struct {
	conn    *grpc.ClientConn
	speaker listen.ListenerClient
	token   string
}

func NewGownClient(port int) (client *GownClient) {
	client = new(GownClient{})

	var opts []grpc.DialOption
	connString := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.NewClient(connString, opts...)
	pkg.FATAL(err)
	client.conn = conn
	client.speaker = listen.NewListenerClient(conn)
	client.token = uuid.New().String()
	return client
}

func (client *GownClient) RequestDownload(url string) {
	id := uuid.New().String()
	link := listen.Link{
		Id:    id,
		Url:   url,
		Token: client.token,
	}
	var callOpts []grpc.CallOption
	reply, err := client.speaker.Download(context.Background(), &link, callOpts...)
	pkg.FATAL(err)
	fmt.Println("Is the download finished? ", reply.Finished)
}
