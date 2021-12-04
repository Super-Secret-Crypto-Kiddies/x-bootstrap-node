package p2p

import (
	"x-bootstrap-node/xutil"

	"github.com/perlin-network/noise"
)

var Node *noise.Node = nil

func InitP2P() {

	var err error = nil
	if Node, err = noise.NewNode(noise.WithNodeBindPort(9871)); err != nil {
		panic(err)
	}

	Node.Handle(func(ctx noise.HandlerContext) error {

		// handle raw []byte messages

		if (string)(ctx.Data()) == "peerinfo" {
			ctx.SendMessage(xutil.PeerInfo{
				Type:       xutil.BootstrapNode,
				Currencies: []string{},
			})
		}
		return nil
	})

	if err := Node.Listen(); err != nil {
		panic(err)
	}
}
