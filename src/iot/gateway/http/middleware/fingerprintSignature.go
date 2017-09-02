package middleware

import (
	"github.com/gin-gonic/gin"
	"encoding/json"
	"fmt"
	"strings"
	"crypto"
	"encoding/base64"
	"github.com/pkg/errors"
	"iot/gateway/util/signature"
	"iot/gateway/logger"
	"net/http"
)

var (
	clientInfoKey = "clientInfo"
	clientSigKey = "App-Signature"
	defaulCode = 10000;
)

func SignatureMiddleWare() gin.HandlerFunc{
	return verify()
}

func verify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientSigStr := getHeader(ctx, clientSigKey)
		clientInfoStr := getHeader(ctx, clientSigKey)
		err := verifySig(ctx, clientSigStr, clientInfoStr)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": defaulCode + 8, "msg": "App-Signature verfiy fail"})
			ctx.Abort()
		}

		//继续处理下一个middleware
		ctx.Next()
	}
}

func getHeader(ctx *gin.Context, key string) string {
	return ctx.Request.Header.Get(key)
}

func verifySig(ctx *gin.Context, sigStr string, clientInfo string) (err error) {
	logger := logger.GetLoggerInstance()
	if len(sigStr) <= 0 {
		fmt.Println("App-Signature is null")
		return errors.New("App-Signature is null")
	}
	logger.Infof("App-Signature=%s", sigStr)

	elementArray := strings.SplitN(sigStr, ".", 2)
	logger.Infof("elementArray=%v \n", elementArray)
	logger.Infof("elementArray[0]=%v \n", elementArray[0])
	logger.Infof("elementArray[1]=%v \n", elementArray[1])
	logger.Infof("elementArray len=%v \n", len(elementArray[1]))

	if len(elementArray) != 2 {
		err = errors.New("invalid app signature")
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(elementArray[1]), &result)
	if err != nil {
		logger.Infof("parse json error, %v", err)
		return
	}
	logger.Infof("result=%v \n", result)

	//验证请求路径是否正确
	path := ctx.Request.URL.Path
	if uri, ok := result["uri"].(string); ok {
		if strings.Compare(path, uri) != 0 {
			logger.Infof("error request path, %v , %s\n",  result["uri"], path)
			err = errors.New("request path incrroect")
			return
		}
	}

	//这里go解析json字符串会对字段的顺序进行重调整，导致Marshal后的json字段顺序与客户端发送来的字段顺序不一致，从而导致md5值不一样
	//jti, ok := result["jti"].(string)
	//if !ok {
	//	logger.Printf("jit error, %v \n", result["jti"])
	//	err = errors.New("jit error")
	//	return
	//}
	//
	//var tlStr string
	//var body map[string]interface{}
	//err3 := ctx.BindJSON(&body)
	//if err3 != nil {
	//	err = err3
	//	logger.Printf("err = %v \n", err)
	//	return
	//}
	//body["mobile"] = "10000000000"
	//body["code"] = "1234"
	//jsonByte, _ := json.Marshal(body)
	//logger.Printf("============, %v , jsonByte = %s\n", body, string(jsonByte))
	//
	//ttTrim := base64urlEncode(jsonByte)
	////ttTrim := base64urlEncode([]byte(`{"mobile":"10000000000","code":1234}`))
	//tlStr = ttTrim + jti + "041d3b28ef034472a9b7bbeb60c8d588"
	//logger.Printf("ttTrim = %s, tlStr=%s, jsonByte=%s\n", ttTrim, tlStr, string(jsonByte))
	//md := md5.Sum([]byte(tlStr))
	//md5str := fmt.Sprintf("%x", md)
	//logger.Printf("tlStr = %s, md5str=%v, md=%v \n", tlStr, md5str, md)
	//
	//newRdb := strings.ToLower(Substr(string(md5str), 0, 8))
	//rdb, ok := result["rbd"].(string)
	//if !ok {
	//	logger.Printf("rdb error, %v \n", result["rdb"])
	//	err = errors.New("rdb error")
	//	return
	//}
	//
	//if strings.Compare(newRdb, rdb) != 0 {
	//	err = errors.New("rdb compare fail")
	//	logger.Printf("new rdb: %v, rdb: %v \n", newRdb, rdb)
	//	return
	//}


	resStr, err3 := base64urlDecode(elementArray[0])
	if err3 != nil {
		logger.Infof("err3=%v", err3)
		return err3
	}

	logger.Infof("ttt1=%X, elementArray[0]=%v\n", resStr, elementArray[0])
	err = signature.CipherInstance.Verify([]byte(elementArray[1]), []byte(resStr), crypto.SHA256)
	if err != nil {
		logger.Infof("err4=%v \n", err)
		return err
	}

	logger.Infof("verify successful\n")

	return nil
}

func base64urlDecode(str string) ([]byte, error) {
	str1 := strings.Replace(strings.Replace(str, "-", "+", -1), "_", "/", -1)
	isEqualZero := len(str1) % 4
	padLen := 0
	if isEqualZero != 0 {
		padLen = (4 - len(str1) % 4)
	}
	padStr := strings.Repeat("=", padLen)
	rsStr := str1 + padStr
	strings.Join([]string{str1, padStr}, "")
	return base64.StdEncoding.DecodeString(rsStr)
}

func base64urlEncode(str []byte) string  {
	str1 := strings.TrimRight(strings.Replace(strings.Replace(string(str[:]), "+", "-", -1), "/", "_", -1), "=")
	return base64.StdEncoding.EncodeToString([]byte(str1))
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
