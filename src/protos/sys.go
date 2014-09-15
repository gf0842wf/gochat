package protos

// 包括握手,心跳,登录

import (
	"fmt"
)

import (
	"bytes"
	"encoding/binary"
	"errors"
	"types"
)

func handle_shake(user *types.User, msg []byte) (ack []byte, err error) {
	subType := msg[0]
	buf := new(bytes.Buffer)
	if subType == 0 {
		binary.Write(buf, binary.BigEndian, byte(1))
		binary.Write(buf, binary.BigEndian, user.Coder.CryptKey)
		ack = buf.Bytes()
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

	//uid, ok := (*(pObj.MsgJson))["uid"]
	//if !ok {
	//	err = errors.New("no uid field")
	//	return
	//}
	//password, ok := (*(pObj.MsgJson))["password"]
	//if !ok {
	//	err = errors.New("no password field")
	//	return
	//}

	//if true { // TODO: 向hub服务器发送用户名密码请求登录(http接口?)
	//	fmt.Println("Login:", uid.(uint32), password.(string))
	//	fmt.Println("Logined")
	//	user.UID = uid.(uint32)
	//	user.Password = password.(string)
	//	user.Online = true
	//	user.Logined = true
	//	// TODO: 通知hub服务器该用户登录
	//} else {
	//	err = errors.New("login failed")
	//}

	return
}
