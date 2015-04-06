package bubble

import (. "agent")

const (
	BZ_USERLOGINREQ = 1
	BZ_USERLOGINACK = 2
	BZ_USERJOINREQ = 3
	BZ_USERJOINACK = 4
	BZ_ROOMREADYNTF = 5
	BZ_ROOMCLOSENTF = 6
	BZ_ROOMUSERCHG = 7
	BZ_SETBUBBLEREQ = 8
	BZ_BUBBLEBOMB = 9
	BZ_ROOMUSERSTATUSCHG = 10
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

type BVector2 struct {
	x int32
	y int32
}

type RoomUser struct {
	pos *BVector2
	direction int32
	status int32
}

type RoomReadyNtf struct {
	roomId int32
	uPosAll []*RoomUser
	uIdx int32
}

type RoomCloseNtf struct {
	t int32
}

type RoomUserChg struct {
	uIdx int32
	user *RoomUser
}

type Bubble struct {
	id int32
	pos *BVector2
	power int32
	keeptime int32
}

type SetBubble struct {
	b *Bubble
	uIdx int32
}

type BubbleBomb struct {
	id int32
	destroyTiles []*BVector2
	destroyUsers []int32
}

type RoomUserStatusChg struct {
	id int32
	status int32
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
func BzReadBVector2(datai []byte) (data []byte, ret *BVector2, err error) {
	data = datai
	ret = &BVector2{}
	data, ret.x, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.y, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteBVector2(datai []byte, ret *BVector2) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.x)
	data, err = BzWriteint32(data, ret.y)
	return
}
func BzReadRoomUser(datai []byte) (data []byte, ret *RoomUser, err error) {
	data = datai
	ret = &RoomUser{}
	data, ret.pos, err = BzReadBVector2(data)
 	if err != nil {
 		return
 	}
	data, ret.direction, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.status, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteRoomUser(datai []byte, ret *RoomUser) (data []byte, err error) {
	data = datai
	data, err = BzWriteBVector2(data, ret.pos)
	data, err = BzWriteint32(data, ret.direction)
	data, err = BzWriteint32(data, ret.status)
	return
}
func BzReadRoomReadyNtf(datai []byte) (data []byte, ret *RoomReadyNtf, err error) {
	data = datai
	ret = &RoomReadyNtf{}
	data, ret.roomId, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	var uPosAll_v *RoomUser
	data, uPosAll_size, err := BzReaduint16(data)
	for i := 0; i < int(uPosAll_size); i++ {
		data, uPosAll_v, err = BzReadRoomUser(data)
	 	if err != nil {
	 		return
	 	}
		ret.uPosAll = append(ret.uPosAll, uPosAll_v)
	}
 	if err != nil {
 		return
 	}
	data, ret.uIdx, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteRoomReadyNtf(datai []byte, ret *RoomReadyNtf) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.roomId)
	data, err = BzWriteuint16(data, uint16(len(ret.uPosAll)))
	for _, uPosAll_v := range ret.uPosAll {
		data, err = BzWriteRoomUser(data, uPosAll_v)
	}
	data, err = BzWriteint32(data, ret.uIdx)
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
func BzReadRoomUserChg(datai []byte) (data []byte, ret *RoomUserChg, err error) {
	data = datai
	ret = &RoomUserChg{}
	data, ret.uIdx, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.user, err = BzReadRoomUser(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteRoomUserChg(datai []byte, ret *RoomUserChg) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.uIdx)
	data, err = BzWriteRoomUser(data, ret.user)
	return
}
func BzReadBubble(datai []byte) (data []byte, ret *Bubble, err error) {
	data = datai
	ret = &Bubble{}
	data, ret.id, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.pos, err = BzReadBVector2(data)
 	if err != nil {
 		return
 	}
	data, ret.power, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.keeptime, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteBubble(datai []byte, ret *Bubble) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.id)
	data, err = BzWriteBVector2(data, ret.pos)
	data, err = BzWriteint32(data, ret.power)
	data, err = BzWriteint32(data, ret.keeptime)
	return
}
func BzReadSetBubble(datai []byte) (data []byte, ret *SetBubble, err error) {
	data = datai
	ret = &SetBubble{}
	data, ret.b, err = BzReadBubble(data)
 	if err != nil {
 		return
 	}
	data, ret.uIdx, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteSetBubble(datai []byte, ret *SetBubble) (data []byte, err error) {
	data = datai
	data, err = BzWriteBubble(data, ret.b)
	data, err = BzWriteint32(data, ret.uIdx)
	return
}
func BzReadBubbleBomb(datai []byte) (data []byte, ret *BubbleBomb, err error) {
	data = datai
	ret = &BubbleBomb{}
	data, ret.id, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	var destroyTiles_v *BVector2
	data, destroyTiles_size, err := BzReaduint16(data)
	for i := 0; i < int(destroyTiles_size); i++ {
		data, destroyTiles_v, err = BzReadBVector2(data)
	 	if err != nil {
	 		return
	 	}
		ret.destroyTiles = append(ret.destroyTiles, destroyTiles_v)
	}
 	if err != nil {
 		return
 	}
	var destroyUsers_v int32
	data, destroyUsers_size, err := BzReaduint16(data)
	for i := 0; i < int(destroyUsers_size); i++ {
		data, destroyUsers_v, err = BzReadint32(data)
	 	if err != nil {
	 		return
	 	}
		ret.destroyUsers = append(ret.destroyUsers, destroyUsers_v)
	}
 	if err != nil {
 		return
 	}
	return
}
func BzWriteBubbleBomb(datai []byte, ret *BubbleBomb) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.id)
	data, err = BzWriteuint16(data, uint16(len(ret.destroyTiles)))
	for _, destroyTiles_v := range ret.destroyTiles {
		data, err = BzWriteBVector2(data, destroyTiles_v)
	}
	data, err = BzWriteuint16(data, uint16(len(ret.destroyUsers)))
	for _, destroyUsers_v := range ret.destroyUsers {
		data, err = BzWriteint32(data, destroyUsers_v)
	}
	return
}
func BzReadRoomUserStatusChg(datai []byte) (data []byte, ret *RoomUserStatusChg, err error) {
	data = datai
	ret = &RoomUserStatusChg{}
	data, ret.id, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.status, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteRoomUserStatusChg(datai []byte, ret *RoomUserStatusChg) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.id)
	data, err = BzWriteint32(data, ret.status)
	return
}
