package main

import (
	"flag"

	"fileman-net/internal/client"
	"fileman-net/internal/common"
	"fileman-net/internal/server"
)

func main() {
	var ctx common.AppContext

	flag.BoolVar(&ctx.RunAsServer, "server", common.DefaultRunAsServer, "run as server")

	flag.UintVar(&ctx.Port, "port", common.DefaultPort, "port to connect")
	flag.StringVar(&ctx.Addr, "address", common.DefaultAddr, "address to connect")

	flag.StringVar(&ctx.ServerWorkingDir, "server-wd", common.DefaultServerPath, "server working dir")

	flag.Parse()

	if ctx.RunAsServer {
		server.RunServer(&ctx)
	} else {
		client.RunClient(&ctx)
	}
}
