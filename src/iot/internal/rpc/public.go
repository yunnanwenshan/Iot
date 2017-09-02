//rpc 提供comet to router，router to comet的RPC调用
package rpc

import (
	"errors"
)

var (
	ERROR_RESPONSE = errors.New("error response code")
)

const (
	RPC_RET_SUCCESS = 0
	RPC_RET_FAILED  = -1
)

const (
	STATE_OFFLINE = 0
	STATE_ONLINE  = 1
)

/* Node向Router注册
* 方向:Node->Router
 */
type NodeRegister struct {
	NodeId string //node ID
	TcpAddr string //node对外开放tcp服务地址
	RpcAddr string //反连地址(node rpc服务地址)
}

/* 鉴权
*方向:Node->Router
 */
type AuthRequest struct {
	Id       string
	Termtype int
	Code     string
}

/*推送请求
* 方向:router->Node
 */
type MsgUpwardRequst struct {
	Id       string
	Termtype int
	Msg      string
}

/* 用户socket状态通知
* 方向:Node->Router
 */
type StateNotify struct {
	Id       string
	Termtype int
	Token    string
	NodeId  string //附着Node ID
	State    int    //1-online 0-offline
}

/* 踢人下线
*方向:router->Node
 */
type KickRequst struct {
	Id       string
	Termtype int
	Token    string
	/* added by liang @ 2016-07-11
	0x01:重复登录，同终端类型客户端登录
	0x02:互斥登录，Android/iOS互斥
	0x03:session超时
	*/
	Reason int
}

/*推送请求
* 方向:router->Node
 */
type PushRequst struct {
	Tp         int //消息类型
	Flag       int //IOS 声音提示
	Id         string
	Termtype   int
	AppleToken string
	Msg        string
}

/* 心跳检测
* 方向：RPC客户端->RPC服务端
 */
type Ping struct {
}

//公共应答
type Response struct {
	Code int //0-成功 -1-失败
}
