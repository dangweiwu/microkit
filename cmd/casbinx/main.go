package main

import (
	"fmt"
	"github.com/dangweiwu/microkit/casbinx"
	"github.com/dangweiwu/microkit/db/mysqlx"
	"log"
)

func main() {
	dbcli, err := mysqlx.NewClient(mysqlx.Config{
		DbName:   "goservice",
		User:     "root",
		Password: "a12346",
		Host:     "127.0.0.1:13306",
	})
	if err != nil {
		panic(err)
	}

	cs, err := casbinx.NewCasbinGorm(dbcli)
	if err != nil {
		panic(err)
	}

	ok, e := cs.Enforce("admin", "/api/v1/user/list", "get")
	log.Println("enforce", ok, e)

	fmt.Println(cs.AddPolicy("admin", []casbinx.ApiPolice{
		{Api: "/api/v1/user/list", Method: "get"},
		{Api: "/api/v1/user/list", Method: "post"},
	}))

	ok, e = cs.Enforce("admin", "/api/v1/user/list", "get")
	log.Println("enforce", ok, e)

	//time.Sleep(time.Second)
	fmt.Println(cs.AddPolicy("user1", []casbinx.ApiPolice{
		{Api: "/api/v1/user/list", Method: "get"},
		{Api: "/api/v1/user/list", Method: "post"},
	}))
	cs.RemoveRolePolicy("admin")
	ok, e = cs.Enforce("admin", "/api/v1/user/list", "get")
	log.Println("enforce", ok, e)
	//
	//cs.Clear()
	//ok, e = cs.Enforce("admin", "/api/v1/user/list", "get")
	//log.Println("enforce", ok, e)
	//cs.Casbin.LoadPolicy()
	//ok, e = cs.Enforce("admin", "/api/v1/user/list", "get")
	//log.Println("enforce", ok, e)

}
