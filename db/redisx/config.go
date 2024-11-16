package redisx

type Config struct {
	Addr     string `yaml:"Addr" validate:"required"`
	Password string `yaml:"Password"`
	Db       int    `default:"0"`
}
