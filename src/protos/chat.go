package protos

// 包括聊天,请求离线消息,通知

import (
	"errors"
	"fmt"
)

import (
	"misc/zpack"
	"share"
	"types"
)

func handle_chat(user *types.User, msg []byte) (ack []byte, err error) {
	if !user.Coder.Shaked && !user.Logined {
		err = errors.New("not shaked or not logined")
		return
	}

	_s := ">IIBQ4B"
	s := fmt.Sprint(_s, (len(msg) - zpack.CalcSize(_s)), "B")

	IIBQ4BnB := zpack.Unpack(s, msg)
	to := IIBQ4BnB[0]

	target := share.Clients.Get(to.(uint32))
	if target == nil {
		fmt.Println("forward:", to)
		// TODO: 不在本服,发送到hub服务器,由它转发
	} else {
		targetUser := target.(*types.User)
		targetUser.MQ <- msg
	}
	return
}
