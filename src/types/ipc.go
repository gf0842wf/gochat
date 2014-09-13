package types

// 用户之间通信载体

import (
	"encoding/json"
)

type IPCObj struct {
	SrcID  int32   // 发送方用户ID
	DstID  int32   // 接收放用户ID
	AuxIDs []int32 // 目标用户ID集合(用于组播)
	No     int16   // 服务号
	JsObj  []byte  // 投递的 JSON STRING
	Time   int64   // 发送时间
}

// IPCObj to json
func (obj *IPCObj) Json() []byte {
	jsobj, _ := json.Marshal(obj)
	return jsobj
}
