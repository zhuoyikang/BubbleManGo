package hub

import (. "agent")

const (
	BZ_USERLOGINREQ = 1
)

type UserLoginReq struct {
	udid string
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
