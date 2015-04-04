/*------------------------------------------------------------------------------
 房间管理，等待2个玩家加入后立刻新建房间.
------------------------------------------------------------------------------*/
package bubble


import (
	"fmt"
)



// 房间管理，可以建立和销毁一个房间.
type RoomMgr struct {
	mq chan Msg  // 消息处理
	w chan Msg // 正在等待的chan。
}

func MakeRoomMgr() *RoomMgr {
	//channel大小为1，同时只有1个等待者.
	return &RoomMgr{mq: make(chan Msg, 1)}
}

// 玩家加入消息
func (mgr *RoomMgr) MsgJoin(msg Msg) {
	if mgr.w == nil {
		mgr.w = msg.d.(chan Msg)
	} else {
		u1p := msg.d.(chan Msg)
		u2p := mgr.w
		mgr.w = nil
		room := MakeRoom(u1p, u2p)
		room.NotifyReady()
		go room.Run()
	}
}

// 玩家退出消息,这个消息只会由waitmq发过来。
func (mgr *RoomMgr) MsgQuit(msg Msg) {
	if mgr.w == nil {
		return
	}
	mgr.w = nil
}

// 分发消息
func (mgr *RoomMgr) dispatchMsg(msg Msg) {
	switch {
	case msg.t == MSG_T_JOIN:
		mgr.MsgJoin(msg)
	case msg.t == MSG_T_QUIT:
		mgr.MsgQuit(msg)
	default:
		fmt.Printf("%s\n", "x 1")
	}
}

// 等待并返回消息.
func (mgr *RoomMgr) waitMsg() (msg Msg, status bool) {
	if mgr.w == nil {
		msg, status = <-mgr.mq
		return
	} else {
		select {
		case msg, status = <-mgr.w:
			//老子退出了，转为一个客户端退出消息.
			if status == false {
				msg = Msg{t: MSG_T_QUIT, d: mgr.w}
				status = true
			}
			//新人来了
		case msg, status = <-mgr.mq:
			return
		}
	}
	return
}

// A玩家加入，放到等待
// B玩家加入，将其和A分配房间.
func (mgr *RoomMgr) Run() {
	for {
		msg, status := mgr.waitMsg()
		if status == false {
			break
		}
		mgr.dispatchMsg(msg)
	}
	fmt.Printf("%s\n", "room mgr run quit")
}

//关闭mgr.mq， roomMgr的routine会自动退出
func (mgr *RoomMgr) Stop() {
	close(mgr.mq)
}

//全局的
var roomMgr *RoomMgr

//
func init() {
	roomMgr = MakeRoomMgr()
	go roomMgr.Run()
	fmt.Printf("%s\n", "init")
}
