package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

var (
	listenAddr = flag.String("listenAddr", ":3301", "set http server listen Addr")
)

// 我们简单写一个http服务，然后使用consul来做健康检查
func main() {
	flag.Parse()
	log.Println("the server is running on port:", *listenAddr)
	http.HandleFunc("/heath", ConuslHeath)
	http.ListenAndServe("172.20.115.6"+*listenAddr, nil)
}

type Message struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func ConuslHeath(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	msg := &Message{Code: 200, Msg: "Hello"}
	json.NewEncoder(w).Encode(msg)
}
