package bubble

import (
	"agent"
	"github.com/ugorji/go/codec"
	"fmt"
)

// map 语法。
var PacketHandlerMap = agent.HandlerMap{
	1: LoginPktHandler,
}

//登陆包
func LoginPktHandler(s *agent.Session, b []byte) int {
	var login LoginPkt
	dec := codec.NewDecoderBytes(b, s.H)
	dec.Decode(&login)

	fmt.Printf("%s\n", "login")
	return 0;
}
