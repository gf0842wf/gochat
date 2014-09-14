package share

import (
	"types"
)

var Clients map[uint32]*types.User // 这个应该加锁,如果是多核的话
