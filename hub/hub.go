package hub

import (
	"agent"
	"fmt"
)

type HubGs struct {
}

// 给一个初始化的机会
func (b HubGs) Start(s *agent.Session) {
	fmt.Printf("%s\n", "start")
}

// 给一个初始化的机会
func (b HubGs) Stop(s *agent.Session) {
	if s.U == nil {
		return
	}
	u := s.U.(*UserData)
	u.Stop()
}

// 分发处理客户端消息. clientRoutine只处理登录等消息，登录后消息分发给user处理。
func (b HubGs) Handler(s *agent.Session, t int, data []byte) bool {

	//unlock
	defer func() {
		if s.U != nil && t != BZ_USERLOGINREQ {
			//u := s.U.(*UserData)
			//u.mutex.Unlock()
		}
	}()

	var hc agent.Handler
	if h, s := packetHandlerMap[t]; s == true {
		hc = h
	} else {
		hc = DefaultHandler
	}

	//lock
	if s.U != nil {
		//u := s.U.(*UserData)
		//u.mutex.Lock()
	}

	if hc(s, t, data) == 0 {
		return true
	} else {
		return false
	}
}

// map 语法。
var packetHandlerMap = agent.HandlerMap{
	BZ_USERLOGINREQ: LoginPktHandler,
}

//默认包处理函数
func DefaultHandler(s *agent.Session, t int, b []byte) int {
	//clientMsg := ClientMsg{t: t, d: b}
	//msg := Msg{t: MSG_T_CLIENT, d: clientMsg}
	//s.U.(*UserData).mq <- msg
	fmt.Printf("%s\n", "default handler")
	return 0
}

//登陆成功，初始化用户数据，并分配单独的共组routine。
func LoginPktHandler(s *agent.Session, t int, b []byte) int {
	_, login, _ := BzReadUserLoginReq(b)
	fmt.Printf("Login Pkt:%v\n", login)
	s.U = MakeUserData(s, login.udid)
	go s.U.Run()
	return 0
}
