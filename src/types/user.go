package types

import (
	"misc/codec"
)

type User struct {
	UID      int32  // 用户id
	SID      string // 用户所在服
	Nickname string // 昵称
	Password string // 密码

	IsActive bool // 是否在线

	Mac    string // 玩家MAC地址
	OSType int64  // 系统类型

	Coder *codec.Coder

	MQ chan IPCObj // User之间通信队列
}

func NewUser() *User {
	mq = make(chan IPCObj)
	sid = "0" // TODO: SID通过配置获取
	return &User{IsActive: true, MQ: mq, SID: sid}
}
