package types

import (
//"protos"
)

type User struct {
	UID      uint32 // 用户id
	SID      string // 用户所在服
	Nickname string // 昵称
	Password string // 密码

	Online  bool // 是否在线
	Logined bool // 是否登录

	Mac    string // 用户MAC地址
	OSType int64  // 系统类型

	Coder *Coder

	MQ chan Msg //IPCObj // User之间通信队列
}

func NewUser() *User {
	mq := make(chan Msg)
	coder := NewCoder()
	sid := "0" // TODO: SID通过配置获取
	return &User{Online: true, MQ: mq, SID: sid, Coder: coder}
}
