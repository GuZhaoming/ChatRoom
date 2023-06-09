package process

import (
	"encoding/json"
	"fmt"
	"project/chatroom/client/utils"
	"project/chatroom/common/message"
)

type SmsProcess struct {
}

// 发送群发的消息
func (This *SmsProcess) SendGroupMes(content string) (err error) {
	//1.创建一个Mes
	var mes message.Message
	mes.Type = message.SmsMesType

	//2.创建一个SmsMes实例
	var smsMes message.SmsMes
    smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	//3.序列化
	data, err := json.Marshal(smsMes)
	if err != nil{
		fmt.Println("SendGroupMes json.Marshal err ",err.Error())
		return
	}
	mes.Data = string(data)

	//4.对mes再次序列化
	data, err = json.Marshal(mes)
	if err != nil{
		fmt.Println("SendGroupMes json.Marshal err ",err.Error())
		return
	}

	//5.将mes发送给服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}

	//6.发送
	err = tf.WritePkg(data)
	if err != nil{
		fmt.Println("SendGroupMes err ",err.Error())
		return
	}
	return
}

//发送私聊消息
func (This *SmsProcess) SendPrivateMes(userId int,content string)(err error){
	//1.
	var mes message.Message
	mes.Type = message.SmsPrivateMesType

	//2.创建一个SmsPrivateMes实例
	var smsPrivateMes message.SmsPrivateMes
	smsPrivateMes.Content = content
	smsPrivateMes.FromUserId = CurUser.UserId
	smsPrivateMes.ToUserId = userId

	//3.
    data, err := json.Marshal(smsPrivateMes)
	if err != nil{
		fmt.Println("SendPrivateMes json.Marshal err ",err.Error())
		return
	}

	mes.Data = string(data)

	//4.
    data, err = json.Marshal(mes)
	if err != nil{
		fmt.Println("SendPrivateMes json.Marshal err ",err.Error())
		return
	}

	//5.将mes发送给服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}

    err = tf.WritePkg(data)
	if err != nil{
		fmt.Println("SendPrivateMes err ",err.Error())
		return
	}
	return

}