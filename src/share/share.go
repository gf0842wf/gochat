package share

import (
	"misc/zmap"
)

// uid 2 user
var Clients *zmap.SafeMap

func init() {
	Clients = zmap.NewSafeMap()
}
