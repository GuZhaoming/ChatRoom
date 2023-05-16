package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"project/chatroom/client/utils"
	"project/chatroom/common/message"
)

type UserProcess struct {
}

// 给关联一个用户登录的方法
// 登录
func (This *UserProcess) Login(userId int, userPwd string) (err error) {
	//1.链接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2.创建mes结构体实例用于发送消息
	var mes message.Message
	mes.Type = message.LoginMesType

	//3.创建一个loginMes结构体实例
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	//4，将loginMes实例序列化为二进制数据
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}

	//5. 把data赋给mes.Data字段
	mes.Data = string(data)

	//6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}

	//7.这个是个data是我们要发送的信息
	//7.1先把data的长度发送给服务器
	pkgLen := uint32(len(data))
	var buf [4]byte
	//将一个uint32类型的整数值pkgLen转换为一个长度为4的切片，并存到buf切片中
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//发送信息长度到服务端
	n, err := conn.Write(buf[:4])
	if err != nil || n != 4 {
		fmt.Println("conn.Write(bytes) err=", err)
		return
	}
	//fmt.Printf("客户端发送信息的长度ok,len(data)=%d,string(data)=%v\n", len(data), string(data))

	//发送消息本身到服务端Write(切片)
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) err=", err)
		return
	}

	//处理服务器返回消息
	//创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) err", err)
	}

	//将mes的Data部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)

	if loginResMes.Code == 200 {
		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline

		//可以显示当前在线用户的列表，遍历loginResMes.UsersId
		fmt.Println("当前在线用户列表如下:")
		for _, v := range loginResMes.UsersId {
			//如何不显示自己的id
			if v == userId {
				continue
			}
			fmt.Println("用户id:\t", v)
			//完成客户端的onlineUser完成初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")

		//登陆成功后启动一个协程
		//该协程保持和服务端的通讯，如果服务器有数据推给客户端
		//则接收并显示在客户端的终端
		go serverProcrssMes(conn)
		//1.显示登陆成功后的菜单[循环显示菜单]
		for {
			ShowMenu(userId)
		}
	} else {
		fmt.Println(loginResMes.Error)	
	}
	return
}

// 注册
func (This *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	//1.链接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2.准备通过conn发送信息给服务端
	var mes message.Message
	mes.Type = message.RegisterMesType

	//3.创建一个RegisterMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	//4，将registerMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}

	//5. 把data赋给mes.Data字段
	mes.Data = string(data)

	//6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}

	//创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}

	//发送data给服务器
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息错误 err", err)
		return
	}

	mes, err = tf.ReadPkg() //这个mes就是RegisterResMes
	if err != nil {
		fmt.Println("readPkg(conn) err", err)
		return
	}

	//将mes的Data部分反序列化成RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功,重新登录一把")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}

	return
}





