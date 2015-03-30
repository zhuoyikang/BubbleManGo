package bubble

import (. "agent")

const (
	BZ_USERLOGINREQ = 1
	BZ_USERLOGINACK = 2
	BZ_USERJOINREQ = 3
	BZ_USERJOINACK = 4
	BZ_USERJOINNTF = 5
	BZ_ROOMCLOSENTF = 6
)

type UserLoginReq struct {
	udid string
}

type UserLoginAck struct {
	udid string
	name string
	level int32
}

type UserJoinReq struct {
	udid string
}

type UserJoinAck struct {
	udid string
	name string
	level int32
}

type RoomReadyNtf struct {
	t int32
}

type RoomCloseNtf struct {
	t int32
}

func BzReadUserLoginReq(datai []byte) (data []byte, ret *UserLoginReq, err error) {
	data = datai
	ret = &UserLoginReq{}
	data, ret.udid, err = BzReadstring(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteUserLoginReq(datai []byte, ret *UserLoginReq) (data []byte, err error) {
	data = datai
	data, err = BzWritestring(data, ret.udid)
	return
}
func BzReadUserLoginAck(datai []byte) (data []byte, ret *UserLoginAck, err error) {
	data = datai
	ret = &UserLoginAck{}
	data, ret.udid, err = BzReadstring(data)
 	if err != nil {
 		return
 	}
	data, ret.name, err = BzReadstring(data)
 	if err != nil {
 		return
 	}
	data, ret.level, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteUserLoginAck(datai []byte, ret *UserLoginAck) (data []byte, err error) {
	data = datai
	data, err = BzWritestring(data, ret.udid)
	data, err = BzWritestring(data, ret.name)
	data, err = BzWriteint32(data, ret.level)
	return
}
func BzReadUserJoinReq(datai []byte) (data []byte, ret *UserJoinReq, err error) {
	data = datai
	ret = &UserJoinReq{}
	data, ret.udid, err = BzReadstring(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteUserJoinReq(datai []byte, ret *UserJoinReq) (data []byte, err error) {
	data = datai
	data, err = BzWritestring(data, ret.udid)
	return
}
func BzReadUserJoinAck(datai []byte) (data []byte, ret *UserJoinAck, err error) {
	data = datai
	ret = &UserJoinAck{}
	data, ret.udid, err = BzReadstring(data)
 	if err != nil {
 		return
 	}
	data, ret.name, err = BzReadstring(data)
 	if err != nil {
 		return
 	}
	data, ret.level, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteUserJoinAck(datai []byte, ret *UserJoinAck) (data []byte, err error) {
	data = datai
	data, err = BzWritestring(data, ret.udid)
	data, err = BzWritestring(data, ret.name)
	data, err = BzWriteint32(data, ret.level)
	return
}
func BzReadRoomReadyNtf(datai []byte) (data []byte, ret *RoomReadyNtf, err error) {
	data = datai
	ret = &RoomReadyNtf{}
	data, ret.t, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteRoomReadyNtf(datai []byte, ret *RoomReadyNtf) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.t)
	return
}
func BzReadRoomCloseNtf(datai []byte) (data []byte, ret *RoomCloseNtf, err error) {
	data = datai
	ret = &RoomCloseNtf{}
	data, ret.t, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteRoomCloseNtf(datai []byte, ret *RoomCloseNtf) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.t)
	return
}
