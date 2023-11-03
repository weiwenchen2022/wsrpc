# wsrpc
a tiny websocket rpc library.

## Uasage

### define service

```go
type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (*Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (*Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}
```

### start a server

```go
s := wsrpc.NewServer("")
s.Register(new(Arith))
http.ListenAndServe(":3000", s)
```

### start a client

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
cl, err := wsrpc.Dial(ctx, "http://localhost:3000","")
if err != nil {
	log.Fatal("dialing:", err)
}
defer cl.Close()

// Synchronous call
args := &server.Args{7,8}
var reply int
err = client.Call("Arith.Multiply", args, &reply)
if err != nil {
	log.Fatal("arith error:", err)
}
fmt.Printf("Arith: %d * %d = %d", args.A, args.B, reply)

// Asynchronous call
quotient := new(Quotient)
divCall := client.Go("Arith.Divide", args, quotient, nil)
<-divCall.Done	// will be equal to divCall
// check errors, print, etc.
```