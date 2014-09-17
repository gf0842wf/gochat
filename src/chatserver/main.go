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
		//delete(bot.Manager.Clients, bot.User.UID)
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
		case data := <-bot.User.MQ:
			fmt.Println("MQ:", data)
			// 加密
			bot.User.Coder.Encode(data)
			bot.PutData(data)
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
