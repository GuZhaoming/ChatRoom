package utils

import (
	"fmt"
	"net"
	"encoding/binary"
	"encoding/json"
	"project/chatroom/common/message"
)

//这里将这些方法关联到结构体中
type Transfer struct{
	//分析他应该有哪些字段
	Conn net.Conn
	Buf [8096]byte

}

// 从服务端读取发送的数据
func (This *Transfer)ReadPkg() (mes message.Message, err error) {
	// buf := make([]byte, 8096)
	fmt.Println("读取客户端发送的数据。。。")
	_, err = This.Conn.Read(This.Buf[:4])
	if err != nil {
		return
	}
	//根据buf转成一个uint32类型
	pkgLen := binary.BigEndian.Uint32(This.Buf[0:4])

	n, err := This.Conn.Read(This.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		return
	}
	//把pkgLen包反序列化message.Message
	err = json.Unmarshal(This.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err= ", err)
		return
	}
	return
}

// 发送data方法
func (This *Transfer)WritePkg( data []byte) (err error) {
	pkgLen := uint32(len(data))
	//var buf [4]byte
	binary.BigEndian.PutUint32(This.Buf[:4], pkgLen)
	n, err := This.Conn.Write(This.Buf[:4]) //放[]byte
	if err != nil || n != 4 {
		fmt.Println("conn.Write(bytes) err=", err)
		return
	}

	//发送data信息本身
	n, err = This.Conn.Write(data) //放[]byte
	if err != nil || n != int(pkgLen) {
		fmt.Println("conn.Write(bytes) err=", err)
		return
	}
	return
}