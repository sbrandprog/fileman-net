package main

import (
	"filemannet/client"
	"filemannet/common"
	"filemannet/server"
	"flag"
)

func main() {
	var ctx common.AppContext

	flag.UintVar(&ctx.Port, "port", common.DefaultPort, "port to connect")
	flag.StringVar(&ctx.Addr, "address", common.DefaultAddr, "address to connect")
	flag.BoolVar(&ctx.RunAsServer, "server", false, "run as server")

	flag.Parse()

	if ctx.RunAsServer {
		server.RunServer(&ctx)
	} else {
		client.RunClient(&ctx)
	}
}
