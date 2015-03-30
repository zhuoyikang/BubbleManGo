package agent

import (
	"fmt"
	"net"
	"os"
	"sync"
)

/*------------------------------------------------------------------------------
 agent:监听端口分发进程
------------------------------------------------------------------------------*/

type Handler func(*Session, []byte) int
type HandlerMap map[int]Handler

type TcpAgent struct {
	ip             string
	hmap           HandlerMap
	connectionPool map[net.Conn]*Session
	listener       net.Listener
	wg             *sync.WaitGroup
}

//创建一个agent.
func MakeTcpAgent(ip string, h HandlerMap) *TcpAgent {
	agent := TcpAgent{ip: ip}
	agent.wg = &sync.WaitGroup{}
	agent.connectionPool = make(map[net.Conn]*Session)
	agent.hmap = h
	return &agent
}

//开始工作
func (agent *TcpAgent) Run() {
	listener, err := net.Listen("tcp", agent.ip)
	if err != nil {
		os.Exit(1)
	}

	defer func() {
		agent.wg.Done()
		listener.Close()
	}()

	agent.listener = listener
	agent.wg.Add(1)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error Accept %s\n", err.Error())
			return
		}
		session := MakeSession(conn)
		agent.connectionPool[conn] = &session
		go session.HandleClient(agent)
	}
}

// 同步关闭所有的连接和连接处理routine.
func (agent *TcpAgent) Close() {
	agent.listener.Close()
	for k, _ := range agent.connectionPool {
		k.Close()
	}
	agent.wg.Wait()
}
