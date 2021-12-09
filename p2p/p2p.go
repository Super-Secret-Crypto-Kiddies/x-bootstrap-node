package p2p

import (
	"fmt"
	"time"
	"x-bootstrap-node/xutil"

	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
)

var Node *noise.Node = nil

const p2pPort = 9871

func InitP2P() {

	var err error = nil
	if Node, err = noise.NewNode(
		noise.WithNodeBindPort(p2pPort),
		// CHANGE THIS PLS
		noise.WithNodeAddress(fmt.Sprintf("143.198.176.143:%d", p2pPort)),
	); err != nil {
		panic(err)
	}

	k := kademlia.New()
	Node.Bind(k.Protocol())

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

	go func() {
		for {
			fmt.Printf("Peers: %s\n\n", k.Discover())
			time.Sleep(30 * time.Second)
		}
	}()

}
