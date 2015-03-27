package bubble

import (
	"agent"
	"fmt"
)

// map 语法。
var PacketHandlerMap = agent.HandlerMap{
	1: LoginPktHandler,
}

//登陆包
func LoginPktHandler(s *agent.UserData, b []byte) int {
	fmt.Printf("%s\n", "login")
	return 0;
}
