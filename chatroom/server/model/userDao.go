package model

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"project/chatroom/common/message"
)

// 在服务器启动后，就初始化一个UserDao实例
// 把它做成全局的变量，在需要和redis操作时，直接使用
var (
	MyUserDao *UserDao
)

// 定义一个UserDao结构体
// 完成对User结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

// 使用工厂模式，创建一个UserDao的实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

// 1、根据用户id返回一个user实例+err
func (This *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	//通过给的id去redis查询
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		//错误
		if err == redis.ErrNil {
			//表示在users哈希中，没有找到对应的id
			err = ErroR_USER_NOTEXISTS
		}
		return
	}
	
    user =&User{}
	//这里需要把res反序列化User实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		return
	}
	return
}

// 完成登录的校验
// 1.Login完成对用户的校验
// 2.id,pwd都正确，则返回user实例
// 3.错误则返回对应的错误
func (This *UserDao) Login(userId int, userPwd string) (user *User, err error) {

	//先从UserDao的连接池中取出一根连接
	conn := This.pool.Get()
	defer conn.Close()
	user, err = This.getUserById(conn, userId)
	if err != nil {
		return
	}

	//证明获取到用户，但是密码可能不对
	if user.UserPwd != userPwd {
		err = ErroR_USER_PWD
		return
	}

	return
}

func (This *UserDao) Register(user *message.User) (err error) {

	//先从UserDao的连接池中取出一根连接
	conn := This.pool.Get()
	defer conn.Close()
	_, err = This.getUserById(conn, user.UserId)
	if err == nil {
		err = ErroR_USER_EXISTS
		return
	}

	//这时，说明id再redis中还不存在
	data, err := json.Marshal(user)
	if err != nil {
		return
	}

	//入库
	_, err = conn.Do("HSet", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存用户错误 err =", err)
		return
	}
	return
}
