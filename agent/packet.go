package agent

import (
	"errors"
)

func BzReadbyte(datai []byte) (data []byte, ret byte, err error) {
	data = datai
	if 1 > len(data) {
		err = errors.New("read byte failed")
		return
	}
	ret = data[0]
	data = data[1:]
	return
}

func BzWritebyte(datai []byte, v byte) (data []byte, err error) {
	data = datai
	data = append(data, byte(v))
	return
}

func BzReaduint16(datai []byte) (data []byte, ret uint16, err error) {
	data = datai
	if 2 > len(data) {
		err = errors.New("read uint16 failed")
		return
	}

	buf := data[0:2]
	ret = uint16(buf[0])<<8 | uint16(buf[1])
	data = data[2:]
	return
}

func BzWriteuint16(datai []byte, v uint16) (data []byte, err error) {
	data = datai
	data = append(data, byte(v>>8), byte(v))
	return
}

func BzReadint16(datai []byte) (data []byte, ret int16, err error) {
	if 2 > len(data) {
		err = errors.New("read uint16 failed")
		return
	}

	buf := data[0:2]
	ret = int16(buf[0])<<8 | int16(buf[1])
	data = data[2:]
	return
}

func BzWriteint16(datai []byte, v int16) (data []byte, err error) {
	data = datai
	data = append(data, byte(v>>8), byte(v))
	return
}

func BzReaduint32(datai []byte) (data []byte, ret uint32, err error) {
	data = datai
	if 4 > len(data) {
		err = errors.New("read uint32 failed")
		return
	}

	buf := data[0:4]
	ret = uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 |
		uint32(buf[3])

	data = data[4:]
	return
}

func BzWriteuint32(datai []byte, v uint32) (data []byte, err error) {
	data = datai
	data = append(data, byte(v>>24), byte(v>>16),
		byte(v>>8), byte(v))
	return
}

func BzReadint32(datai []byte) (data []byte, ret int32, err error) {
	data, ret1, err := BzReaduint32(datai)
	ret = int32(ret1)
	return
}

func BzWriteint32(datai []byte, v int32) (data []byte, err error) {
	return BzWriteuint32(datai, uint32(v))
}

func BzReadstring(datai []byte) (data []byte, ret string, err error) {
	data, size, err := BzReaduint16(datai)
	if err != nil {
		return
	}
	if int(size) > len(data) {
		err = errors.New("read string failed")
	}

	bytes := data[0:int(size)]
	ret = string(bytes)
	data = data[int(size):]
	return
}

func BzWritestring(datai []byte, str string) (data []byte, err error) {
	bytes := []byte(str)
	data, err = BzWriteuint16(datai, uint16(len(bytes)))
	data = append(data, bytes...)
	return
}

// // 创建一个完整的数据包
// func MakePacketData(api uint16, datai []byte) (data []byte) {
// 	length := 4 + len(datai)
// 	data = append(data, byte(length>>8), byte(length))
// 	data = append(data, byte(api>>8), byte(api))
// 	data = append(data, datai...)

// 	return
// }
