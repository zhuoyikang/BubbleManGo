package bubble

import (
	"agent"
	"fmt"
)


// map 语法。
var PacketHandlerMap = agent.HandlerMap{
	BZ_USERLOGINREQ: LoginPktHandler,
	BZ_USERJOINREQ:  JoinPktHandler,
}

//登陆包
func LoginPktHandler(s *agent.Session, b []byte) int {
	_, login, _ := BzReadUserLoginReq(b)
	fmt.Printf("login pkt:%v\n", login)
	s.U = MakeUserData(s, login.udid)
	go s.U.Run()
	return 0
}

//加入房间
func JoinPktHandler(s *agent.Session, b []byte) int {
	_, join, _ := BzReadUserJoinReq(b)
	fmt.Printf("join pkt:%v\n", join)
	msg := Msg{t:MSG_T_JOIN, d: &(s.U.(*UserData).mq) }
	roomMgr.mq <- msg
	return 0
}
