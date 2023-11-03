package wsrpc

import (
	"context"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/weiwenchen2022/wsrpc/examples/arith"
)

func Benchmark(b *testing.B) {
	rawurl, closeFn := setupTest(b)
	defer closeFn()

	u, err := url.Parse(rawurl)
	if err != nil {
		b.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cl, err := Dial(ctx, u.Hostname()+":"+u.Port(), "")
	if err != nil {
		b.Fatal(err)
	}
	defer cl.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		args := &arith.Args{A: 7, B: 8}
		var reply int
		err = cl.Call("Arith.Multiply", args, &reply)
		if err != nil {
			b.Fatal(err)
		}
		if args.A*args.B != reply {
			b.Fatalf("%d * %d = %d want %d", args.A, args.B, reply, args.A*args.B)
		}
	}
}

func setupTest(t testing.TB) (url string, closeFn func()) {
	s := NewServer("")
	s.Register(new(arith.Arith))
	server := httptest.NewServer(s)
	return server.URL, func() {
		server.Close()
	}
}
