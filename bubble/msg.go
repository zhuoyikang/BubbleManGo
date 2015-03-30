package bubble

const (
	MSG_T_TCP_BIN = 1   // 消息转发.此类消息直接从进程中发送出去.\
	MSG_T_JOIN = 2      // 请求加入房间
	MSG_T_ROOM_CLOSE = 3 //房间关闭
	MSG_T_ROOM_READY = 4 //房间建立
)

// 消息.
type Msg struct {
	t int         //类型
	d interface{} //内容
}
