package main

import (
	"fmt"
	"iot/internal/logs"
	"iot/internal/protocol"
	"iot/internal/rpc"
)

//RPC 异步句柄
func (p *Node) RpcAsyncHandle(request interface{}) {
	defer func() {
		if r := recover(); r != nil {
			logs.Logger.Error("recover ", r)
		}
	}()
	msg := request.(*rpc.PushRequst)
	logs.Logger.Info("Receive From Router type=", msg.Tp, " id=", msg.Id, " palt=", msg.Termtype, " msg=", msg.Msg)
	switch request.(type) {
	case *rpc.PushRequst:
		switch msg.Tp {
		case protocol.MSGTYPE_MESSAGE:
			p.message(msg)
		}
	}
}

//RPC 同步句柄
func (p *Node) RpcSyncHandle(request interface{}) int {
	defer func() {
		if r := recover(); r != nil {
			logs.Logger.Error("recover ", r)
		}
	}()
	switch request.(type) {
	case *rpc.KickRequst:
		{
			msg := request.(*rpc.KickRequst)
			logs.Logger.Info("Receive From Router Kick it=", msg.Id, " plat=", msg.Termtype, " reson=", msg.Reason, " token=", msg.Token)
			p.kick(msg)
		}
	}
	return -1
}


func (p *Node) message(msg *rpc.PushRequst) {
	id := msg.Id
	plat := msg.Termtype
	ptlType := protocol.PROTOCOL_TYPE_BINARY
	logs.Logger.Debug("[>>>MESSAGE]  id=", msg.Id, " plat=", plat)

	ids := fmt.Sprintf("%s-%d", id, plat)
	if sess := p.pool.findSessions(ids); sess != nil {
		var sendMsg protocol.ImDown
		protocol.SetMsgType(&sendMsg.Header, protocol.MSGTYPE_MESSAGE)
		protocol.SetEncode(&sendMsg.Header, sess.encode)
		sendMsg.Tid = uint32(sess.nextTid())
		sendMsg.Len = uint32(len(msg.Msg) + 1)
		sendMsg.Msg = append(sendMsg.Msg, []byte(msg.Msg)...)

		buf := protocol.Pack(&sendMsg, ptlType)
		logs.Logger.Debug("[MESSAGE>>>]  to id=", id, " plat=", plat, " Tid=", sendMsg.Tid)
		if err := p.write(sess.conn, buf); err != nil {
			logs.Logger.Error("MESSAGE>>> write error:", err)
		} else {
			//创建事务并保存
			trans := newTransaction()
			trans.tid = int(sendMsg.Tid)
			trans.msgType = protocol.MSGTYPE_MESSAGE
			//			trans.webOnline = iWebOnline
			trans.msg = append(trans.msg, []byte(msg.Msg)...) //mem leak
			sess.saveTransaction(trans)
			sess.checkTrans(trans)
		}
	} else {
		logs.Logger.Debug("[>>>MESSAGE]Not find session id=", msg.Id, " plat=", plat)
	}
}

func (p *Node) kick(msg *rpc.KickRequst) {
	id := msg.Id
	plat := msg.Termtype
	ptlType := protocol.PROTOCOL_TYPE_BINARY
	logs.Logger.Debug("[>>>KICK]  id=", msg.Id, " plat=", plat, " reason=", msg.Reason)

	ids := fmt.Sprintf("%s-%d", id, plat)
	if sess := p.pool.findSessions(ids); sess != nil {
		if sess.token != msg.Token {
			return
		}
		var sendMsg protocol.Kick
		protocol.SetMsgType(&sendMsg.Header, protocol.MSGTYPE_KICK)
		protocol.SetEncode(&sendMsg.Header, sess.encode)
		sendMsg.Tid = uint32(sess.nextTid())
		sendMsg.Len = 1
		sendMsg.Reason = uint8(msg.Reason)

		buf := protocol.Pack(&sendMsg, ptlType)
		logs.Logger.Debug("[KICK>>>]  to id=", id, " plat=", plat, " Tid=", sendMsg.Tid, " reason=", sendMsg.Reason)
		if err := p.write(sess.conn, buf); err != nil {
			logs.Logger.Error("KICK>>> write error:", err)
		}

		//上报router
		p.rpcCli.Notify(sess.id, sess.plat, sess.token, rpc.STATE_OFFLINE, p.nodeId)
		//清除session
		sess.destroy()
		p.pool.deleteSessions(ids)
		p.pool.deleteIds(connString(sess.conn))
		//关闭连接
		p.closeConn(sess.conn)
	}
}
