package process

import (
	"fmt"
	"project/chatroom/client/model"
	"project/chatroom/common/message"
)

//客户端要维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
var CurUser model.CurUser   //登录成功后，完成对当前用户的初始化

// 在客户端显示当前在线的用户
func outputOnlineUser() {
	// 遍历一把onlineUsers
	fmt.Println("当前在线用户列表")
	for id := range onlineUsers {
		fmt.Println("用户ID:", id)
	}
}

//编写一个方法，处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes){
	//适当优化
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		// 如果当前用户不在onlineUsers中，则创建一个新的用户并添加到map中
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user
	if notifyUserStatusMes.Status == message.UserOffline {
		delete(onlineUsers, notifyUserStatusMes.UserId)
	}

	outputOnlineUser()
}


