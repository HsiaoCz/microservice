package main

import (
	"mircoservice/go-kit/demo1/service"
	"net/http"

	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func main() {
	user := service.UserService{}
	endp := service.GenUserEndpoint(user)

	serverHandler := httpTransport.NewServer(endp, service.DecodeUserRequest, service.EncodeUserResponse)
	r := mux.NewRouter()
	r.Handle(`/user/{uid:\d+}`, serverHandler)
	// r.Methods("GET").Path(`/user/{uid:\d+}`).Handler(serverHandler)
	http.ListenAndServe(":3202", r)
}
