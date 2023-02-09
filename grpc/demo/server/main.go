package main

import (
	"context"
	"fmt"
	"log"
	"mircoservice/grpc/demo/pb"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// hello server端
// grpc 首先一个结构体
type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	//服务端接收metadata
	// 这里ok代表什么呢？
	// 另外 直接打印会将所有的东西都打印出来
	// 这里如果只需要一部分的内容的话,可以筛选以下
	// 这里单个key对应的值是一个[]string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println("get metadata error")
	}
	value, ok := md["woqu"]
	if !ok {
		fmt.Println("no such key")
	}
	fmt.Println("woqu:", value)
	return &pb.HelloResponse{Reply: in.Name}, nil
}

func main() {
	//监听端口
	listen, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()                       //创建服务器
	pb.RegisterHelloServiceServer(s, &server{}) //在服务端注册

	//启动服务
	err = s.Serve(listen)
	if err != nil {
		log.Fatal(err)
	}
}
