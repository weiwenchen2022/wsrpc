package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/weiwenchen2022/wsrpc"
	"github.com/weiwenchen2022/wsrpc/examples/arith"
)

var addr = flag.String("addr", ":1234", "http service address")

func main() {
	flag.Parse()

	s := wsrpc.NewServer("")
	s.Register(new(arith.Arith))

	log.Println("server serving at", *addr)
	log.Fatal(http.ListenAndServe(*addr, s))
}
