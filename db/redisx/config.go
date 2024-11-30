package redisx

type Config struct {
	Addr     string `yaml:"addr" validate:"required"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db" default:"0"`
}
