package wsrpc

import (
	"context"
	"net/rpc"
	"net/rpc/jsonrpc"
	"net/url"
	"path"

	"nhooyr.io/websocket"
)

// Dial connects to an websocket RPC server
// at the specified address and path.
// If rpcPath is empty use DefaultWsRPCPath.
func Dial(ctx context.Context, addr, rpcPath string) (*rpc.Client, error) {
	if rpcPath == "" {
		rpcPath = DefaultWsRPCPath
	}
	u, err := url.Parse("ws://" + path.Join(addr, rpcPath))
	if err != nil {
		return nil, err
	}

	c, _, err := websocket.Dial(ctx, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return jsonrpc.NewClient(websocket.NetConn(context.Background(), c, websocket.MessageText)), nil
}
