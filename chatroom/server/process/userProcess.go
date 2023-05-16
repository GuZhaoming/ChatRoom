package process

import (
	"encoding/json"
	"fmt"
	"net"
	"project/chatroom/common/message"
	"project/chatroom/server/model"
	"project/chatroom/server/utils"
)

type UserProcess struct {
	//字段：表示连接
	Conn net.Conn
	//增加字段，表示连接的用户ID
	UserId int
}

// 编写通知所有在线用户的方法
// 通知其他用户我上线了
func (This *UserProcess) NotifyOtherOnlineUser(userId int) {
	//遍历onlineUsers,然后一个一个发送
	for id, up := range userMgr.onlineUsers {
		if id == userId {
			continue
		}
		//开始通知，单独写一个方法
		up.NotifyMeOnline(userId)
	}
}


// 处理上线通知请求
func (This *UserProcess) NotifyMeOnline(userId int) {
	//组装NotifyuserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline
	//将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}
	mes.Data = string(data)
	data, err = json.Marshal(mes) // 对整个消息体进行序列化
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}

	//发送，创建transfer实例
	tf := &utils.Transfer{
		Conn: This.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err = ", err)
		return
	}
}



// 专门处理登录请求
func (This *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//1.先从mes中取出mes.Data,并直接反序列化成loginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		return
	}

	//1.先声明一个resMes信息对象（消息类型，消息内容）
	var serMes message.Message
	serMes.Type = message.LoginResMesType

	//2.再声明一个loginReMes状态码对象，完成赋值
	var loginResMes message.LoginResMes

	//1.需要到redis数据库完成验证
	//使用model.MyUserDao到redis验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {

		if err == model.ErroR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ErroR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}
	} else {
		fmt.Println(user, "登录成功")
		loginResMes.Code = 200

		//将当前在线用户的id放入到loginResMes.UsersId
		This.UserId = loginMes.UserId

		//这里，因为用户登录成功，我们把用户放入userMgr里
		userMgr.AddOnlineUser(This)

		//通知其他在线用户，我上线了
		This.NotifyOtherOnlineUser(loginMes.UserId)

		for id := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
	}

	//3.将loginResMes序列化
	data, err := json.Marshal(loginResMes) //data：json编码的字节数组
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}
	//4.将data赋值给resMes
	serMes.Data = string(data)
	//5.对serMes进行序列化，准备发送
	data, err = json.Marshal(serMes)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}
	//6，发送data
	tf := &utils.Transfer{
		Conn: This.Conn,
	}
	err = tf.WritePkg(data)
	return
}

// 处理注册请求
func (This *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		return
	}

	//1.先声明一个resMes信息对象（消息类型，消息内容）
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	//1.需要到redis数据库完成验证
	//使用model.MyUserDao到redis验证
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ErroR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ErroR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误"
		}
	} else {
		registerResMes.Code = 200
	}

	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}
	//4.将data赋值给resMes
	resMes.Data = string(data)
	//5.对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}
	//6，发送data
	//使用分层模式，我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: This.Conn,
	}

	err = tf.WritePkg(data)
	return

}




