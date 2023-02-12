package main

import (
	"log"
	"mircoservice/go-kit/demo2/service"

	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	svc := service.AddService{}

	sumHandler := httptransport.NewServer(
		service.MakeSumEndpoint(svc),
		service.DecodeSumRequest,
		service.EncodeResponse,
	)

	concatHandler := httptransport.NewServer(
		service.MakeConcatEndpoint(svc),
		service.DecodeCountRequest,
		service.EncodeResponse,
	)

	http.Handle("/sum", sumHandler)
	http.Handle("/concat", concatHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
