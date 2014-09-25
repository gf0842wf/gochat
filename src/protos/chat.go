package protos

// 包括聊天,请求离线消息

import (
	"encoding/binary"
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
		err = errors.New("not shaked or not logined.")
		return
	}

	// 1(msgtype) + 4(to) + 4(from) + 1(online) + 4(gid) + 8(st) + n

	// msg[0]: msgType

	// to uid
	to := binary.BigEndian.Uint32(msg[1:5])
	// 把聊天发送时间改为服务器时间
	binary.BigEndian.PutUint64(msg[14:22], uint64(time.Now().UnixNano()/1000000))

	fmt.Println("from:", binary.BigEndian.Uint32(msg[5:9]), "to:", to)

	target := share.Clients.Get(to)
	if target == nil {
		fmt.Println("forward:", to)
		// TODO: 不在本服,发送到hub服务器,由它转发
	} else {
		targetUser := target.(*types.User)
		targetUser.MQ <- msg
	}

	return
}

func handle_getoffchat(user *types.User, msg []byte) (ack []byte, err error) {
	if !user.Coder.Shaked && !user.Logined {
		err = errors.New("not shaked or not logined.")
		return
	}

	// msg[0]: msgType

	// maxSize
	maxSize := binary.BigEndian.Uint32(msg[1:5])

	uid := user.UID
	// TODO: 通过uid获取离线消息data
	data := func(uid uint32, maxSize uint32) (data []byte) {
		return zpack.Pack('>', []interface{}{byte(3), uid, uint32(10002), byte(0), uint32(2), uint64(time.Now().UnixNano() / 1000000), byte(98), byte(99), byte(100), byte('\b'), byte('\r'), byte('\n'), uint32(10002), byte(0), uint32(2), uint64(time.Now().UnixNano() / 1000000), byte(100), byte(101), byte(102)})
	}(uid, maxSize)

	user.MQ <- data

	return
}
