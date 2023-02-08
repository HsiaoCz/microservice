package main

import (
	"context"
	"flag"
	"log"
	"mircoservice/grpc/demo/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

//hello client端

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", ":3001", "the addr to connect")
	name = flag.String("name", defaultName, "name to hello")
)

func main() {
	flag.Parse()

	// 连接到server端,此处禁用安全传输
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewHelloServiceClient(conn)
	//执行rpc调用
	md := metadata.Pairs("woqu", "nima")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("hello:%s", r.GetReply())
}
