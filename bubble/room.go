//房间

package bubble

import (
	"fmt"
)

const (
	ROOM_MAP_WIDTH = 18   //格子宽
	ROOM_MAP_HEIGHT = 10  //格子高.
)


//一个炸弹房间
type Room struct {
	mq chan Msg
	u1 chan Msg //玩家1
	u2 chan Msg //玩家2
}

func MakeRoom(u1 chan Msg, u2 chan Msg) *Room {
	mq := make(chan Msg, 2)
	return &Room{u1: u1, u2: u2, mq: mq}
}

//房间建立时通知各个玩家.
func (r *Room) NotifyReady() {
	r.u1 <- Msg{t: MSG_T_ROOM_READY, d:int32(0)}
	r.u2 <- Msg{t: MSG_T_ROOM_READY, d:int32(1)}
}

//其中一个玩家退出，给另一个玩家发送通知消息，并关闭room。
func (r *Room) NotifyClose(u chan Msg) {
	msg := Msg{t: MSG_T_ROOM_CLOSE}
	u <- msg
}

//关闭房间
func (r *Room) Stop() {
	close(r.mq)
}

//监听各个channel，处理消息
func (r *Room) Run() {
	defer close(r.mq)

	var msg Msg
	var status bool

J:
	for {
		select {
		case msg, status = <-r.mq:
			fmt.Printf("mq %v %v\n", msg, status)
			if status == false {
				//房间退出，给两个玩家发送通知消息.
				r.NotifyClose(r.u1)
				r.NotifyClose(r.u2)
				break J
			}
		case msg, status = <-r.u1:
			if status == false {
				r.NotifyClose(r.u2)
				break J
			}
			fmt.Printf("u1 %v %v\n", msg, status)
		case msg, status = <-r.u2:
			if status == false {
				r.NotifyClose(r.u1)
				break J
			}
			fmt.Printf("u2 %v %v\n", msg, status)
		}
		fmt.Printf("%s\n", "hehere")
	}
}

// 房间管理，可以建立和销毁一个房间.
type RoomMgr struct {
	mq chan Msg  // 消息处理
	w  *chan Msg // 正在等待的chan。
}

func MakeRoomMgr() *RoomMgr {
	return &RoomMgr{mq: make(chan Msg, 2)}
}

// 执行.
func (mgr *RoomMgr) Run() {
	var msg Msg
	var status bool
	for {
		if mgr.w == nil {
			msg, status = <-mgr.mq
			//退出
			if status == false {
				break
			}
			switch {
			case msg.t == MSG_T_JOIN:
				mgr.w = msg.d.(*chan Msg)
			default:
				fmt.Printf("%s\n", "x 1")
			}
		} else {
			select {
			case msg, status = <-*mgr.w:
				//老子退出了
				if status == false {
					mgr.w = nil
				} else {
					fmt.Printf("%s\n", "shoud bb mgr run")
				}
				//新人来了
			case msg, status = <-mgr.mq:
				switch {
				case msg.t == MSG_T_JOIN:
					u1p := msg.d.(*chan Msg)
					u2p := mgr.w
					mgr.w = nil
					room := MakeRoom(*u1p, *u2p)
					room.NotifyReady()
					go room.Run()
				default:
					fmt.Printf("%s\n", "x 2")
				}

			}
		}
	}
}

var roomMgr *RoomMgr

//
func init() {
	roomMgr = MakeRoomMgr()
	go roomMgr.Run()
	fmt.Printf("%s\n", "init")
}
