package p2p

import (
	"context"
	"fmt"
	"time"
	"x-bootstrap-node/xutil"

	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
	"go.uber.org/zap"
)

var Node *noise.Node
var ip = GetIPAddress()
var port uint16 = 9871

func InitP2P() {
	logger, _ := zap.NewDevelopment()

	address := fmt.Sprintf("%s:%d", ip, port)

	fmt.Printf("Running on: %s\n", address)

	Node, err := noise.NewNode(
		noise.WithNodeBindPort(port),
		noise.WithNodeAddress(address),
		noise.WithNodeLogger(logger),
	)

	if err != nil {
		panic(err)
	}

	k := kademlia.New()
	Node.Bind(k.Protocol())

	xutil.RegisterNodeMessages(Node)

	Node.Handle(handle)

	if err := Node.Listen(); err != nil {
		panic(err)
	}

	go func() {
		for {
			updatePeers(Node, k)
			time.Sleep(time.Minute)
		}
	}()
}

func handle(ctx noise.HandlerContext) error {
	if (string)(ctx.Data()) == "peerinfo" {
		ctx.SendMessage(xutil.PeerInfo{
			Type:       xutil.BootstrapNode,
			Currencies: []string{},
		})
	}

	return nil
}

func updatePeers(node *noise.Node, k *kademlia.Protocol) {
	peers := k.Table().Peers()

	for _, peer := range peers {
		node.Ping(context.Background(), peer.Address)
	}
}
