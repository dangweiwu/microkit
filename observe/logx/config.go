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
	LogName      string `yaml:"logName" default:"log.log"` // 日志文件路径，默认 os.TempDir()
	MaxSize      int    `yaml:"maxSize" default:"10"`      // 每个日志文件保存10M，默认 100M
	MaxBackNum   int    `yaml:"maxBackNum" default:"15" `  // 保留30个备份，默认不限
	MaxAge       int    `yaml:"maxAge" default:"7"`        // 保留7天，默认不限
	Compress     bool   `yaml:"compress" default:"false"`  // 是否压缩，默认不压缩
	Level        string `yaml:"level" default:"debug" validate:"oneof=error debug panic info warn"`
	OutType      string `yaml:"outType" default:"console" validate:"oneof=console file all"` // 输出到哪 console all file
	Formatter    string `yaml:"formatter" default:"txt" validate:"oneof=txt json"`           //json or txt
	HasTimestamp bool   `yaml:"hasTimestamp" default:"false"`
	Caller       bool   `yaml:"caller" default:"false"`      //启用堆栈
	Development  bool   `yaml:"development" default:"false"` // 记录行号
}
