package mysqlx

type Config struct {
	User     string `yaml:"User" validate:"required" default:"root"`
	Password string `yaml:"Password" validate:"required" default:"123456"`
	Host     string `yaml:"Host" default:"localhost:3306"`
	DbName   string `yaml:"DbName" validate:"required"`
	LogFile  string `yaml:"LogFile"`                           //日志位置
	LogLevel int    `yaml:"LogLevel" validate:"oneof=1 2 3 4"` //1 slice 2 err 3 war 4 info
}
