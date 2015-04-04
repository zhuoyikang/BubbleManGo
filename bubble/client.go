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
	BZ_ROOMUSERCHG: RoomCastPktHandler,
	BZ_SETBUBBLEREQ: MsgToRoomPktHandler,
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
	msg := Msg{t: MSG_T_JOIN, d: (s.U.(*UserData).mq)}
	roomMgr.mq <- msg
	return 0
}

// 将消息转发给room,room会转发房间里的其他人|不处理逻辑。
func RoomCastPktHandler(s *agent.Session, t int, b []byte) int {
	u := s.U.(*UserData)
	rMsg := RoomMsg{t: t, d: b, uid:u.uid}
	fmt.Printf("cast %d to room %v\n", t, rMsg)
	u.roomMq <- Msg{t: MSG_T_TCP_BIN, d: rMsg}
	return 0
}

// 将消息转发给room.room会处理逻辑.
func MsgToRoomPktHandler(s *agent.Session, t int, b []byte) int {
	u := s.U.(*UserData)
	rMsg := RoomMsg{t: t, d: b, uid:u.uid}
	fmt.Printf("transfer %d to room %v\n", t, rMsg)
	u.roomMq <- Msg{t: MSG_T_ROOM_MSG, d: rMsg}
	return 0
}
