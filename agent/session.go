package agent

import (
	"fmt"
	"io"
	"net"
)

/*------------------------------------------------------------------------------
 session:单个连接由session处理。
------------------------------------------------------------------------------*/

type UD interface {
	Run()
	Stop()
}

type Session struct {
	conn net.Conn //网络连接,of course.
	U    UD       //玩家数据,of course.
}

func MakeSession(conn net.Conn) Session {
	s := Session{}
	s.conn = conn
	s.U = nil
	return s
}

func ShowBytes(s string, data []byte) {
	fmt.Printf("show byte(%d) %s:", len(data), s)
	for _, i := range data {
		fmt.Printf("%x ", i)
	}
	fmt.Printf("\n")
}

// 发数据包.
func (s *Session) SendBytes(data []byte) error {
	ShowBytes("SendBytes", data)
	w := s.conn
	want := len(data)
	n := 0
	for {
		ret, err := w.Write(data)
		if err != nil {
			return err
		}
		n += ret
		switch {
		case n == want:
			return nil
		case n > want:
			return nil
		}
	}
}

//读取一个完整的数据包.
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
func (s *Session) HandleClient(agent *TcpAgent) {
	agent.wg.Add(1)
	defer func() {
		s.conn.Close()
		agent.wg.Done()
		delete(agent.connectionPool, s.conn)
	}()
	for {
		t, data := s.ReadPkt(s.conn)
		fmt.Printf("r %d %v\n", t, data)
		if t < 0 {
			//读包异常，直接退出
			if s.U != nil {
				s.U.Stop()
			}
			break
		}
		//如果没有对应的处理函数，忽略,
		if h, err := agent.hmap[t]; err == false {
			fmt.Printf("%s\n", "unkown pkt")
			continue
		} else {
			h(s, data)
		}
	}
	fmt.Printf("%s\n", "handle client stop")
}
