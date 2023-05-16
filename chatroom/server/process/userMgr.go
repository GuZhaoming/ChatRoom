package process

import (
	"fmt"
)

// UserMgr实例在服务器端有且只有一个
// 很多地方都会用到，因此定义为全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess // 存储在线用户的map，key为用户ID，value为用户实例
}

// 完成对userMgr初始化
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// 完成对onlineUsers添加
func (This *UserMgr) AddOnlineUser(up *UserProcess) {
	This.onlineUsers[up.UserId] = up
}

// 删除
func (This *UserMgr) DelOnlineUser(userId int) {
	delete(This.onlineUsers, userId)
}

// 返回当前所有在线的用户
func (This *UserMgr) GetOnlineUser() map[int]*UserProcess {
	return This.onlineUsers   //返回存储在线用户的map
}

// 根据id返回对应的值
func (This *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	
	up, ok := This.onlineUsers[userId]  //通过用户ID从onlineUsers中获取对应的用户实例
	if !ok {
		err = fmt.Errorf("用户%d 不存在", userId)
		return
	}
	return
}
