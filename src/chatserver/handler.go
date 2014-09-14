package main

// 转接network protocol 和 internal IPC 等

import (
	"errors"
	"fmt"
)

import (
	"protos"
	"types"
)

// network protocol
func HandleNetProto(bot *Bot, data []byte) (ack []byte, err error) {
	var obj types.Msg
	bot.User.Coder.Decode(data)
	err, _ = types.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("HandleNetProto decode err:", err.Error())
		return
	}

	msgType := obj.MsgType
	fmt.Println("msgType:", msgType)
	if handle, ok := protos.NetProtoHandlers[msgType]; ok {
		ack, err = handle(bot.User, &obj)
	} else {
		err = errors.New("unknown msgType")
	}

	return
}

// internal IPC
func HandleIPCProto(bot *Bot, data []byte) (ack []byte, err error) {
	return
}

// 定时消息
func HandleTMProto(bot *Bot, data []byte) (ack []byte, err error) {
	return
}
