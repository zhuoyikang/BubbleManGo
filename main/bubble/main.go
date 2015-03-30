package main

import (
	"agent"
	"bubble"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

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
			fmt.Printf("%s\n", "hub")
			agentBubble.Close()
			os.Exit(1)
		case syscall.SIGINT:
			fmt.Printf("%s\n", "int")
			agentBubble.Close()
			os.Exit(2)
		case syscall.SIGTERM:
		}
	}
}

var agentBubble* agent.TcpAgent



func main() {
	go signalHandler()


	agentBubble = agent.MakeTcpAgent("127.0.0.1:3004", bubble.PacketHandlerMap)
	agentBubble.Run()

}
