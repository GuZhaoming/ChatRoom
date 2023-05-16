package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
	SmsPrivateMesType       = "SmsPrivateMes"
	LogoutMesType ="LogoutMes"
)

const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息的内容
}

//定义两个消息
type LoginMes struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

type LoginResMes struct {
	Code    int    `json:"code"`  //500表示该用户未注册 200成功
	UsersId []int  `json:"users"` //保存用户id的切片
	Error   string `json:"error"` //错误信息
}

type RegisterMes struct {
	User User `json:"user"` //类型为User结构体
}

type RegisterResMes struct {
	Code  int    `json:"code"`  //400表示该用户已经占用 200成功
	Error string `json:"error"` //错误信息
}

//配合服务器端推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId int `json:"userId"`
	Status int `json:"status"`
}

//增加一个SmsMes  发送的消息
type SmsMes struct {
	Content string `json:"content"`
	User           //继承
}

//私聊
type SmsPrivateMes struct {
	FromUserId int    `json:"fromUserId"`
	ToUserId   int    `json:"toUserId"`
	Content    string `json:"content"`
	User              //继承
}
type LogoutMes struct {
	UserId int // 用户ID
	// 其他字段...
}
