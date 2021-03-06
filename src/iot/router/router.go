package main

import (
	"iot/internal/logs"
	"iot/internal/rds"
	"iot/internal/rpc"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/widuu/goini"
)

type Router struct {
	//rpc服务
	routerRpcAddr string
	maxRpcInFight int
	rpcServer     *rpc.RpcServer

	//负载路由
	httpBindAddr string
	nodeExit    chan string //nodeId exit channel
	pool         *Pool

	//离线消息
	//store *redisstore.Storager
	store *rds.Storager

	//系统控制
	exit chan struct{}
	wg   sync.WaitGroup
}

func (p *Router) Init() {
	conf := goini.SetConfig("./config.ini")
	logs.Logger.Debug("--------OnInit--------")
	//RPC
	{
		p.routerRpcAddr = conf.GetValue("router", "rpcAddr")
		s := conf.GetValue("router", "rpcServerCache")
		p.maxRpcInFight, _ = strconv.Atoi(s)
		logs.Logger.Debug("----router rpc addr=", p.routerRpcAddr, " cache=", p.maxRpcInFight)
	}

	//HTTP
	{
		p.httpBindAddr = conf.GetValue("http", "bindAddr")
		logs.Logger.Debug("----http addr=", p.httpBindAddr, " cache=", p.maxRpcInFight)
	}

	p.nodeExit = make(chan string)

	p.pool = new(Pool)
	p.pool.nodes = make(map[string]*node)
	//	p.pool.sessions = make(map[string]*session)

	//REDIS
	{
		dbconn := conf.GetValue("redis", "conn")
		password := conf.GetValue("redis", "password")
		password = strings.TrimSpace(password)
		databaseS := conf.GetValue("redis", "database")
		database, err := strconv.Atoi(databaseS)
		if err != nil {
			database = 0
		}
		p.store = rds.NewStorager(dbconn, password, database)
		logs.Logger.Debug("----redis addr=", dbconn, " password:", password, " database:", database)
	}

	//开启统计输出
	go p.stat()

	logs.Logger.Debug("--------Init success--------")
}

func (p *Router) Start() {
	defer func() {
		if r := recover(); r != nil {
			logs.Logger.Critical("Start.recover:", r)
			go p.Start()
		}
	}()
	p.rpcServer = rpc.NewRpcServer(p.routerRpcAddr, p.maxRpcInFight, p.RpcSyncHandle, p.RpcAsyncHandle)

	//处理node异常中断 清除node以及node上注册的用户
	go func() {
		for {
			select {
			case id := <-p.nodeExit:
				p.pool.deleteNode(id)
				p.store.OfflineNode(id)
			}
		}
	}()

	p.startHttpServer()

	logs.Logger.Debug("--------Start Router success--------")
}

func (p *Router) Stop() error {
	debug.PrintStack()
	close(p.exit)
	return nil
}

//newRpcClient 返回一个RPC客户端
func (p *Router) NewRpcClient(name, addr string, ch chan int) (*rpc.RpcClient, error) {
	c, err := rpc.NewRpcClient(name, addr, ch)
	if err != nil {
		logs.Logger.Error("NewRpcClient ", err)
		return c, err
	}
	return c, err
}

func (p *Router) stat() {
	t := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-t.C:
			logs.Logger.Debug("Registered ", len(p.pool.nodes), " comets with ", p.store.SessionCount(), " sessions")
		}
	}
}
