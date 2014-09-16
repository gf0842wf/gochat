package endpoint

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

/* tcp Endpoint:
创建 EndPoint 对象后需要调用 Init, InitCBs 和 Start 方法
对外暴露的接口:
1.使用 PutData 发送
2.读取 RecvBox 管道处理收到的消息
3.设置 OnConnectionLost 回调函数来处理断开连接, 可以为nil
4.设置 RawRecv 回调函数来处理解析收到的消息, 可以为nil
5.设置 RawSend 回调函数来处理封包发送的消息, 可以为nil
*/

type EndPoint struct {
	Conn *net.TCPConn

	SendBox chan []byte // 发送缓冲管道
	RecvBox chan []byte // 接收缓冲管道
	Ctrl    chan bool   // 控制结束 EndPoint 所有协程的

	Heartbeat int64 // 心跳超时(s), < 0表示不设置心跳

	OnConnectionLost interface{} // func(err error)  回调, err: 断开错误信息
	RawRecv          interface{} // func() (n int, err error) 回调
	RawSend          interface{} // func([]byte) 回调
}

func (ep *EndPoint) Init(conn *net.TCPConn, heartbeat int64, sendBufSize int, recvBufSize int) {
	ep.Conn = conn
	ep.Heartbeat = heartbeat
	ep.Ctrl = make(chan bool)
	ep.SendBox = make(chan []byte, sendBufSize)
	ep.RecvBox = make(chan []byte, recvBufSize)
}

// 初始化回调函数
func (ep *EndPoint) InitCBs(OnConnectionLost interface{}, RawRecv interface{}, RawSend interface{}) {
	if OnConnectionLost == nil {
		ep.OnConnectionLost = ep.onConnectionLost
	} else {
		ep.OnConnectionLost = OnConnectionLost
	}

	if RawRecv == nil {
		ep.RawRecv = ep.rawRecv
	} else {
		ep.RawRecv = RawRecv
	}

	if RawSend == nil {
		ep.RawSend = ep.rawSend
	} else {
		ep.RawSend = RawSend
	}
}

func (ep *EndPoint) PutData(data []byte) {
	ep.SendBox <- data
}

func (ep *EndPoint) recvData() {
	var err error

	for {
		if ep.Heartbeat > 0 {
			ep.Conn.SetReadDeadline(time.Now().Add(time.Duration(ep.Heartbeat) * time.Second))
		}
		_, err = ep.RawRecv.(func() (int, error))()
		if err != nil {
			break
		}
	}
	ep.OnConnectionLost.(func(err error))(err)
}

// 如果封包方式不同,需要修改这个函数,或者通过外部传入回调的方式
func (ep *EndPoint) rawRecv() (n int, err error) {
	// 接受解析长度了
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

func (ep *EndPoint) onData(data []byte) {
	ep.RecvBox <- data
}

func (ep *EndPoint) onConnectionLost(err error) {
	fmt.Println("[EP] Connection Lost:", err.Error())
	ep.Ctrl <- false
}

func (ep *EndPoint) sendData() {
	for {
		select {
		case data := <-ep.SendBox:
			ep.RawSend.(func([]byte))(data)
		case <-ep.Ctrl:
			defer close(ep.SendBox)
			defer close(ep.RecvBox)
			// 准备关闭连接, 要发完剩下的消息
			for data := range ep.SendBox {
				ep.RawSend.(func([]byte))(data)
			}
			ep.Conn.Close()

			fmt.Println("[EP] Close connection:", ep.Conn.LocalAddr)

			return
		}
	}
}

// 如果封包方式不同,需要修改这个函数,或者通过外部传入回调的方式
func (ep *EndPoint) rawSend(data []byte) {
	// 发送封装长度
	// header
	header := make([]byte, 4)
	length := len(data)
	binary.BigEndian.PutUint32(header, uint32(length))
	data = append(header, data...)
	n, err := ep.Conn.Write(data)
	if err != nil {
		fmt.Println("[EP] Error send reply, bytes:", n, "reason:", err)
		return
	}
	fmt.Println("raw send:", data)
}

func (ep *EndPoint) Start() {
	go ep.recvData()
	go ep.sendData()
}
