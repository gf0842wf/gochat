package main

// 节点服务器,客户端连接到此

import (
	"fmt"
	"net"
	"net/tcpserver"
	"net/tcpserver/endpoint"
)

var SID uint32

func init() {
	SID = 1
}

type Bot struct {
	endpoint.EndPoint

	UID     uint32
	Manager *TCPServerManager
}

func (bot *Bot) OnConnectionLost(err error) {
	fmt.Println("Connection Lost:", err.Error())
	fmt.Println(bot.Manager.Clients)

	bot.Ctrl <- false
	delete(bot.Manager.Clients, bot.ID)
}

func (bot *Bot) Handle() {
	for {
		select {
		case data := <-bot.RecvBox:
			fmt.Println("Recv:", string(data))
			bot.PutData(data)
			// to do something
		}
	}
}

type TCPServerManager struct {
	Address string
	Clients map[uint32]*Bot // 这个应该加锁,如果是多核的话
}

func (m *TCPServerManager) connectionHandler(conn *net.TCPConn) {
	bot := &Bot{}
	bot.Init(conn, 10, 16, 12)
	bot.InitCBs(bot.OnConnectionLost, nil, nil)
	bot.ID = SID
	bot.Manager = m
	SID++

	m.Clients[bot.ID] = bot

	go bot.Handle()
	bot.Start()
}

func (m *TCPServerManager) Start() {
	server := tcpserver.NewStreamServer(m.Address, m.connectionHandler)
	server.Start()
}

func main() {
	manager := &TCPServerManager{Address: ":7005"}
	manager.Clients = make(map[uint32]*Bot, 100)
	go manager.Start()

	waiting := make(chan bool)
	<-waiting

}
