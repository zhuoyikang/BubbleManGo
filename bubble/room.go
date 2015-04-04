//房间

package bubble

import (
	"fmt"
	"time"
)

const (
	ROOM_MAP_WIDTH  = 18 //格子宽
	ROOM_MAP_HEIGHT = 10 //格子高.
)

//一个炸弹房间
type Room struct {
	bubbleId   int32
	mm         [ROOM_MAP_HEIGHT][ROOM_MAP_WIDTH]int
	bubbleList []Bubble
	mq         chan Msg
	u1_id      int
	u2_id      int
	u1         chan Msg //玩家1
	u2         chan Msg //玩家2
}

// 新建一个房间
func MakeRoom(u1 chan Msg, u2 chan Msg) *Room {
	mq := make(chan Msg, 2)
	fmt.Printf("u1 %v u2 %v\n", u1, u2)
	room := Room{u1: u1, u2: u2, mq: mq}
	room.mm = [ROOM_MAP_HEIGHT][ROOM_MAP_WIDTH]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0},
		{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
		{0, 0, 1, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 1, 0, 0},
		{0, 0, 1, 3, 3, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 3, 3, 1, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 1, 0, 0},
		{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
		{0, 0, 0, 0, 0, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	room.bubbleList = make([]Bubble, 0)
	return &room
}

//房间建立时通知各个玩家.
func (r *Room) NotifyReady() {
	u1d := RoomReadMsg{id: 0, roomMq: r.mq}
	u2d := RoomReadMsg{id: 1, roomMq: r.mq}
	r.u1 <- Msg{t: MSG_T_ROOM_READY, d: u1d}
	r.u2 <- Msg{t: MSG_T_ROOM_READY, d: u2d}
	r.u1_id = 0
	r.u2_id = 1
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

// 转发消息给同一房间的其他玩家.
func (r *Room) MsgTcpBin(msg Msg) int {
	castMsg := msg.d.(RoomMsg)
	if castMsg.uid == r.u1_id {
		r.u2 <- msg
	} else {
		r.u1 <- msg
	}
	return 0
}

//处理客户端逻辑消息，并转发给别人.
func (r *Room) MsgSetBubble(msg Msg) int {
	castMsg := msg.d.(RoomMsg)
	_, setBubble, _ := BzReadSetBubble(castMsg.d)
	r.bubbleId += 1
	setBubble.b.id = r.bubbleId

	now := int32(time.Now().UnixNano() / 1000000000)
	setBubble.b.keeptime += now

	fmt.Printf("MsgSetBubble %v \n", setBubble.b)
	r.bubbleList = append(r.bubbleList, *setBubble.b)

	dbyte, _ := BzWriteSetBubble(make([]byte, 0), setBubble)
	castMsg.d = dbyte

	msg.d = castMsg
	//修改格式可让接收者转发给客户端。
	msg.t = MSG_T_TCP_BIN
	r.u2 <- msg
	r.u1 <- msg

	return 0
}

// 消息分发
func (r *Room) MsgRoom(msg Msg) int {
	fmt.Printf("MsgRoom dd %s\n", msg)
	castMsg := msg.d.(RoomMsg)
	switch {
	case castMsg.t == BZ_SETBUBBLEREQ:
		return r.MsgSetBubble(msg)
	}
	return 0
}

// 房间消息，各种逻辑

// 分发消息
func (r *Room) dispatchMsg(msg Msg) int {
	switch {
	case msg.t == MSG_T_QUIT:
		return r.MsgQuit(msg)
	case msg.t == MSG_T_TCP_BIN:
		return r.MsgTcpBin(msg)
	case msg.t == MSG_T_ROOM_MSG:
		return r.MsgRoom(msg)
	}
	return 0
}

// 泡泡爆炸逻辑
func (r *Room) BubbleBomb(b Bubble) {
	fmt.Printf("bb bomb %d\n", b.id)
	var bubbleBomb BubbleBomb
	bubbleBomb.id = b.id
	bubbleBomb.destroyTiles = make([]*BVector2,0)
	bubbleBomb.destroyUsers = make([]int32,0)

	bytes, _ := BzWriteBubbleBomb(make([]byte, 0), &bubbleBomb)
	var msg Msg
	msg.t = MSG_T_TCP_BIN
	msg.d = RoomMsg{t:9, uid:0, d:bytes}

	r.u2 <- msg
	r.u1 <- msg
}

//检查泡泡，设置爆炸.
func (r *Room) BubbleCheck() {
	var length int
	var bubble Bubble
	now := int32(time.Now().UnixNano() / 1000000000)

	for {
		length = len(r.bubbleList)
		if length == 0 {
			return
		}
		bubble = r.bubbleList[0]
		fmt.Printf("bubble:%d keeptime:%d now:%d\n", bubble.id, bubble.keeptime, now)

		if bubble.keeptime <= now {
			r.BubbleBomb(bubble)
		} else {
			//直接退出，没有更多可以爆的呢。
			return
		}
		r.bubbleList = r.bubbleList[1:]
	}
}

// 需要设置心跳，如果超出时间没有任何消息，主动退出。
func (r *Room) Run() {
	var msg Msg
	var status bool
	bubbleChkTimer := time.After(time.Second)

	for {
		select {
		case msg, status = <-r.mq:
			fmt.Printf("room mq %v %v\n", msg, status)
			if status == false {
				return
			}
			//返回0时退出.
			if r.dispatchMsg(msg) < 0 {
				return
			}
		case <-bubbleChkTimer:
			//定时检查泡泡爆炸
			r.BubbleCheck()
			bubbleChkTimer = time.After(time.Second)
		}

	}
}
