package serverToken

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
	"iot/gateway/logger"
)

type MyCustomClaims struct {
	userId string `json:"user_id"`
	env string `json:"env"`
	ttu string `json:"ttu"`
	jwt.StandardClaims
}

var (
	tokenSig []byte
)

func init()  {
	tokenSig = []byte("%cSMkR8gbr8y&3L03AiF&7D2W3@KfAlgi%")
}

func GenerateToken(userId string) (jwtToken string, err error) {
	logger := logger.GetLoggerInstance()
	// Create the Claims
	claims := MyCustomClaims{
		userId,
		"test",
		"token",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(365).Unix(),
			Issuer:    "test",
			Audience: "ios",
			IssuedAt: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken, err = token.SignedString(tokenSig)
	logger.Infof("generate jwt: %s, userid: %s", jwtToken, userId)
	return
}

func ExtractUserId(serverToken string) (userId string, err error) {
	logger := logger.GetLoggerInstance()
	token, err1 := jwt.ParseWithClaims(serverToken, &MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return tokenSig, nil
		})
	if err1 != nil {
		err = err1
		return "", err
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		logger.Infof("userid: %s,  %d", claims.userId, claims.StandardClaims.ExpiresAt)
		return  claims.userId, nil
	}

	return "", errors.New("parse token error")
}

func VlidateServerToken(serverToken string) (bool, error) {
	if len(serverToken) <= 0 {
		return false, errors.New("Sever-Token is null")
	}

	token, err := jwt.ParseWithClaims(serverToken, &MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return tokenSig, nil
		})
	if err != nil {
		return false, err
	}

	if token.Valid {
		return true, nil
	}

	return false, errors.New("Server-token verify fail")
}
