package protos

import (
	"types"
)

var NetProtoHandlers map[byte]func(*types.User, []byte) (resp []byte, err error)

func init() {
	NetProtoHandlers = map[byte]func(*types.User, []byte) (resp []byte, err error){
		0: handle_shake,
		1: handle_nop,
		2: handle_login,

		//3: handle_offline,
		4: handle_chat,
		//6: handle_notify,

		//5: handle_cmd,

	}
}
