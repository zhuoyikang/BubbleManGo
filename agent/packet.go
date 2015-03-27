package agent

// import (
// 	"fmt"
// 	"io"
// )

// type Packet struct {
// }

// //读取一个完整的数据包.
// // t:类型
// // d:数据
// func (pkt Packet) ReadPkt(reader io.Reader) (t int, buffer []byte) {
// 	//两个字节的类型.
// 	header := []byte{0, 0}

// 	//前两个字节包长度 >= 4
// 	n, err := io.ReadFull(reader, header)
// 	if err != nil || n != 2 {
// 		return -1, nil
// 	}
// 	pkt_length := uint16(header[0])<<8 | uint16(header[1])

// 	//后两个字节包类型
// 	n, err = io.ReadFull(reader, header)
// 	if err != nil || n != 2 {
// 		return -1, nil
// 	}

// 	t = int(uint16(header[0])<<8 | uint16(header[1]))

// 	//包可以无内容
// 	if pkt_length <= 4 {
// 		return t, nil
// 	}

// 	buffer = make([]byte, pkt_length-4)
// 	n, err = io.ReadFull(reader, buffer)

// 	if err != nil || uint16(n) != (pkt_length-4) {
// 		return -1, nil
// 	}

// 	//返回包内容和包内容.
// 	return t, buffer
// }

// func xx(s *Session, b []byte) (ret int) {
// 	fmt.Printf("%s\n", "haha")
// 	return
// }

// // map 语法。
// var packetHandlerMap = map[int]func(*Session, []byte) int{
// 	1: xx,
// }

// //在本routine中分发处理包.
// func (pkt Packet) DispatchPkt(t int, buffer []byte) {
// 	//如果没有对应的处理函数，忽略,
// 	if h, err := packetHandlerMap[t]; err == false {
// 		fmt.Printf("%s\n", "unkown pkt")
// 		return
// 	}
// 	h()
// }
