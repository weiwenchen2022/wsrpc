// wsrpc is a minimal and idiomatic WebSocket rpc library for Go.
package wsrpc

import (
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"

	"nhooyr.io/websocket"
)

const DefaultWsRPCPath = "/_wsRPC_"

type Server struct {
	// serveMux routes the various endpoints to the appropriate handler.
	serveMux http.ServeMux

	rpcServer *rpc.Server
}

// NewServer returns a new Server at the specified rpc path.
// If rpcPath is empty use DefaultWsRPCPath.
func NewServer(rpcPath string) *Server {
	s := &Server{
		rpcServer: rpc.NewServer(),
	}
	if rpcPath == "" {
		rpcPath = DefaultWsRPCPath
	}
	s.serveMux.HandleFunc(rpcPath, s.rpcHandler)

	return s
}

func (s *Server) rpcHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.CloseNow()

	s.rpcServer.ServeCodec(jsonrpc.NewServerCodec(websocket.NetConn(r.Context(), c, websocket.MessageText)))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *Server) Register(rcvr any) error {
	return s.rpcServer.Register(rcvr)
}

func (s *Server) RegisterName(name string, rcvr any) error {
	return s.rpcServer.RegisterName(name, rcvr)
}
