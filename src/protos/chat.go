package protos

// 包括聊天,请求离线消息,通知

import (
	"errors"
	"fmt"
	"time"
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

	_s := ">BIIBQ4B"
	s := fmt.Sprint(_s, (len(msg) - zpack.CalcSize(_s)), "B")

	BIIBQ4BnB := zpack.Unpack(s, msg)

	// BIIBQ4BnB[0]: msgType
	to := BIIBQ4BnB[1]

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

func handle_offchat(user *types.User, msg []byte) (ack []byte, err error) {
	if !user.Coder.Shaked && !user.Logined {
		err = errors.New("not shaked or not logined")
		return
	}

	s := ">BI"
	BI := zpack.Unpack(s, msg)

	// BI[0]: msgType
	maxSize := BI[1]

	uid := user.UID
	// TODO: 通过uid获取离线消息data
	data := func(uid uint32, maxSize uint32) (data []byte) {
		return zpack.Pack('>', []interface{}{byte(3), uid, uint32(10002), byte(0), uint64(time.Now().Unix()), byte(0), byte(0), byte(0), byte(0), byte(98), byte(99), byte(100), byte('\xef'), byte('\xff'), uint32(10002), byte(0), uint64(time.Now().Unix()), byte(0), byte(0), byte(0), byte(0), byte(100), byte(101), byte(102)})
	}(uid, maxSize.(uint32))

	user.MQ <- data

	return
}
