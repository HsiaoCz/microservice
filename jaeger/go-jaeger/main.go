package main

import (
	"log"
	"time"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
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
	span := tracer.StartSpan("go-grpc-hello")
	time.Sleep(time.Second)
	defer span.Finish()
}
