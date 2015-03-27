package main

import (
	"agent"
	"bubble"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//
func signalHandler() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGINT)
	for {
		msg := <-signalChan
		switch msg {
		case syscall.SIGHUP:
			os.Exit(1)
		case syscall.SIGINT:
			os.Exit(2)
		case syscall.SIGTERM:
		}
	}
}

func main() {
	go signalHandler()
	agentTcp := agent.MakeTcpAgent("127.0.0.1:3004",bubble.PacketHandlerMap)
	go agentTcp.Run()
	time.Sleep(time.Second * 10)
	fmt.Printf("%s\n", "begin to close")
	agentTcp.Close()
	fmt.Printf("%s\n", "end close")
}
