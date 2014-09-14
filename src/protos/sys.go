package protos

// 包括握手,心跳,登录

import (
	"fmt"
)

import (
	"errors"
	"types"
)

func handle_shake(user *types.User, pObj *types.Msg) (ack []byte, err error) {
	subType, ok := (*(pObj.MsgJson))["type"]
	if !ok {
		err = errors.New("no type field")
		return
	}
	pAckObj := types.NewSendMsg(0)
	if subType.(string) == "PRE" {
		(*(pAckObj.MsgJson))["type"] = "REQ"
		(*(pAckObj.MsgJson))["key"] = user.Coder.CryptKey
		ack, err = types.Marshal(pAckObj)
		user.Coder.Encode(ack)
		return
	} else if subType.(string) == "ACK" {
		user.Coder.Shaked = true
		fmt.Println("Shaked:", user)
		return
	}

	return
}

func handle_login(user *types.User, pObj *types.Msg) (ack []byte, err error) {
	if !user.Coder.Shaked {
		err = errors.New("not shaked")
		return
	}

	uid, ok := (*(pObj.MsgJson))["uid"]
	if !ok {
		err = errors.New("no uid field")
		return
	}
	password, ok := (*(pObj.MsgJson))["password"]
	if !ok {
		err = errors.New("no password field")
		return
	}

	if true { // TODO: 向hub服务器发送用户名密码请求登录(http接口?)
		fmt.Println("Login:", uid.(uint32), password.(string))
		fmt.Println("Logined")
		user.UID = uid.(uint32)
		user.Password = password.(string)
		user.Online = true
		user.Logined = true
		// TODO: 通知hub服务器该用户登录
	} else {
		err = errors.New("login failed")
	}

	return
}
