package redis

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
	"github.com/polaris1119/config"
	"fmt"
	"strings"
	"iot/gateway/logger"
	"iot/internal/rds"
)

var (
	redisInstance *ConnPool
	redisForDeviceInstance *rds.Storager
)

// ConnPool is for Cache fd
type ConnPool struct {
	redisPool *redis.Pool
}

// Init func create Cache fd by REDIS configration map:
//	var REDIS = map[string]string{
//		"host":         "127.0.0.1:6379",
//		"database":     "0",
//		"password":     "",
//		"maxOpenConns": "0",
//		"maxIdleConns": "0",
//	}

//获取实例
func RedisInstance() *ConnPool  {
	if redisInstance == nil {
		Init()
	}

	return redisInstance
}

//获取设备绑定实例
func RedisForDeviceInstance() *rds.Storager {
	if redisForDeviceInstance != nil {
		return redisForDeviceInstance
	}
	initRedisForDevice()

	return redisForDeviceInstance
}

//初始化redis pool
func Init() {
	var redisHost string

	logger := logger.GetLoggerInstance()
	redisConfig, err := config.ConfigFile.GetSection("redis")
	env, _ := config.ConfigFile.GetSection("global")
	//qconfConfig, _ := config.ConfigFile.GetSection("qconf")
	if err != nil {
		fmt.Printf("redis init fail, err = %v", err)
		logger.Errorf("redis init fail, err = %v", err)
		panic("redis init fail")
	}

	if strings.Compare(env["env"], "debug") == 0 {
		redisHostInfo, ok := redisConfig["host"]
		redisHost = redisHostInfo
		if !ok {
			logger.Errorf("get host fail, ok = %v", ok)
			panic("get host fail")
		}
	//} else {
		//redisHostInfo, err := qconf.GetHost(redisConfig["host"], qconfConfig["qconf"])
		//redisHost = redisHostInfo
		//if err != nil {
		//	logger.Errorf("qconf parse error, err = %v", err)
		//	panic("qconf parse error")
		//}
		//logger.Infof("qconf addr: %s, ip addr: %s", redisConfig["host"], redisHost)
	}

	params := map[string]string {
		"host": redisHost,
		"password": "",
		"maxOpenConns": "10000",
		"maxIdleConns": "50",
		"database": "0",
	}

	fmt.Printf("redis host info, %v", params)

	redisInstance = initPool(params)

	logger.Info("redis init successful")
	fmt.Println("redis init successful")
}

func initRedisForDevice() {
	var redisHost string
	logger := logger.GetLoggerInstance()
	redisConfig, err := config.ConfigFile.GetSection("redis")
	env, _ := config.ConfigFile.GetSection("global")
	//qconfConfig, _ := config.ConfigFile.GetSection("qconf")
	if err != nil {
		fmt.Printf("redis init fail, err = %v", err)
		logger.Errorf("redis init fail, err = %v", err)
		panic("redis init fail")
	}

	if strings.Compare(env["env"], "debug") == 0 {
		redisHostInfo, ok := redisConfig["host"]
		redisHost = redisHostInfo
		if !ok {
			logger.Errorf("get host fail, ok = %v", ok)
			panic("get host fail")
		}
	}

	redisForDeviceInstance = rds.NewStorager(redisHost, "", 0)
	logger.Infof("redis for device init successful")
}

func initPool(REDIS map[string]string) *ConnPool {
	Cache := &ConnPool{}
	maxOpenConns, _ := strconv.ParseInt(REDIS["maxOpenConns"], 10, 64)
	maxIdleConns, _ := strconv.ParseInt(REDIS["maxIdleConns"], 10, 64)
	database, _ := strconv.ParseInt(REDIS["database"], 10, 64)

	Cache.redisPool = newPool(REDIS["host"], REDIS["password"], int(database), int(maxOpenConns), int(maxIdleConns))
	if Cache.redisPool == nil {
		panic("init redis failed！")
	}
	return Cache
}

func newPool(server, password string, database, maxOpenConns, maxIdleConns int) *redis.Pool {
	return &redis.Pool{
		MaxActive:   maxOpenConns, // max number of connections
		MaxIdle:     maxIdleConns,
		IdleTimeout: 10 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			//if _, err := c.Do("AUTH", password); err != nil {
			//	c.Close()
			//	return nil, err
			//}
			//if _, err := c.Do("select", database); err != nil {
			//	c.Close()
			//	return nil, err
			//}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// Close pool
func (p *ConnPool) Close() error {
	err := p.redisPool.Close()
	return err
}

// Do commands
func (p *ConnPool) Do(command string, args ...interface{}) (interface{}, error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	return conn.Do(command, args...)
}

// SetString for string
func (p *ConnPool) SetString(key string, value interface{}) (interface{}, error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	return conn.Do("SET", key, value)
}

// GetString for string
func (p *ConnPool) GetString(key string) (string, error) {
	// get one connection from pool
	conn := p.redisPool.Get()
	// put connection to pool
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}

// GetBytes for bytes
func (p *ConnPool) GetBytes(key string) ([]byte, error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	return redis.Bytes(conn.Do("GET", key))
}

// GetInt for int
func (p *ConnPool) GetInt(key string) (int, error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("GET", key))
}

// GetInt64 for int64
func (p *ConnPool) GetInt64(key string) (int64, error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("GET", key))
}

// DelKey for key
func (p *ConnPool) DelKey(key string) (interface{}, error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	return conn.Do("DEL", key)
}

// ExpireKey for key
func (p *ConnPool) ExpireKey(key string, seconds int64) (interface{}, error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	return conn.Do("EXPIRE", key, seconds)
}

// Keys for key
func (p *ConnPool) Keys(pattern string) ([]string, error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do("KEYS", pattern))
}

// KeysByteSlices for key
func (p *ConnPool) KeysByteSlices(pattern string) ([][]byte, error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	return redis.ByteSlices(conn.Do("KEYS", pattern))
}

// SetHashMap for hash map
func (p *ConnPool) SetHashMap(key string, fieldValue map[string]interface{}) (interface{}, error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	return conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(fieldValue)...)
}

// GetHashMapString for hash map
func (p *ConnPool) GetHashMapString(key string) (map[string]string, error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	return redis.StringMap(conn.Do("HGETALL", key))
}

// GetHashMapInt for hash map
func (p *ConnPool) GetHashMapInt(key string) (map[string]int, error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	return redis.IntMap(conn.Do("HGETALL", key))
}

// GetHashMapInt64 for hash map
func (p *ConnPool) GetHashMapInt64(key string) (map[string]int64, error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	return redis.Int64Map(conn.Do("HGETALL", key))
}
