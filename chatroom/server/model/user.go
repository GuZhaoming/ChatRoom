package model

//定义一个用户的结构体
//为了序列化和反序列化成功，
//用户信息的json字符串key和结构体的字段对应的tag名字一致
type User struct{
	UserId int `json:"userId"`
	UserPwd string `json:"userPwd"`
	UserName string `json:"userName"`
}