package main

// 转接network protocol 和 internal IPC 等

import (
	"errors"
	"fmt"
)

import (
	"protos"
)

// network protocol
func HandleNetProto(bot *Bot, data []byte) (ack []byte, err error) {
	msgType := data[0]
	fmt.Println("UID:", bot.User.UID, "msgType:", msgType)
	if handle, ok := protos.NetProtoHandlers[msgType]; ok {
		ack, err = handle(bot.User, data)
	} else {
		err = errors.New("unknown msgType")
	}

	return
}

// forward chat message protocol
func HandleForwardProto(bot *Bot, data []byte) (ack []byte, err error) {
	return
}

// offline chat message protocol
func HandleOffchatProto(bot *Bot, data []byte) (ack []byte, err error) {
	return
}

// timer message protocol
func HandleTickProto(bot *Bot, data []byte) (ack []byte, err error) {
	return
}
