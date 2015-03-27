package agent

import (
	"fmt"
	"io"
	"net"
)

/*------------------------------------------------------------------------------
 session:单个连接由session处理。
------------------------------------------------------------------------------*/

type Session struct {
	conn net.Conn  //网络连接,of course.
	u    *UserData //玩家数据,of course.
}

//读取一个完整的数据包.
// t:类型
// d:数据
func (s *Session) ReadPkt(reader io.Reader) (t int, buffer []byte) {
	//两个字节的类型.
	header := []byte{0, 0}

	//前两个字节包长度 >= 4
	n, err := io.ReadFull(reader, header)
	if err != nil || n != 2 {
		return -1, nil
	}
	pkt_length := uint16(header[0])<<8 | uint16(header[1])

	//后两个字节包类型
	n, err = io.ReadFull(reader, header)
	if err != nil || n != 2 {
		return -1, nil
	}

	t = int(uint16(header[0])<<8 | uint16(header[1]))

	//包可以无内容
	if pkt_length <= 4 {
		return t, nil
	}

	buffer = make([]byte, pkt_length-4)
	n, err = io.ReadFull(reader, buffer)

	if err != nil || uint16(n) != (pkt_length-4) {
		return -1, nil
	}

	//返回包内容和包内容.
	return t, buffer
}

// 处理客户端连接
func (session *Session) HandleClient(agent *TcpAgent) {
	agent.wg.Add(1)
	fmt.Printf("%s %v\n", "begin", session.conn)
	defer func() {
		session.conn.Close()
		agent.wg.Done()
		delete(agent.connectionPool, session.conn)
		fmt.Printf("%s %v\n", "end", session.conn)
	}()

	for {
		t, data := session.ReadPkt(session.conn)
		if t < 0 {
			//读包异常，直接退出
			break
		}
		//如果没有对应的处理函数，忽略,
		if h, err := agent.hmap[t]; err == false {
			fmt.Printf("%s\n", "unkown pkt")
			return
		} else {
			h(session.u, data)
		}
	}
}
