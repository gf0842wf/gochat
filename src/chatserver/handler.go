package main

// 转接network protocol 和 internal IPC 等

import (
	"errors"
	"fmt"
)

import (
	"protos"
	//"types"
)

// network protocol
func HandleNetProto(bot *Bot, data []byte) (ack []byte, err error) {
	msgType := data[0]
	fmt.Println("msgType:", msgType)
	if handle, ok := protos.NetProtoHandlers[msgType]; ok {
		ack, err = handle(bot.User, data) // 未去掉msgType
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
