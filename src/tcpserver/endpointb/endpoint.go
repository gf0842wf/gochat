package endpointb

import (
	"bufio"
	"net"
)

/* tcp Endpoint:
和 endpoint 相比, endpointb 不再使用管道缓存了, 直接用 bufio
这样可以更加精准的控制 读取数据解析过程的conn超时

外部组合是这个类后, 需要操作 BR BW和Conn来读写和控制超时
*/

type EndPoint struct {
	Conn *net.TCPConn

	BR *bufio.Reader // 接受buf: BW.Read...
	BW *bufio.Writer // 发送buf: BW.Write..., BW.Flush
}

func (ep *EndPoint) Init(conn *net.TCPConn) {
	ep.Conn = conn
	ep.BR = bufio.NewReader(conn)
	ep.BW = bufio.NewWriter(conn)
}

func (ep *EndPoint) Close() {
	ep.Conn.Close()
}
