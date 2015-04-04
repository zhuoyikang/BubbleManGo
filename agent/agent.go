package agent

/*
事件类型
*/
type Handler func(*Session, int, []byte) int
type HandlerMap map[int]Handler

/*
用户必须实现自己的agent的回调。
*/
type Agent interface {
	Start(*Session)
	Stop(*Session)
	Handler(*Session, int, []byte) bool
}
