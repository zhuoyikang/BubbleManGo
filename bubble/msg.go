package bubble

const (
	MSG_T_CLIENT = 0    //客户端消息
	MSG_T_TCP_BIN = 1   // 消息转发.此类消息直接从进程中发送出去.\
	MSG_T_JOIN = 2      // 请求加入房间
	MSG_T_ROOM_CLOSE = 3 //房间关闭
	MSG_T_ROOM_READY = 4 //房间建立
	MSG_T_ROOM_USER_CHG = 4
	MSG_T_QUIT = 3      // 请求退出房间
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
	id int
	roomMq chan Msg
}

// RoomCast Msg
type RoomCastMsg struct {
	t int
	uid int
	d []byte
}
