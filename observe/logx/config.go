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
	LogName      string `default:"log.log"` // 日志文件路径，默认 os.TempDir()
	MaxSize      int    `default:"10"`      // 每个日志文件保存10M，默认 100M
	MaxBackNum   int    `default:"15" `     // 保留30个备份，默认不限
	MaxAge       int    `default:"7"`       // 保留7天，默认不限
	Compress     bool   `default:"false"`   // 是否压缩，默认不压缩
	Level        string `default:"debug" validate:"one_of=error,debug,panic,info,warn"`
	OutType      string `default:"console" validate:"one_of=console,file,all"` // 输出到哪 console all file
	Formatter    string `default:"txt" validate:"one_of=txt,json"`             //json or txt
	HasTimestamp bool   `default:"false"`
	Caller       bool   `default:"false"` //启用堆栈
	Development  bool   `default:"false"` // 记录行号
}
