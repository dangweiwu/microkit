package main

import (
	"fmt"
	"microkit/yamlconfig"
)

type Config struct {
	Host     string
	User     string `json:"user" default:"user"`
	Password string `validate:"empty=false"`
	Role     string `validate:"one_of=admin,publisher,author"`
}

func main() {
	cfg := &Config{}

	func() {
		defer func() {
			e := recover()
			fmt.Println(e)
		}()

		yamlconfig.MustLoad("config1.yaml", cfg)
	}()

	func() {
		defer func() {
			e := recover()
			fmt.Println(e)
		}()

		yamlconfig.MustLoad("config2.yaml", cfg)
	}()
	yamlconfig.MustLoad("config3.yaml", cfg)
	fmt.Println(cfg)
}
