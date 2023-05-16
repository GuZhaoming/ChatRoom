package model

import (
	"net"
	"project/chatroom/common/message"
)


//curUser做成全局的，维护当前登录的用户
type CurUser struct {
	Conn net.Conn
	message.User
}