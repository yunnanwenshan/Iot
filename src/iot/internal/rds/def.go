package rds

type Sessions struct {
	Id      string     `json:"id"`
	NodeId string     `json:"node"`
	Sess    []*Session `json:"sess"`
}

type Session struct {
	Id       string `json:"id"`
	Plat     int    `json:"plat"`
	Online   bool   `json:"online"`
	Login    bool   `json:"login"`
	AuthCode string `json:"authcode"`
}