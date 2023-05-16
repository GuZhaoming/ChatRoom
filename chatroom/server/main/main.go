package main

import (
	"fmt"
	"net"
	"project/chatroom/server/model"
	"time"
)

// 处理和客户端的通讯·
func mainprocess(conn net.Conn) {
	defer conn.Close()

	//调用总控,创建一个总控实例
	processor := &Processor{
		Conn: conn,
	}
	err := processor.readClientMsg()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误",err)
		return
	}
}

func init(){
   //当服务器启动时，初始化redis的连接池
   initPool("localhost:6379",16,0,300*time.Second)
   //调用此方法，创建一个和redis打交道的实例
   initUserDao()
}

//编一个函数，完成对UserDao的初始化任务
func initUserDao(){
	//注意初始化的顺序问题
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	fmt.Println("服务器在8889端口监听")
	//创建一个TCP服务器监听端口
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.listen err", err)
		return
	}
	defer listen.Close()
	for {
		fmt.Println("等待客户端连接服务器")
		//循环等待连接，创建一个表示连接的实例
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.accept err", err)
		}
		//一旦链接成功，则启动一个协程和客户端保持通讯
		go mainprocess(conn)
	}
}
