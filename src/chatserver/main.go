package main

// 聊天服务器(一般多个),客户端连接到此

import (
	"fmt"
	"net"
)

import (
	"share"
	"tcpserver"
	"tcpserver/endpoint"
	"types"
)

type Bot struct {
	endpoint.EndPoint

	User    *types.User
	Manager *TCPServerManager
}

func (bot *Bot) OnConnectionLost(err error) {
	fmt.Println("Connection Lost:", err.Error())

	bot.Ctrl <- false
	if bot.User.UID > 0 {
		// TODO: 发送hub服务器该用户下线
		fmt.Println("LOST:", bot.User.UID)
		share.Clients.Delete(bot.User.UID)
	}
}

func (bot *Bot) Handle() {
	for {
		select {
		case data := <-bot.RecvBox:
			// 解密
			bot.User.Coder.Decode(data)
			ack, err := HandleNetProto(bot, data)
			if err != nil {
				// 断开连接
				fmt.Println(err.Error())
				bot.Conn.Close()
				return
			}
			if ack != nil {
				// 加密
				bot.User.Coder.Encode(ack)
				bot.PutData(ack)
			}
		case data := <-bot.User.MQ: // 转发聊天消息
			fmt.Println("MQ:", data)
			// 加密
			bot.User.Coder.Encode(data)
			bot.PutData(data)
		case data := <-bot.SendErrBox: // 发送失败的消息,视为离线消息
			// 解密
			bot.User.Coder.Decode(data)
			// TODO: 离线消息入库,只入库聊天消息,把聊天消息状态改为离线
			fmt.Println("offchat:", data)
		}
	}
}

type TCPServerManager struct {
	Address string
}

func (m *TCPServerManager) connectionHandler(conn *net.TCPConn) {
	bot := &Bot{}
	bot.Init(conn, 300, 16, 12, 8)
	bot.InitCBs(bot.OnConnectionLost, nil, nil)
	bot.Manager = m
	user := types.NewUser(16) // 这个参数是聊天发送队列
	bot.User = user

	go bot.Handle()
	bot.Start()
}

func (m *TCPServerManager) Start() {
	server := tcpserver.NewStreamServer(m.Address, m.connectionHandler)
	server.Start()
}

func main() {
	manager := &TCPServerManager{Address: ":7005"}
	go manager.Start()

	waiting := make(chan bool)
	<-waiting

}
