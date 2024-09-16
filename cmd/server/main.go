package main

import (
	"context"
	"flag"
	"log"
	"microkit/proto/gohello"
	"net"

	"google.golang.org/grpc"
)

var port = flag.String("port", ":8080", "The port to listen on")

type HelloServer struct {
	gohello.UnimplementedHelloServer
}

func (h *HelloServer) SayHello(ctx context.Context, req *gohello.Reqmsg) (*gohello.Response, error) {
	return &gohello.Response{Message: *port + " ::hello " + req.Name}, nil
}

func main() {
	//port := flag.String("port", "", "The port to listen on")
	flag.Parse()
	//服务
	src := grpc.NewServer()
	gohello.RegisterHelloServer(src, &HelloServer{})
	l, err := net.Listen("tcp", *port)
	if err != nil {
		panic(err)
	}

	log.Println("listen" + *port)
	if err := src.Serve(l); err != nil {
		panic(err)
	}
}
