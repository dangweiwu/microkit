package main

import (
	"fmt"
	"os"

	"github.com/dangweiwu/microkit/yamlconfig"
)

type Config struct {
	Host        string `yaml:"Host"`
	User        string `yaml:"User" default:"user"`
	Password    string `yaml:"Password" validate:"required"`
	Role        string `yaml:"Role" validate:"oneof=admin publisher author"`
	SuperUser   string `yaml:"SuperUser"`
	EnvPassword string `yaml:"EnvPassword"`
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

	cfg = &Config{}
	os.Setenv("ENV_PASSWORD", "password123456")
	yamlconfig.MustLoad("config4.yaml", cfg)
	fmt.Println("配置文件环境变量替换", cfg)
}
