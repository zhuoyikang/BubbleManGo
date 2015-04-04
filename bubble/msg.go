package bubble

const (
	MSG_T_CLIENT     = iota //客户端消息
	MSG_T_TCP_BIN           // 消息转发.此类消息直接从进程中发送出去.\
	MSG_T_JOIN              // 请求加入房间
	MSG_T_ROOM_CLOSE        //房间关闭
	MSG_T_ROOM_READY        //房间建立
	MSG_T_ROOM_USER_CHG
	MSG_T_QUIT     // 请求退出房间
	MSG_T_ROOM_MSG //房间消息.
)

// 消息.
type Msg struct {
	t int         //类型
	d interface{} //内容
}

//客户端消息
type ClientMsg struct {
	t int
	d []byte
}

// 房间准备好了
type RoomReadMsg struct {
	id     int
	roomMq chan Msg
}

// RoomCast Msg
type RoomMsg struct {
	t   int
	uid int
	d   []byte
}

// type MsgHandler func() int
// type MsgHandlerMap map[int]Handler

