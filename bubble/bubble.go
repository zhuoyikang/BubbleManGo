package bubble

import (
	"agent"
	"fmt"
)

type BubbleGs struct {
}

// 给一个初始化的机会
func (b BubbleGs) Start(s *agent.Session) {
	fmt.Printf("%s\n", "start")
}

// 给一个初始化的机会
func (b BubbleGs) Stop(s *agent.Session) {
	if s.U != nil {
		s.U.(*UserData).Stop()
	}
	fmt.Printf("%s\n", "stop")
}

// 分发处理客户端消息. clientRoutine只处理登录等消息，登录后消息分发给user处理。
func (b BubbleGs) Handler(s *agent.Session, t int, data []byte) bool {
	var hc agent.Handler
	if h, s := packetHandlerMap[t]; s == true {
		hc = h
	} else {
		hc = DefaultHandler
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
	clientMsg := ClientMsg{t: t, d: b}
	msg := Msg{t: MSG_T_CLIENT, d: clientMsg}
	s.U.(*UserData).mq <- msg
	return 0
}

//登陆成功，初始化用户数据，并分配单独的共组routine。
func LoginPktHandler(s *agent.Session, t int, b []byte) int {
	_, login, _ := BzReadUserLoginReq(b)
	fmt.Printf("login pkt:%v\n", login)
	s.U = MakeUserData(s, login.udid)
	go s.U.Run()
	return 0
}
