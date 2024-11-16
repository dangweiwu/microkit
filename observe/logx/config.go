package logx

const (
	//level
	ERROR = "error"
	DEBUG = "debug"
	PANIC = "panic"
	INFO  = "info"
	WARN  = "warn"

	//输出类型
	JSON = "json"
	TXT  = "txt"

	//输出位置
	CONSOLE = "console"
	FILE    = "file"
	ALL     = "all"
)

type Config struct {
	//分割配置
	LogName      string `yaml:"LogName" default:"log.log"` // 日志文件路径，默认 os.TempDir()
	MaxSize      int    `yaml:"MaxSize" default:"10"`      // 每个日志文件保存10M，默认 100M
	MaxBackNum   int    `yaml:"MaxBackNum" default:"15" `  // 保留30个备份，默认不限
	MaxAge       int    `yaml:"MaxAge" default:"7"`        // 保留7天，默认不限
	Compress     bool   `yaml:"Compress" default:"false"`  // 是否压缩，默认不压缩
	Level        string `yaml:"Level" default:"debug" validate:"oneof=error debug panic info warn"`
	OutType      string `yaml:"OutType" default:"console" validate:"oneof=console file all"` // 输出到哪 console all file
	Formatter    string `yaml:"Formatter" default:"txt" validate:"oneof=txt json"`           //json or txt
	HasTimestamp bool   `yaml:"HasTimestamp" default:"false"`
	Caller       bool   `yaml:"Caller" default:"false"`      //启用堆栈
	Development  bool   `yaml:"Development" default:"false"` // 记录行号
}
