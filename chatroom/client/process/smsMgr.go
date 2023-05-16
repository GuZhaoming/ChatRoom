package process

import (
	"encoding/json"
	"fmt"
	"project/chatroom/common/message"
)

//群聊
func outputGroupMes(mes *message.Message){
	//显示
	//1.反序列化mes.Data
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data),&smsMes )
	if err != nil{
		fmt.Println("json.Unmarshal err= ",err)
	}

	//显示信息
	info:= fmt.Sprintf("用户id:\t%d 对大家说:\t%s",smsMes.UserId,smsMes.Content)
    fmt.Println(info)
	fmt.Println()
}

//私聊
func outputPrivateMes(mes *message.Message){
	var smsPrivateMes message.SmsPrivateMes
	err := json.Unmarshal([]byte(mes.Data),&smsPrivateMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	info := fmt.Sprintf("用户id:\t%d 对你说:\t%s", smsPrivateMes.FromUserId,smsPrivateMes.Content)
	fmt.Println(info)
	fmt.Println()
}