package rpc

import (
	"iot/gateway/logger"
	"github.com/polaris1119/config"
	"iot/internal/rpc"
)


var nodeApiRpcClient *NodeApi = nil

//与router交互
type NodeApi struct {
	nodeId       string //node api id 手工配置
	routerRpcAddr string //router RPC服务地址
	rpcClient        *rpc.RpcClient
	rpcStateChan  chan int //RPC链接状态通知通道
}

func GetNodeApiInstance() (*NodeApi) {
	log := logger.GetLoggerInstance()
	if nodeApiRpcClient == nil {
		nodeApiRpcClient = new(NodeApi)
		nodeApiRpcClient.Init()
		nodeApiRpcClient.NewClientRpc()
		nodeApiRpcClient.rpcStateChan <- 1
		log.Debugf("node api init successful")
	}

	return nodeApiRpcClient
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
		log.Debug("-------------node api start rpc connected to router rpc server---------------")
		log.Debugf("---------rpc server:%s, nodeId: %s", node.routerRpcAddr, node.nodeId)
		client, err := rpc.NewRpcClient(node.nodeId, node.routerRpcAddr, node.rpcStateChan)
		if err != nil {
			log.Errorf("connect to router fail, addr: %s", node.routerRpcAddr)
			panic(err)
		}
		node.rpcClient = client
		go node.CheckRpc()
		log.Debugf("-------------node api end rpc connected to router rpc server-----------------")
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
						err := node.rpcClient.ReConnect()
						if err != nil {
							log.Errorf("reconnected to the router failed")
							return
						}
					}
				case 1:
					{
						node.rpcClient.StartPing()
						log.Debug("node api started, ping the router>>>>>>>>")
					}
				}
			}
		}
	}(node.rpcStateChan)
}

func (node *NodeApi) GetNode(uid string) (int, error) {
	log := logger.GetLoggerInstance()
	log.Debug("GetNode begin >>>>>>")

	nodeId, err := node.rpcClient.GetNodeByUid(uid)
	log.Debugf("getNode end >>>>, nodeId: %d, err :%v", nodeId, err)
	return nodeId, err
}