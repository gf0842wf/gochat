### EndPoint ###

`endpointb/EndPoint`: 使用bufio进行简单的封装, 只提供了 `BR` `BW` `Conn` `Close` 


Example:

    package main
    
    import (
    	"errors"
    	"fmt"
    	"net"
    	"net/tcpserver"
    	"net/tcpserver/endpointb"
    	"time"
    )
    
    type Bot struct {
    	endpointb.EndPoint
    }
    
    func connectionHandler(conn *net.TCPConn) {
    	bot := &Bot{}
    	bot.Init(conn)

    	// 对bot的BR, BW, Conn进行操作吧
    	bot.Conn.SetReadDeadline(time.Now().Add(10))

    	b1, err := bot.BR.ReadByte()
    	CheckError(err, "Read error!")
    	fmt.Println(b1)
    
    	bot.BW.WriteByte(100)
    	err = bot.BW.Flush()
    	CheckError(err, "Write error!")
    }
    
    func CheckError(err error, name string) {
    	if err != nil {
    		panic(errors.New(fmt.Sprintf("%s: %s", name, err.Error())))
    	}
    }
    
    func main() {
    	server := tcpserver.NewStreamServer(":7005", connectionHandler)
    	server.Start()
    }
