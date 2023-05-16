package main

import (
	"fmt"
	"io"
	"net"
	"project/chatroom/common/message"
	"project/chatroom/server/process"
	"project/chatroom/server/utils"
)

// 先创建一个Processor的总控结构体
type Processor struct {
	//与客户端建立的连接，这里的连接通过（客户端发来的连接）服务端接收的连接来传进来
	Conn net.Conn
}

// serverProcessMes 用于处理客户端发来的消息
func (This *Processor) serverProcessMes(mes *message.Message) (err error) {

	fmt.Println("mes=", mes)

	switch mes.Type {
	case message.LoginMesType: // 处理登录
		// 创建一个UserProcess实例
		up := process.UserProcess{
			Conn: This.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType: // 处理注册
		up := process.UserProcess{
			Conn: This.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType: // 处理群发
		// 创建一个SmsProcess实例
		smsProcess := process.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	case message.SmsPrivateMesType: //处理私聊
		smsProcess := process.SmsProcess{}
		smsProcess.SendPrivateMes(mes)
	default:
		fmt.Println("消息类型不存在，无法处理")
	}
	return
}

// process2 用于循环读取客户端发送的消息
func (This *Processor) readClientMsg() (err error) {
	for {
		// 读取数据包封装成函数
		// 创建一个Transfer实例完成读包
		tf := &utils.Transfer{
			Conn: This.Conn,
		}
		mes, err := tf.ReadPkg() //读取客户端发来的消息
		if err != nil {
			//错误1.连接关闭了，Read()方法会返回一个io.EOF错误
			if err == io.EOF {
				fmt.Println("客户端退出，服务器也退出")
				return err
			} else {
				//错误2.读包问题
				fmt.Println("readPkg err", err)
				return err
			}
		}

		err = This.serverProcessMes(&mes) // 处理客户端发来的消息
		if err != nil {
			return err
		}
	}
}
