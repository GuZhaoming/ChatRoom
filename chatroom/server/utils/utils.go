package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"project/chatroom/common/message"
)

// 定义一个名为Transfer的结构体，用于数据读写
type Transfer struct {
	Conn net.Conn   //表示客户端与服务器的连接
	Buf  [8096]byte //定义一个缓冲区，用于读写数据
}

//从客户端读取数据
func (This *Transfer) ReadPkg() (mes message.Message, err error) {
	fmt.Println("读取客户端发送的数据...")
	// 读取前4个字节，表示要读取的数据的长度
	_, err = This.Conn.Read(This.Buf[:4])
	if err != nil {
		return
	}
	//根据读取到的长度解码出包的长度
	pkgLen := binary.BigEndian.Uint32(This.Buf[0:4])

	//读取pkgLen长度的数据
	n, err := This.Conn.Read(This.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		return
	}

	//将读取到的数据进行反序列化，得到message.Message类型
	err = json.Unmarshal(This.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err= ", err)
		return
	}
	return
}

// 向客户端发送数据
func (This *Transfer) WritePkg(data []byte) (err error) {
	// 计算待发送数据的长度
	pkgLen := uint32(len(data))

	// 将数据的长度编码成4个字节，放入缓冲区的前4个字节中
	binary.BigEndian.PutUint32(This.Buf[0:4], pkgLen)

	// 发送长度信息
	n, err := This.Conn.Write(This.Buf[:4])
	if err != nil || n != 4 {
		fmt.Println("conn.Write(bytes) err=", err)
		return
	}

	// 发送数据信息本身
	n, err = This.Conn.Write(data)
	if err != nil || n != int(pkgLen) {
		fmt.Println("conn.Write(bytes) err=", err)
		return
	}
	return
}
