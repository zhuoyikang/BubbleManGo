package agent

import (
	"fmt"
	"github.com/ugorji/go/codec"
	"io"
	"net"
)

/*------------------------------------------------------------------------------
 session:单个连接由session处理。
------------------------------------------------------------------------------*/

type Session struct {
	conn net.Conn  //网络连接,of course.
	U    *UserData //玩家数据,of course.
	H    codec.Handle
}

func MakeSession(conn net.Conn) Session {
	s := Session{}
	s.conn = conn
	s.H = new(codec.MsgpackHandle)
	return s
}

// 发包.
func (s *Session) SendPkt(t int, v interface{}) error {
	b := make([]byte, 64)
	enc := codec.NewEncoderBytes(&b, s.H)
	header := make([]byte, 2)
	if error := enc.Encode(v); error != nil {
		fmt.Printf("SendPkt %s\n", "encode2 error")
		return error
	}
	header[0] = byte(t >> 8)
	header[1] = byte(t)
	data := append(header, b...)
	return s.SendBytes(data)
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
		if t < 0 {
			//读包异常，直接退出
			break
		}
		//如果没有对应的处理函数，忽略,
		if h, err := agent.hmap[t]; err == false {
			fmt.Printf("%s\n", "unkown pkt")
			return
		} else {
			h(s, data)
		}
	}
}
