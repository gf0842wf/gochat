package share

import (
	"misc/zmap"
	//"types"
)

//var Clients map[uint32]*types.User // 这个应该加锁,如果是多核的话

// uid 2 user
var Clients *zmap.SafeMap

func init() {
	Clients = zmap.NewSafeMap()
}
