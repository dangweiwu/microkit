package mysqlx

type Config struct {
	User     string `yaml:"user" validate:"required" default:"root"`
	Password string `yaml:"password" validate:"required" default:"123456"`
	Host     string `yaml:"host" default:"localhost:3306"`
	DbName   string `yaml:"dbName" validate:"required"`
	LogFile  string `yaml:"logFile"`                           //日志位置
	LogLevel int    `yaml:"logLevel" validate:"oneof=1 2 3 4"` //1 slice 2 err 3 war 4 info
}
