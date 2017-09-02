package common

import (
	"github.com/polaris1119/config"
	"iot/gateway/util/qconf"
	"strings"
)

func getHost(path string, serviceName string) (host string, err error)  {
	//存在并发访问config.ConfigFile出现问题: fatal error: concurrent map read and map write
	lock.Lock()
	defer lock.Unlock()

	host = ""
	env, _ := config.ConfigFile.GetSection("global")
	envEngine, _ := config.ConfigFile.GetSection(serviceName)
	host = envEngine["host"]

	//开发环境使用的是ip地址，无法使用qconf
	if strings.Compare(env["env"], "debug") == 0 {
		return "http://" + host + path, nil
	}

	qconfConfig, _ := config.ConfigFile.GetSection("qconf")
	host, err = qconf.GetHost(host, qconfConfig["qconfIdc"]);
	if err != nil {
		return
	}

	host = "http://" + host + path

	return
}

//engine服务qconf地址到ip地址转换
func GetEngineHost(path string) (host string, err error) {
	serviceName := "engine"
	return getHost(path, serviceName)
}