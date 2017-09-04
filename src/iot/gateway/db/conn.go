// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package db

import (
	_ "database/sql"
	"errors"
	"fmt"

	. "github.com/polaris1119/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"iot/gateway/logger"
	"os"
)

var (
	ConnectDBErr = errors.New("connect db error")
	UseDBErr     = errors.New("use db error")

	sqlLogFile = "sql.log"
	masterDB *xorm.Engine
	dns string
)

func init() {
}

func fillDns(mysqlConfig map[string]string) {
	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConfig["user"],
		mysqlConfig["password"],
		mysqlConfig["host"],
		mysqlConfig["port"],
		mysqlConfig["dbname"],
		mysqlConfig["charset"])
}

func initEngine() error {
	var err error

	logger := logger.GetLoggerInstance()
	masterDB, err = xorm.NewEngine("mysql", dns)
	if err != nil {
		return err
	}

	maxIdle := ConfigFile.MustInt("mysql", "max_idle", 2)
	maxConn := ConfigFile.MustInt("mysql", "max_conn", 10)

	masterDB.SetMaxIdleConns(maxIdle)
	masterDB.SetMaxOpenConns(maxConn)

	showSQL := ConfigFile.MustBool("xorm", "show_sql", false)
	logLevel := ConfigFile.MustInt("xorm", "log_level", 1)

	file := sqlLogger()
	sqlLogFile := xorm.NewSimpleLogger(file)
	masterDB.SetLogger(sqlLogFile)
	masterDB.ShowSQL(showSQL)
	masterDB.Logger().SetLevel(core.LogLevel(logLevel))
	fmt.Println("database init finished")
	logger.Info("database init finished")
	// 启用缓存
	// cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	// MasterDB.SetDefaultCacher(cacher)

	return nil
}

func StdMasterDB() *xorm.Engine {
	if masterDB == nil {
		mysqlConfig, err := ConfigFile.GetSection("mysql")
		if err != nil {
			fmt.Println("get mysql config error:", err)
			panic("mysql init fail")
		}

		fillDns(mysqlConfig)

		// 启动时就打开数据库连接
		if err = initEngine(); err != nil {
			panic(err)
		}
	}

	return masterDB
}

func sqlLogger() *os.File {
	file, err := os.OpenFile(sqlLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("init sql.log file fail, err = %v", err)
		return nil;
	}

	return file
}
