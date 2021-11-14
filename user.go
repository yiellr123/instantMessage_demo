package main

import "net"

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

	} else {
		u.server.BroadCast(u, msg)
	}
}
