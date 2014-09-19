package tcpclient

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

type Connection struct {
	Addr *net.TCPAddr
	Conn *net.TCPConn

	sendBufSize int
	recvBufSize int

	SendBox chan []byte // 发送缓冲管道(外部不直接用,用PutData)
	RecvBox chan []byte // 接收缓冲管道(外部直接用)
	Ctrl    chan bool   // 控制结束 EndPoint 所有协程的
}

func (ep *Connection) PutData(data []byte) {
	ep.SendBox <- data
}

func (ep *Connection) recvData() {
	var err error

	for {
		_, err = ep.RawRecv()
		if err != nil {
			break
		}
	}
	ep.onConnectionLost(err)
}

func (ep *Connection) RawRecv() (n int, err error) {
	// header
	header := make([]byte, 4)
	n, err = io.ReadFull(ep.Conn, header)
	if err != nil {
		err = errors.New("[EP] Error recv header:" + strconv.Itoa(n) + ":" + err.Error())
		return
	}

	// data
	length := binary.BigEndian.Uint32(header)
	data := make([]byte, length)
	n, err = io.ReadFull(ep.Conn, data)
	if err != nil {
		err = errors.New("[EP] Error recv msg:" + strconv.Itoa(n) + ":" + err.Error())
		return
	}
	ep.onData(data)

	return
}

func (ep *Connection) onData(data []byte) {
	ep.RecvBox <- data
}

func (ep *Connection) onConnectionLost(err error) {
	fmt.Println("[EP] Connection Lost:", err.Error())
	// 断开后重连
	ep.Ctrl <- false
	ep.Reconnect()
}

func (ep *Connection) sendData() {
	for {
		select {
		case data := <-ep.SendBox:
			ep.RawSend(data)
		case <-ep.Ctrl:
			defer close(ep.SendBox)
			defer close(ep.RecvBox)
			// 准备关闭连接, 要发完剩下的消息
			for data := range ep.SendBox {
				ep.RawSend(data)
			}
			ep.Conn.Close()

			fmt.Println("[EP] Close connection:", ep.Conn.LocalAddr)

			return
		}
	}
}

func (ep *Connection) RawSend(msg []byte) {
	// 发送封装长度
	// header
	header := make([]byte, 4)
	length := len(msg)
	binary.BigEndian.PutUint32(header, uint32(length))
	data := append(header, msg...)
	n, err := ep.Conn.Write(data)
	if err != nil {
		fmt.Println("[EP] Error send reply, bytes:", n, "reason:", err)
		return
	}
	fmt.Println("raw send:", data)
}

func (ep *Connection) Init() {
	ep.Ctrl = make(chan bool)
	ep.SendBox = make(chan []byte, ep.sendBufSize)
	ep.RecvBox = make(chan []byte, ep.recvBufSize)

	go ep.recvData()
	go ep.sendData()
}

func (ep *Connection) Connect() error {
	conn, err := net.DialTCP("tcp", nil, ep.Addr)
	if err == nil {
		ep.Conn = conn
	}

	return err
}

func (ep *Connection) Reconnect() {
	for {
		fmt.Println("Trying reconnect..")
		time.Sleep(time.Second * 5)
		err := ep.Connect()
		if err == nil {
			ep.Init()
			fmt.Println("Reconneced.")
			break
		}
	}
}

func NewConnection(addr string, sendBufSize int, recvBufSize int) *Connection {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		panic(err)
	}

	client := Connection{}
	client.Addr = tcpAddr
	client.sendBufSize = sendBufSize
	client.recvBufSize = recvBufSize

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	client.Init()

	return &client
}
