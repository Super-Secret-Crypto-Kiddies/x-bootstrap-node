package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"x-bootstrap-node/xutil"
	"x-bootstrap-node/xutil/ipassign"

	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var (
		listenAddr       = flag.String("address", "", "listen address (default IPv4 address)")
		listenPort       = flag.Uint("port", 9871, "listen port")
		useLogger        = flag.Bool("log", true, "use logger")
		verbosity        = flag.Uint("verbosity", 4, "log verbosity (0-5)")
		maxDialAttempts  = flag.Uint("dial-attempts", 3, "max dial attempts")
		maxInboundConn   = flag.Uint("inbound", 128, "max inbound connections")
		maxOutboundConn  = flag.Uint("outbound", 128, "max outbound connections")
		numWorkers       = flag.Uint("workers", 0, "number of workers (default number of CPU cores)")
		peerRefreshDelay = flag.Uint("refresh", 60, "delay between peer refresh attempts")
	)

	flag.Parse()

	var options = []noise.NodeOption{
		noise.WithNodeMaxDialAttempts(*maxDialAttempts),
		noise.WithNodeMaxInboundConnections(*maxInboundConn),
		noise.WithNodeMaxOutboundConnections(*maxOutboundConn),
	}

	if *listenAddr == "" {
		ip := ipassign.GetIPv4Address()
		fmt.Println("No listen address specified, listening on", ip)
		addr := noise.WithNodeBindHost(net.ParseIP(*listenAddr))
		options = append(options, addr)
	}

	if *listenAddr != "" {
		fmt.Println("Listening on", *listenAddr)
		addr := noise.WithNodeBindHost(net.ParseIP(*listenAddr))
		options = append(options, addr)
	}

	if *listenPort != 0 {
		fmt.Println("Listening on port", *listenPort)
		port := noise.WithNodeBindPort(uint16(*listenPort))
		options = append(options, port)
	}

	if *useLogger {
		logger, _ := zap.NewDevelopment(zap.AddStacktrace(zapcore.Level(*verbosity)))
		log := noise.WithNodeLogger(logger)
		options = append(options, log)
	}

	if *numWorkers != 0 {
		fmt.Println("Using", *numWorkers, "workers")
		workers := noise.WithNodeNumWorkers(*numWorkers)
		options = append(options, workers)
	}

	node := newNode(options)

	overlay := kademlia.New()
	node.Bind(overlay.Protocol())

	xutil.RegisterNodeMessages(node)

	if err := node.Listen(); err != nil {
		panic(err)
	}

	updatePeers(node, overlay, *peerRefreshDelay)

	c := make(chan interface{})
	<-c
}

func newNode(options []noise.NodeOption) *noise.Node {
	node, err := noise.NewNode(options...)

	if err != nil {
		panic(err)
	}

	return node
}

func updatePeers(node *noise.Node, k *kademlia.Protocol, peerRefreshDelay uint) {
	go func() {
		for {
			for _, peer := range k.Table().Peers() {
				node.Ping(context.Background(), peer.Address)
			}

			time.Sleep(time.Second * time.Duration(peerRefreshDelay))
		}
	}()
}
