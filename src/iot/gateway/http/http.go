package http

import (
	"github.com/boj/redistore"
	"github.com/polaris1119/config"
	"iot/gateway/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gorilla/sessions"
	"fmt"
	"strings"
	"iot/gateway/util/qconf"
)

type StoreRedis struct {
	store *redistore.RediStore
}

var Store = new(StoreRedis)

func init()  {
	var redisHost string
	logger := logger.GetLoggerInstance()
	redisConfig, err := config.ConfigFile.GetSection("redis")
	env, _ := config.ConfigFile.GetSection("global")
	qconfConfig, _ := config.ConfigFile.GetSection("qconf")
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
	} else {
		redisHostInfo, err := qconf.GetHost(redisConfig["host"], qconfConfig["qconf"])
		redisHost = redisHostInfo
		if err != nil {
			logger.Errorf("qconf parse error, err = %v", err)
			panic("qconf parse error")
		}

		logger.Infof("qconf addr: %s, ip addr: %s", redisConfig["host"], redisHost)
	}
	logger.Infof("redis init, url: %v", redisHost)
	st, err := redistore.NewRediStore(10, "tcp", redisHost, "", []byte(config.ConfigFile.MustValue("global", "cookie_secret")))
	if err != nil {
		logger.Warnf("new redis store init fail, err: %s, %s", err.Error(), redisHost)
		panic("redis初始化失败")
	}
	logger.Infof("redis初始化成功, host: %s", redisHost)
	Store.store = st
}

func GetStore() *redistore.RediStore {
	return Store.store
}

func Request(ctx *gin.Context) *http.Request {
	return ctx.Request
}

func ResponseWriter(ctx *gin.Context) http.ResponseWriter  {
	return ctx.Writer
}

func GetCookieSession(ctx *gin.Context) *sessions.Session  {
	session, _ := Store.store.Get(Request(ctx), "user")
	return session
}

func SetCookie(ctx *gin.Context, userName string)  {
	Store.store.Options.HttpOnly = true

	session := GetCookieSession(ctx)
	if ctx.PostForm("remember_me") != "1" {
		session.Options = &sessions.Options{
			Path: "/",
			HttpOnly: true,
		}
	}
	session.Values["username"] = userName
	req := Request(ctx)
	res := ctx.Writer
	session.Save(req, res)
}

