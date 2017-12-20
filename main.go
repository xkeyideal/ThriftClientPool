package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/xkeyideal/ThriftClientPool/client"
	"github.com/xkeyideal/ThriftClientPool/server"
	"github.com/xkeyideal/ThriftClientPool/thriftPool"
)

func mainServer() {
	go server.RunServer()
	go server.RunHttpServer()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-signals

	os.Exit(0)
}

func mainClient() {
	client.GlobalRpcPool = thriftPool.NewThriftPool("localhost", "9999", 10, 32, 600, client.Dial, client.Close)
}

func main() {
	mainServer()
}
