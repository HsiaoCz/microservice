package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mircoservice/grpc/demo/pb"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
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
	cfg := jaegercfg.Configuration{
		// 采样的设置
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.21.17.5:6831",
		},
		ServiceName: "grpc-hello",
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()
	fmt.Println(tracer)
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
