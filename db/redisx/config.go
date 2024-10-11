package redisx

type Config struct {
	Addr     string `validate:"empty=false"`
	Password string `validate:""`
	Db       int    `default:"0"`
}
