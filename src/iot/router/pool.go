package main

/**
* 用户池后期改造为redis或者mongodb持久化
 */

import (
	"iot/internal/rpc"
	"sync"
)

type Pool struct {
	m1     sync.Mutex
	nodes map[string]*node
}

//用户session
type session struct {
	id      string
	nodeId string //用户所附着的nodeId
	item    []*item
}

type item struct {
	plat        int    //终端类型
	online      bool   //推送接口是否在线
	authCode    string //业务层鉴权码
	login       bool   //业务层是否已经登录
	deviceToken string //设备token
}

type node struct {
	id        string         //comet id
	rpcClient *rpc.RpcClient //router连接到本node的RPC客户端句柄
	tcpAddr   string         //node对外开放tcp服务地址
	online    int            //node在线统计
	ch        chan int       //node rpc 状态通知chan
}

func (p *Pool) insertNode(id string, c *node) {
	p.m1.Lock()
	defer p.m1.Unlock()
	p.nodes[id] = c
}

func (p *Pool) findNode(id string) *node {
	p.m1.Lock()
	defer p.m1.Unlock()
	c := p.nodes[id]
	return c
}

func (p *Pool) deleteNode(id string) {
	p.m1.Lock()
	defer p.m1.Unlock()
	delete(p.nodes, id)
}

func (p *Pool) nodeAdd(id string) {
	p.m1.Lock()
	defer p.m1.Unlock()
	c := p.nodes[id]
	if c != nil {
		c.online = c.online + 1
	}
}

func (p *Pool) nodeSub(id string) {
	p.m1.Lock()
	defer p.m1.Unlock()
	c := p.nodes[id]
	if c != nil {
		c.online = c.online - 1
	}
}

//选择负载最低的comet
func (p *Pool) balancer() *node {
	p.m1.Lock()
	defer p.m1.Unlock()
	minLoad := 0
	var c *node
	for _, v := range p.nodes {
		if minLoad == 0 {
			minLoad = v.online
			c = v
		}
		if v.online < minLoad {
			minLoad = v.online
			c = v
		}
	}
	return c
}
