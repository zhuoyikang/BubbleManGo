//房间

package bubble

import (
	"fmt"
)

const (
	ROOM_MAP_WIDTH  = 18 //格子宽
	ROOM_MAP_HEIGHT = 10 //格子高.
)

//一个炸弹房间
type Room struct {
	mq chan Msg
	u1 chan Msg //玩家1
	u2 chan Msg //玩家2
}

func MakeRoom(u1 chan Msg, u2 chan Msg) *Room {
	mq := make(chan Msg, 2)
	fmt.Printf("u1 %v u2 %v\n", u1, u2)
	return &Room{u1: u1, u2: u2, mq: mq}
}

//房间建立时通知各个玩家.
func (r *Room) NotifyReady() {
	u1d := RoomReadMsg{id: 0, roomMq: r.mq}
	u2d := RoomReadMsg{id: 0, roomMq: r.mq}
	r.u1 <- Msg{t: MSG_T_ROOM_READY, d: u1d}
	r.u2 <- Msg{t: MSG_T_ROOM_READY, d: u2d}
}

// 这时候玩家channel有可能已经被关闭，需要捕获异常。
func (r *Room) notifyClose(mq chan Msg) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Notify close Panice %v\n", mq)
		}
	}()
	msg := Msg{t: MSG_T_ROOM_CLOSE}
	mq <- msg
}

// 房间关闭通知玩家。
func (r *Room) NotifyClose() {
	r.notifyClose(r.u1)
	r.notifyClose(r.u2)
}

//关闭房间
func (r *Room) Stop() {
	close(r.mq)
}

// 有玩家退出，告诉其他玩家，并自己退出
func (r *Room) MsgQuit(msg Msg) int {
	r.NotifyClose()
	r.Stop()
	return -1
}

// 分发消息
func (r *Room) dispatchMsg(msg Msg) int {
	switch {
	case msg.t == MSG_T_QUIT:
		return r.MsgQuit(msg)
	}
	return 0
}

// 监听各个channel，处理消息
// 需要设置心跳，如果超出时间没有任何消息，主动退出。
func (r *Room) Run() {
	var msg Msg
	var status bool

	for {
		msg, status = <-r.mq
		fmt.Printf("room mq %v %v\n", msg, status)
		if status == false {
			return
		}
		//返回0时退出.
		if r.dispatchMsg(msg) < 0 {
			return
		}
	}
}
