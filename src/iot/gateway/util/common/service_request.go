package common

import (
	"github.com/parnurzeal/gorequest"
	"iot/gateway/logger"
	"time"
	"github.com/bitly/go-simplejson"
	"github.com/polaris1119/config"
	"math/rand"
	"strconv"
	"crypto/sha1"
	"fmt"
	"strings"
)

func HttpRequest(uri string, params interface{}) (res *simplejson.Json, err error)  {
	//存在并发访问config.ConfigFile出现问题: fatal error: concurrent map read and map write
	lock.Lock()
	defer lock.Unlock()

	begin := time.Now().UnixNano()
	logger := logger.GetLoggerInstance()
	request := gorequest.New()

	//_, file, line, ok := runtime.Caller(0)
	//if ok {
	//	logger.Infof("request before, uri: %s, params: %v, file: %s, line: %d", uri, params, file, line)
	//}

	config, _ := config.ConfigFile.GetSection("global")
	logger.Infof("request before, uri: %s, params: %v, domain: %s", uri, params, config["domain"])

	resp, body, errs := request.Post(uri).
		Set("Content-Type", "application/json").
		Set("Origin", config["domain"]).
		Send(params).
		End()

	logger.Infof("request after, uri: %s, params: %v, body: %v", uri, params, body)

	//计算接口所花费的时间
	end := time.Now().UnixNano()
	diff := (end - begin) / (1000 * 1000)
	logger.Infof("Request service, uri: %s, cost time: %v ms", uri, diff)

	if errs != nil {
		logger.Infof("call interface fail, body: %v, resp: %v, error: %v", body, resp, errs)
		return nil, errs[0]
	}

	return simplejson.NewJson([]byte(body))
}

func HttpRequestEngine(uri string, params interface{}) (res *simplejson.Json, err error)  {
	begin := time.Now().UnixNano()
	logger := logger.GetLoggerInstance()
	request := gorequest.New()

	serviceToken, err := generateToken()
	if err != nil {
		logger.Errorf("generate service token, err = %v", err.Error())
		return nil, err
	}

	logger.Infof("request before, uri: %s, params: %v, serviceToken: %s", uri, params, serviceToken)

	resp, body, errs := request.Post(uri).
		Set("Content-Type", "application/json").
		Set("Service-Token", serviceToken).
		Send(params).
		End()

	logger.Infof("request after, uri: %s, params: %v, body: %v", uri, params, body)

	//计算接口所花费的时间
	end := time.Now().UnixNano()
	diff := (end - begin) / (1000 * 1000)
	logger.Infof("Request service, uri: %s, cost time: %v ms", uri, diff)

	if errs != nil {
		logger.Infof("call interface fail, body: %v, resp: %v, error: %v", body, resp, errs)
		return nil, errs[0]
	}

	return simplejson.NewJson([]byte(body))
}

//内部服务Service-token的产生
func generateToken() (string, error) {
	//存在并发访问config.ConfigFile出现问题: fatal error: concurrent map read and map write
	lock.Lock()
	defer lock.Unlock()

	logger := logger.GetLoggerInstance()
	engineConfig, err := config.ConfigFile.GetSection("engine")
	if err != nil {
		logger.Errorf("get engine key fail, err = %s", err.Error())
		return "", err
	}
	key := engineConfig["engine_key"]
	serviceName := engineConfig["engine_name"]
	randInt := rand.Intn(10000)
	payload := serviceName + "." + strconv.Itoa(randInt)

	//sha1加密
	h := sha1.New()
	h.Write([]byte(payload + key))
	bs := h.Sum(nil)

	return payload + "." + fmt.Sprintf("%x", bs), nil
}

//Service-token的验证
func verifyToken(token string) (error) {
	//存在并发访问config.ConfigFile出现问题: fatal error: concurrent map read and map write
	lock.Lock()
	defer lock.Unlock()

	logger := logger.GetLoggerInstance()
	engineConfig, err := config.ConfigFile.GetSection("engine")
	if err != nil {
		logger.Errorf("get engine key fail, err = %s", err.Error())
		return err
	}
	key := engineConfig["engine_key"]
	//serviceName := engineConfig["engine_name"]
	tokenArr := strings.Split(token, ".")
	h := sha1.New()
	h.Write([]byte(tokenArr[0] + key))
	bs := h.Sum(nil)
	sign := fmt.Sprintf("%x", bs)

	if !strings.EqualFold(sign, tokenArr[1]) {
		logger.Errorf("service token verify fail")
		return nil
	}

	return nil
}
