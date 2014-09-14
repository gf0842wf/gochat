package protos

//import (
//	"encoding/json"
//	"fmt"
//)

import (
	"types"
)

// ----------收发的消息-------------

//type Msg struct {
//	MsgType byte
//	MsgBody []byte
//	MsgJson *map[string]interface{}
//}

//func NewSendMsg(msgType byte) *Msg {
//	return &Msg{MsgType: msgType}
//}

//func Unmarshal(data []byte, pobj *Msg) (err error, kind byte) {
//	// kind: 0-json消息, 1-二进制消息
//	msgType := data[0]
//	pobj.MsgType = msgType
//	if msgType < 100 {
//		fmt.Println(pobj.MsgType)
//		err = json.Unmarshal(data[1:], pobj.MsgJson)
//		kind = 0
//	} else if msgType == 100 {
//		// 二进制消息
//		pobj.MsgBody = data[1:]
//		kind = 1
//	}

//	return
//}

//func Marshal(pobj *Msg) (data []byte, err error) {
//	if pobj.MsgType < 100 {
//		lmsgType := []byte{pobj.MsgType}
//		data, err = json.Marshal(pobj.MsgJson)
//		data = append(lmsgType, data...)
//	} else if pobj.MsgType == 100 {
//		// 二进制消息
//		data = pobj.MsgBody
//	}

//	return
//}

var NetProtoHandlers map[byte]func(*types.User, *types.Msg) (resp []byte, err error)

func init() {
	NetProtoHandlers = map[byte]func(*types.User, *types.Msg) (resp []byte, err error){
		0: handle_shake,
		//1: handle_nop,
		2: handle_login,

		//3: handle_offline,
		//4: handle_chat,
		//6: handle_notify,

		//5: handle_cmd,

	}
}
