package tcpserver

import (
	"fmt"
	"net"
)

/* tcp server:
import (
	"net/tcpserver"
)

func handleClient(conn *net.TCPConn) {
	fmt.Println("conn...")
}

func main() {
	svr := tcpserver.NewStreamServer(":7005", handleClient)
	svr.Start()
}
*/

type StreamServer struct {
	Address           *net.TCPAddr
	Listener          *net.TCPListener
	ConnectionHandler func(conn *net.TCPConn)
}

func (svr *StreamServer) Start() {
	defer svr.Listener.Close()

	for {
		conn, err := svr.Listener.AcceptTCP()
		if err != nil {
			fmt.Println("[SS] Accept failed:", err)
			continue
		}
		go svr.ConnectionHandler(conn)
	}
}

func NewStreamServer(addr string, connectionHandler func(conn *net.TCPConn)) *StreamServer {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}

	fmt.Println("[SS] Server start:", addr)

	svr := StreamServer{Address: tcpAddr, Listener: listener, ConnectionHandler: connectionHandler}

	return &svr
}
