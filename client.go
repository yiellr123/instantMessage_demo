package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) (client *Client) {
	client = &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error")
		return nil
	}
	client.conn = conn
	return
}

// func main() {

// 	client := NewClient("0.0.0.0", 7050)
// 	if client == nil {
// 		fmt.Println(">>>连接服务器失败")
// 		return
// 	}
// 	fmt.Println(">>>连接服务器成功")
// 	//客户端业务逻辑
// 	select {}
// }
