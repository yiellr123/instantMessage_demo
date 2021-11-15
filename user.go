package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage() //死循环 ,所以return后不会被销毁??

	return user
}

//监听当前user channel ,有消息就发给对端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}

// OnLine 用户的上线业务
func (u *User) OnLine() {
	u.server.mapLock.Lock()
	u.server.OnLineMap[u.Name] = u
	u.server.mapLock.Unlock()
	//广播当前用户上线消息
	u.server.BroadCast(u, "已上线")
}

// OffLine 用户的下线业务
func (u *User) OffLine() {
	u.server.mapLock.Lock()
	delete(u.server.OnLineMap, u.Name)
	u.server.mapLock.Unlock()
	u.server.BroadCast(u, "下线")
}

func (u *User) sendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

// DoMessage 用户处理消息的业务
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户有哪些
		u.server.mapLock.Lock()
		for _, user := range u.server.OnLineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线..\n"
			// 不能使用此广播方式,只通过当前的conn来发送给当前用户
			// u.server.BroadCast(u, onlineMsg)

			u.sendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式:rename|张三
		newName := strings.Split(msg, "|")[1] //Split根据|分割成[]string,[0]为左侧,[1]为右侧
		//判断那么是否存在
		_, ok := u.server.OnLineMap[newName]
		if ok {
			u.sendMsg("当前用户名被使用\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnLineMap, u.Name) //将之前初始化name删掉
			u.server.OnLineMap[newName] = u
			u.server.mapLock.Unlock()
			u.Name = newName
			u.sendMsg("您已更新用户名:" + u.Name + "\n")
		}
	} else {
		u.server.BroadCast(u, msg)
	}
}
