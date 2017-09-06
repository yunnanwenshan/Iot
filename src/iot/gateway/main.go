package main

import (
	_ "iot/gateway/http"
	_ "iot/gateway/db"
	_ "iot/gateway/redis"

	"github.com/polaris1119/config"
	"github.com/gin-gonic/gin"
	"github.com/fvbock/endless"

	"fmt"
	"os"
	"strconv"
	"runtime"
	_ "net/http/pprof"
	"iot/gateway/logger"
	"iot/gateway/http/middleware"
	"iot/gateway/http/routes"
)

func main() {
	//获取进程号
	GetPid()

	//配置并发线程数
	ConfigRuntime()

	//设置环境
	config.ConfigFile.BlockMode = true;
	env, _ := config.ConfigFile.GetSection("global")
	gin.SetMode(env["env"])
	fmt.Printf("\033[32m[INFO]\033[0m env=%s\n", env["env"])

	//监控服务
	go func() {
		logger := logger.GetLoggerInstance()
		logger.Infof("start monitor goroutine...., return = %v", endless.ListenAndServe("localhost:6060", nil))
	}()

	r := gin.New()

	//设置中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.MonitorApi())
	//r.Use(middleware.SignatureMiddleWare())

	//静态文件
	serveStatic(r)

	//注册路由
	routes.RegisterRouters(r)

	// Listen and Server in 0.0.0.0:8080
	endless.ListenAndServe(":40001", r)
	//err := endless.ListenAndServe(":8080", r)
	//if err != nil {
	//	fmt.Println("start server fail....")
	//}
	//r.Run(":8080")
}

type staticRootConf struct {
	root   string
	isFile bool
}

var staticFileMap = map[string]staticRootConf{
	"/static/":     {"/static", false},
	"/favicon.ico": {"/static/img/go.ico", true},
}

var filterPrefixs = make([]string, 0, 3)

func serveStatic(e *gin.Engine) {
	for prefix, rootConf := range staticFileMap {
		filterPrefixs = append(filterPrefixs, prefix)

		if rootConf.isFile {
			e.StaticFile(prefix, config.ROOT+rootConf.root)
		} else {
			e.Static(prefix, config.ROOT+rootConf.root)
		}
	}
}

//获取进程id
func GetPid() {
	logger := logger.GetLoggerInstance()
	pid := os.Getpid()
	file, err := os.OpenFile("app.pid", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	if err != nil {
		logger.Errorf("get process ID error, err: %s", err.Error())
		os.Exit(1)
	}
	file.WriteString(strconv.Itoa(pid))

	logger.Infof("write process id successuful, pid = %d", pid)
}

//配置CPU个数
func ConfigRuntime() {
	logger := logger.GetLoggerInstance()
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	logger.Infof("Running with %d CPUs\n", nuCPU)
}
