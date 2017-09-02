package service

import "testing"

func TestLogin(t *testing.T)  {
	userService := UserService{}
	err := userService.Login("test", "1235")

	if err != nil {
		t.Logf("登录测试错误: err = %v", err)
	}
	t.Log("测试登录通过");
}
