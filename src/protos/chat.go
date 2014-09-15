package protos

// 包括聊天,请求离线消息,通知

import (
	"fmt"
)

import (
	"errors"
	"share"
	"types"
)

func handle_chat(user *types.User, pObj *types.Msg) (ack []byte, err error) {
	if !user.Coder.Shaked && !user.Logined {
		err = errors.New("not shaked or not logined")
		return
	}

	//line, ok := (*(pObj.MsgJson))["line"]
	//if !ok {
	//	err = errors.New("no line field")
	//	return
	//}
	//from, ok := (*(pObj.MsgJson))["from"]
	//if !ok {
	//	err = errors.New("no from field")
	//	return
	//}
	to, ok := (*(pObj.MsgJson))["to"]
	if !ok {
		err = errors.New("no to field")
		return
	}
	//st, ok := (*(pObj.MsgJson))["st"]
	//if !ok {
	//	err = errors.New("no st field")
	//	return
	//}
	//flg1, ok := (*(pObj.MsgJson))["flg1"]
	//if !ok {
	//	//err = errors.New("no flg1 field")
	//	//return
	//}
	//flg2, ok := (*(pObj.MsgJson))["flg2"]
	//if !ok {
	//	//err = errors.New("no flg2 field")
	//	//return
	//}
	//flg3, ok := (*(pObj.MsgJson))["flg3"]
	//if !ok {
	//	//err = errors.New("no flg3 field")
	//	//return
	//}
	//ctx, ok := (*(pObj.MsgJson))["ctx"]
	//if !ok {
	//	err = errors.New("no ctx field")
	//	return
	//}

	target := share.Clients.Get(to.(uint32))
	if target == nil {
		// TODO: 不在本服,发送到hub服务器,由它转发
	}
	targetUser := target.(types.User)
	targetUser.MQ <- []byte{1, 2}
	fmt.Println(targetUser)

	return
}
