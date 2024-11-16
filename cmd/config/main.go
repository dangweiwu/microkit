package main

import (
	"fmt"
	"github.com/dangweiwu/microkit/yamlconfig"
)

type Config struct {
	Host      string `yaml:"Host"`
	User      string `yaml:"User" default:"user"`
	Password  string `yaml:"Password" validate:"required"`
	Role      string `yaml:"Role" validate:"oneof=admin publisher author"`
	SuperUser string `yaml:"SuperUser"`
}

func main() {
	cfg := &Config{}
	//
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
		cfg = &Config{}

		yamlconfig.MustLoad("config2.yaml", cfg)
	}()

	cfg = &Config{}

	yamlconfig.MustLoad("config3.yaml", cfg)
	fmt.Println(cfg)
}
