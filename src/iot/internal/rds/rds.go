package rds

import (
	"encoding/json"
	//"iot/internal/logs"
	"sync"

	"github.com/hoisie/redis"
)

//存储器规则
//查找：优先查找内存存储，再查询redis存储
//写入：优先写入redis存储，再写内存存储
type Storager struct {
	//redis存储
	cli *redis.Client

	//内存存储
	mutex    sync.Mutex
	memStore map[string]*Sessions
}

func NewStorager(dial, pswd string, db int) *Storager {
	var store Storager
	var client redis.Client
	client.Addr = dial
	//client.Password = pswd
	client.Db = db
	client.MaxPoolSize = 10
	store.cli = &client
	store.memStore = make(map[string]*Sessions)
	return &store
}

//FindSessions 查找session组合 如果未找到则返回nil
func (p *Storager) FindSessions(id string) (*Sessions, error) {
	//内存拷贝出去 防止多线程操作失败
	var sess Sessions

	//先查询内存
	//p.mutex.Lock()
	//s, ok := p.memStore[id]
	//p.mutex.Unlock()
	//if ok {
	//	sess = *s
	//	return &sess, nil
	//}

	//如未找到则查询redis
	//var sess Sessions
	b, err := p.cli.Get(id)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &sess)
	if err != nil {
		return nil, err
	}
	//p.memStore[id] = &sess
	return &sess, nil
}

//SaveSessions 保持session组合
func (p *Storager) SaveSessions(sess *Sessions) error {
	b, err := json.Marshal(sess)
	if err != nil {
		//logs.Logger.Error("Marshal err:", err, " b=", string(b))
		return err
	}
	err = p.cli.Set(sess.Id, b)
	if err != nil {
		//logs.Logger.Error("Redis err:", err)
		return err
	}
	p.mutex.Lock()
	p.memStore[sess.Id] = sess
	p.mutex.Unlock()
	return nil
}

//SessionOnline 返回指定用户是否在线
func (p *Storager) SessionOnline(id string, plat int) bool {
	var online bool
	sess, _ := p.FindSessions(id)
	if sess == nil {
		return online
	}
	for _, v := range sess.Sess {
		if v.Plat == plat && v.Online == true {
			online = true
			return online
		}
	}
	return online
}

//SessionNode 返回指定用户所连接Node
func (p *Storager) SessionNode(id string) string {
	sess, _ := p.FindSessions(id)
	if sess != nil {
		return sess.NodeId
	}
	return ""
}

func (p *Storager) SessionCount() int {
	return len(p.memStore)
}

func (p *Storager) OfflineNode(node string) {
	p.mutex.Lock()
	for _, sess := range p.memStore {
		if sess.NodeId == node {
			for _, it := range sess.Sess {
				it.Online = false
			}
			b, err := json.Marshal(sess)
			if err != nil {
				//logs.Logger.Error("Marshal err:", err, " b=", string(b))
			}
			if err := p.cli.Set(sess.Id, b); err != nil {
				//logs.Logger.Error("Redis err:", err)
			}
		}
	}
	p.mutex.Unlock()
}