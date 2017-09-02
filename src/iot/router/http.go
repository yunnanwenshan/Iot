package main

import (
	"fmt"
	"iot/internal/logs"
	"net/http"
	_ "net/http/pprof"
)

func (p *Router) startHttpServer() {
	defer func() {
		if r := recover(); r != nil {
			logs.Logger.Error("startHttpServer.recover:", r)
			go p.startHttpServer()
		}
	}()
	go func() {
		http.HandleFunc("/v1/gComet.addr", func(w http.ResponseWriter, r *http.Request) {
			p.loadDispatcher(w, r)
		})
		err := http.ListenAndServe(p.httpBindAddr, nil)
		if err != nil {
			panic(err)
		}
	}()
}

func (p *Router) loadDispatcher(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method != "POST" {
		return
	}

	id := r.FormValue("id")
	addr := p.balancer(id)
	logs.Logger.Debug("[http] load dispatcher id=", id, "addr=", addr, " remote=", r.RemoteAddr)

	//for ajax cross domain
	Origin := r.Header.Get("Origin")
	if Origin != "" {
		w.Header().Add("Access-Control-Allow-Origin", Origin)
		w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "x-requested-with,content-type")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
	}

	fmt.Fprintf(w, addr)
}

//balancer HTTP接口 根据id返回socket地址
func (p *Router) balancer(id string) string {
	sess := p.store.FindSessions(id)
	if sess != nil {
		c := p.pool.findComet(sess.CometId)
		if c != nil {
			return c.tcpAddr
		}
	}
	//系统指配
	c := p.pool.balancer()
	if c != nil {
		return c.tcpAddr
	}
	return ""
}
