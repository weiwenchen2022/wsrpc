package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/weiwenchen2022/wsrpc"
	"github.com/weiwenchen2022/wsrpc/examples/arith"
)

var addr = flag.String("addr", "localhost:1234", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := wsrpc.Dial(ctx, *addr, "")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()

	// Synchronous call
	args := &arith.Args{A: 7, B: 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	log.Printf("Arith: %d * %d = %d", args.A, args.B, reply)

	// Asynchronous call
	quotient := new(arith.Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	<-divCall.Done // will be equal to divCall
	if divCall.Error != nil {
		log.Fatal("arith error:", err)
	}

	divReply := divCall.Reply.(*arith.Quotient)
	log.Printf("Arith: %d / %d = %d, %d", args.A, args.B, divReply.Quo, divReply.Rem)
}
