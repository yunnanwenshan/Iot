package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func success(ctx *gin.Context, data interface{}) error {
	//增加跨域头
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type")
	ctx.Header("content-type", "application/json")

	ctx.JSON(http.StatusOK, data)

	return nil
}

func fail(ctx *gin.Context, code int, msg string) error {
	result := map[string]interface{}{
		"error": map[string]interface{}{
			"code" : code,
			"message" : msg,
			"data": "",
		},
	}

	//增加跨域头
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type")
	ctx.Header("content-type", "application/json")

	ctx.JSON(http.StatusOK, result)

	return nil
}