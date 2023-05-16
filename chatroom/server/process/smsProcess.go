package process

import (
	"encoding/json"
	"fmt"
	"net"
	"project/chatroom/common/message"
	"project/chatroom/server/utils"
)

type SmsProcess struct {
}

//用于转发群聊消息
func (This *SmsProcess) SendGroupMes(mes *message.Message) {

	//遍历服务器端的onlineUsers map[int]*UserProcess
	//将消息转发取出

	//取出mes的内容SmsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err ", err)
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err ", err)
		return
	}
	for id, up := range userMgr.onlineUsers {

		//过滤自己
		if id == smsMes.UserId {
			continue
		}
		This.SendMesToEachOnlineUser(data, up.Conn)
	}
}

//用于将消息发送给每一个在线用户
func (This *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	//创建Transfer实例，发送data
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err= ", err)
	}
}

// 用于转发私聊消息
func (This *SmsProcess) SendPrivateMes(mes *message.Message) {

	//取出mes的内容SmsPrivateMes
	var smsPrivateMes message.SmsPrivateMes
	err := json.Unmarshal([]byte(mes.Data), &smsPrivateMes)
	if err != nil {
		fmt.Println("json.Unmarshal err ", err)
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err ", err)
		return
	}

	//获取接收方用户的UserProcess
	toUserId := smsPrivateMes.ToUserId
	if up, ok := userMgr.onlineUsers[toUserId]; ok {
		This.SendMesToEachOnlineUser(data, up.Conn)
	} else {
		fmt.Println("用户不在线")
	}
}







