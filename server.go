package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnLineMap map[string]*User
	mapLock   sync.RWMutex
	//消息广播channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnLineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "] " + user.Name + ":" + msg
	s.Message <- sendMsg
}

func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message
		//将msg发送给全部在线user
		s.mapLock.Lock()
		for _, cli := range s.OnLineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}
func (s *Server) Handler(conn net.Conn) {
	//当前连接业务
	// fmt.Println("连接建立成功")
	user := NewUser(conn)
	//用户上线,将用户加入到onlinemap中
	s.mapLock.Lock()
	s.OnLineMap[user.Name] = user
	s.mapLock.Unlock()
	//广播当前用户上线消息
	s.BroadCast(user, "已上线")
	//当前handler阻塞  user是指针类型的, 这里的user没有了, map里面的user也会跟着没有?????

	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				s.BroadCast(user, "下线") //关闭连接时,下线
				return
			}
			if err != nil && err != io.EOF { //没有读到末尾
				fmt.Println("Conn Read err:", err)
			}
			//提取用户的消息,去除'\n'
			msg := string(buf[:n-1])
			//将得到的消息进行广播
			s.BroadCast(user, msg)
		}
	}()

	select {}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}
	defer listener.Close()
	//启动监听Message的goroutine
	go s.ListenMessager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listen.accept err:", err)
			continue
		}
		//do handler
		go s.Handler(conn)
	}

}
