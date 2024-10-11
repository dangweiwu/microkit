package main

import (
	"gorm.io/gorm"
	"log"
	"microkit/db/mysqlx"
)

type Demo struct {
	gorm.Model
	Name string
	Age  int
}

func main() {
	config := mysqlx.Config{
		DbName:   "test",
		User:     "root",
		Password: "123456",
		Host:     ":0",
		LogLevel: 1,
		LogFile:  "db.log",
	}

	host, mysqlServer, err := mysqlx.FakeDb(config.DbName)
	if err != nil {
		log.Panicf("fakeDb失败:%v\n", err)
	}
	defer mysqlServer.Close()

	config.Host = host

	dbCli, err := mysqlx.NewClient(config)
	if err != nil {
		log.Panicf("创建client失败:%v\n", err)
	}

	var tables []string

	dbCli.Raw("SHOW TABLES").Scan(&tables)
	log.Println("tables:", tables)
	if err := dbCli.AutoMigrate(&Demo{}); err != nil {
		panic(err)
	}

	dbCli.Raw("SHOW TABLES").Scan(&tables)
	log.Println("tables:", tables)
}
