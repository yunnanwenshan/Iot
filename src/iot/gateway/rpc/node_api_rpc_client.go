package rpc

import (
	"iot/gateway/logger"
	"github.com/polaris1119/config"
	"iot/internal/rpc"
)


var (
	Nodeapi *NodeApi
)

type NodeApi struct {
			       //与router交互
	nodeId       string //node api id 手工配置
	routerRpcAddr string //router RPC服务地址
	rpcCli        *rpc.RpcClient
	rpcStateChan  chan int //RPC链接状态通知通道
}

func init() {
	Nodeapi := new(NodeApi)
	Nodeapi.Init()
	Nodeapi.NewClientRpc()
	Nodeapi.rpcStateChan <- 1
}

func (node *NodeApi) Init() {
	log := logger.GetLoggerInstance()
	configRouter, err := config.ConfigFile.GetSection("router")
	if err != nil {
		log.Errorf("config init fail, err: %s", err.Error())
		panic("config init fail")
	}
	node.nodeId = configRouter["node_id"]
	node.routerRpcAddr = configRouter["rpc_addr"]
	node.rpcStateChan = make(chan int, 1)
}

func (node *NodeApi) NewClientRpc() {
	log := logger.GetLoggerInstance()
	// rpc client to router
	{
		log.Debug("node api start rpc connected to router rpc server")
		client, err := rpc.NewRpcClient(node.nodeId, node.routerRpcAddr, node.rpcStateChan)
		if err != nil {
			log.Errorf("connect to router fail, addr: %s", node.routerRpcAddr)
			panic(err)
		}
		node.rpcCli = client
		go node.CheckRpc()
	}
}

func (node *NodeApi) CheckRpc() {
	go func(ch chan int) {
		log := logger.GetLoggerInstance()
		for {
			select {
			case i := <- ch:
				switch i {
				case 0:
					{
						err := node.rpcCli.ReConnect()
						if err != nil {
							log.Errorf("reconnected to the router failed")
							return
						}
					}
				case 1:
					{
						node.rpcCli.StartPing()
						log.Debug("node api started, ping the router>>>>>>>>")
					}
				}
			}
		}
	}(node.rpcStateChan)
}

func (node *NodeApi) GetNode() (int, error) {
	return node.rpcCli.GetNode()
}