package metricx

type Config struct {
	EndpointUrl string `yaml:"EndpointUrl" validate:"required"` // 链路追踪地址
	Auth        string `yaml:"Auth" validate:"required"`        // 链路追踪认证
	ServerName  string `yaml:"ServerName" validate:"required"`  // 服务名称
	StreamName  string `yaml:"StreamName" default:"default"`
	Interval    int    `yaml:"Interval" default:"60"` //导出时间间隔 单位秒
	IsDebug     bool   `yaml:"IsDebug" default:"false"`
}
