package process

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"project/chatroom/client/utils"
	"project/chatroom/common/message"
)

// 显示登录成功后的界面...
func ShowMenu(userId int) {
	fmt.Printf("----恭喜Id:%d登录成功-----\n",userId)
	fmt.Println("----1.显示在线用户列表-----")
	fmt.Println("----2.群发消息-----")
	fmt.Println("----3.私发信息-----")
	fmt.Println("----4.退出系统-----")
	fmt.Println("请选择(1-4)")

	var key int
	var toUserId int
	var content string

	//发送信息需要经常用到SmsProcess实例，与web开发不同，只要程序不退出，实例一直都在
	smsProcess := &SmsProcess{}

	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		outputOnlineUser()
	case 2:
		fmt.Println("请输入你想对大家说点啥")
		fmt.Scanf("%s\n", &content)
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("私聊")
		fmt.Println("请输入对方的id")
		fmt.Scanf("%d\n", &toUserId)
		fmt.Println("请输入想对对方说的内容")
		fmt.Scanf("%s\n", &content)
		smsProcess.SendPrivateMes(toUserId, content)
	case 4:
		// 退出系统
		fmt.Println("正在退出系统...")
		// 向服务端发送下线消息
		mes := message.Message{
			Type: message.LogoutMesType,
			Data: "",
		}
		data, err := json.Marshal(mes)
		if err != nil {
			fmt.Println("json.Marshal err=", err)
			return
		}
		tf := &utils.Transfer{
			Conn: CurUser.Conn,
		}
		err = tf.WritePkg(data)
		if err != nil {
			fmt.Println("writePkg err=", err)
			return
		}
		os.Exit(0)
	default:
		fmt.Println("输入的选项不正确")
	}
}

// 和服务器保持通讯
func serverProcrssMes(conn net.Conn) {
	//创建一个transfer实例，不停的读取数据
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在等待读取服务端发送的信息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err=", err)
			return
		}
		//如果读取到消息，下一步处理
		switch mes.Type {
		case message.NotifyUserStatusMesType: //功能1.上线提醒
			//1.取出NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			//2.把这个用户的状态保存到map[int]User中
			if notifyUserStatusMes.Status == message.UserOnline {
				// 用户上线
				updateUserStatus(&notifyUserStatusMes)
			}
		case message.SmsMesType: //功能2.群聊
			outputGroupMes(&mes)
		case message.SmsPrivateMesType: //功能3.私聊
			outputPrivateMes(&mes)	
		// 下线操作
		default:
			fmt.Println("服务器端返回了暂时不能处理的消息类型")
		}
		//fmt.Printf("mes=%v\n", mes)
	}
}
