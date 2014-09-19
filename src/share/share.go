package share

import (
	"misc/zmap"
	"tcpclient"
)

// uid 2 user
var Clients *zmap.SafeMap
var HubClient *tcpclient.Connection

func init() {
	Clients = zmap.NewSafeMap()

	// TODO: hub server通过配置获取
	//HubClient = tcpclient.NewConnection("127.0.0.1:7010", 8, 8)
}
