package protos

// 包括握手,心跳,登录

import (
	"errors"
	"fmt"
)

import (
	"misc/zpack"
	"share"
	"types"
)

func handle_shake(user *types.User, msg []byte) (ack []byte, err error) {
	subType := msg[0]
	fmt.Println(msg)
	if subType == 0 {
		ack = zpack.Pack('>', []interface{}{byte(0), byte(1), user.Coder.CryptKey})
	} else if subType == 2 {
		user.Coder.Shaked = true
		fmt.Println("Shaked:")
	} else {
		err = errors.New("handle_shake: unknown subType")
	}

	return
}

func handle_nop(user *types.User, msg []byte) (ack []byte, err error) {
	// TODO: 发送到hub服务器,维持用户在线信息
	return
}

func handle_login(user *types.User, msg []byte) (ack []byte, err error) {
	if !user.Coder.Shaked {
		err = errors.New("handle_login: not shaked")
		return
	}

	s := fmt.Sprint(">I", len(msg)-4, "B")
	InB := zpack.Unpack(s, msg)
	uid, password := InB[0], InB[1]

	if true { // TODO: 向hub服务器发送用户名密码请求登录(http接口?)
		fmt.Println("Login:", uid.(uint32), string(password.([]byte)))
		user.UID = uid.(uint32)
		user.Password = string(password.([]byte))
		user.Online = true
		user.Logined = true
		share.Clients.Set(user.UID, user)
		// TODO: 通知hub服务器该用户登录
		ack = zpack.Pack('>', []interface{}{byte(2), byte(0)})
	} else {
		err = errors.New("login failed")
		//ack = zpack.Pack('>', []interface{}{byte(2), byte(1)})
	}

	return
}
