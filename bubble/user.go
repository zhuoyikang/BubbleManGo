package bubble

import (
	"agent"
	"fmt"
)

//
type UserData struct {
	S    *agent.Session
	mq   chan Msg
	udid string
	name string
}

// 用户登录成功，routine建立.
func MakeUserData(s *agent.Session, udid string) *UserData {
	mq := make(chan Msg)
	return &UserData{S: s, udid: udid, name: "zyk", mq: mq}
}

//玩家异步消息处理进程.
func (u *UserData) Run() {
	var msg Msg
	var status bool

	for {
		msg, status = <-u.mq
		//关闭status,退出即可.
		if status == false {
			break
		}

		fmt.Printf("u msg %v\n", msg)

		switch {
		case msg.t == MSG_T_TCP_BIN:
			bytes := msg.d.([]byte)
			u.S.SendBytes(bytes)
		case msg.t == MSG_T_ROOM_READY:
			ready := RoomReadyNtf{}
			bytes, _ := BzWriteRoomReadyNtf(make([]byte, 0), &ready)
			u.S.SendBytes(bytes)
		case msg.t == MSG_T_ROOM_CLOSE:
			close := RoomCloseNtf{}
			bytes, _ := BzWriteRoomCloseNtf(make([]byte, 0), &close)
			u.S.SendBytes(bytes)
		default:
			fmt.Printf("%s\n", "user unkown msg")
		}
	}
	fmt.Printf("%s\n", "going to die")
}

//玩家异步消息处理进程.
func (u *UserData) Stop() {
	//关闭mq后，用户服务器进程也会关闭.
	close(u.mq)
}
