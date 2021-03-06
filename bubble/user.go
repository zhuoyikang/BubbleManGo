package bubble

import (
	"agent"
	"fmt"
	"sync"
)

// 玩家数据
type UserData struct {
	S  *agent.Session
	mq chan Msg
	//tailMq     chan Msg
	roomMq chan Msg
	udid   string
	name   string
	uid    int
	mutex  *sync.Mutex
}

// 用户登录成功，routine建立.
func MakeUserData(s *agent.Session, udid string) *UserData {
	mq := make(chan Msg)
	//tailMq := make(chan Msg)
	return &UserData{S: s, udid: udid, name: "zyk",
		mutex: &sync.Mutex{},
		mq:    mq, // tailMq:tailMq
	}
}

// 处理客户端消息
func (u *UserData) MsgClient(msg Msg) int {
	m := msg.d.(ClientMsg)
	if h, status := pktMapClient[m.t]; status == true {
		return h(u.S, m.t, m.d)
	}else{
		return ClientHandler(u.S, m.t, m.d)
	}
}

//分发处理各项消息，better use map。
func (u *UserData) MsgDispatch(msg Msg, status bool) int {
	//关闭status,退出即可.
	if status == false {
		return -1
	}

	switch {
	case msg.t == MSG_T_CLIENT:
		return u.MsgClient(msg)
	case msg.t == MSG_T_TCP_BIN:
		return u.MsgTcpBin(msg)
	case msg.t == MSG_T_ROOM_READY:
		return u.MsgRoomReady(msg)
	case msg.t == MSG_T_ROOM_CLOSE:
		return u.MsgRoomClose(msg)
	default:
		fmt.Printf("%s\n", "user unkown msg")
	}
	return 0
}

//玩家异步消息处理进程.
func (u *UserData) Run() {
	var msg Msg
	var status bool
Out:
	for {
		msg, status = <-u.mq
		//u.mutex.Lock()
		if ret := u.MsgDispatch(msg, status); ret < 0 {
			//u.mutex.Unlock()
			break Out
		}
		//u.mutex.Unlock()
	}
	fmt.Printf("%s\n", "going to die")
}

//玩家异步消息处理进程.
func (u *UserData) Stop() {
	fmt.Printf("%s stop %v\n", "UseData", u.roomMq)
	if u.roomMq != nil {
		u.roomMq <- Msg{t: MSG_T_QUIT}
	}
	//关闭mq后，用户服务器进程也会关闭.
	close(u.mq)
}

/*------------------------------------------------------------------------------
 各种业务消息处理
------------------------------------------------------------------------------*/

// 直接发送2进制消息，
func (u *UserData) MsgTcpBin(msg Msg) int {
	castMsg := msg.d.(RoomMsg)
	u.S.SendPkt(uint16(castMsg.t), castMsg.d)
	return 0
}

// 房间准备好。
func (u *UserData) MsgRoomReady(msg Msg) int {
	u1 := &RoomUser{
		pos: &BVector2{x: 1,
			y: 1,
		},
		direction: 5,
		status:    0,
	}
	u2 := &RoomUser{
		pos: &BVector2{x: ROOM_MAP_WIDTH - 2,
			y: 1},
		direction: 5,
		status:    0,
	}
	uAll := []*RoomUser{u1, u2}
	roomInfo := msg.d.(RoomReadMsg)
	ready := RoomReadyNtf{
		roomId:  0,
		uPosAll: uAll,
		uIdx:    int32(roomInfo.id),
	}
	bytes, _ := BzWriteRoomReadyNtf(make([]byte, 0), &ready)
	u.S.SendPkt(BZ_ROOMREADYNTF, bytes)
	u.roomMq = roomInfo.roomMq
	u.uid = roomInfo.id
	return 0
}

// 房间关门闭。
func (u *UserData) MsgRoomClose(msg Msg) int {
	fmt.Printf("%s\n", "use roomr close")

	close := RoomCloseNtf{}
	bytes, _ := BzWriteRoomCloseNtf(make([]byte, 0), &close)
	u.S.SendPkt(BZ_ROOMCLOSENTF, bytes)

	//设置为nil
	u.roomMq = nil
	fmt.Printf("MsgRoomClose %v\n", u.roomMq)
	return 0
}
