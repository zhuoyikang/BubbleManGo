/*------------------------------------------------------------------------------
 处理客户端消息
------------------------------------------------------------------------------*/

package bubble

import (
	"agent"
	"fmt"
)

// map 语法。
var pktMapClient = agent.HandlerMap{
	BZ_USERJOINREQ: JoinPktHandler,
	BZ_ROOMUSERCHG: RoomUserChgHandler,
}

//默认的客户端消息处理函数.
func ClientHandler(s *agent.Session, t int, b []byte) int {
	fmt.Printf("%s\n", "client default handler")
	return 0
}

//加入房间
func JoinPktHandler(s *agent.Session, t int, b []byte) int {
	_, join, _ := BzReadUserJoinReq(b)
	fmt.Printf("join pkt:%v\n", join)
	msg := Msg{t: MSG_T_JOIN, d: &(s.U.(*UserData).mq)}
	roomMgr.mq <- msg
	return 0
}

//玩家状态和位置变化.
func RoomUserChgHandler(s *agent.Session, t int, b []byte) int {
	_, chg, _ := BzReadRoomUserChg(b)
	fmt.Printf("chgx pkt:%v\n", chg)
	//msg := Msg{t: MSG_T_ROOM_USER_CHG, d: chg}
	return 0
}
