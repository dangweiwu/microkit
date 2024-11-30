package metricx

type Config struct {
	EndpointUrl string `yaml:"endpointUrl" validate:"required"` // 链路追踪地址
	Auth        string `yaml:"auth" validate:"required"`        // 链路追踪认证
	ServerName  string `yaml:"serverName" validate:"required"`  // 服务名称
	StreamName  string `yaml:"streamName" default:"default"`
	Interval    int    `yaml:"interval" default:"60"` //导出时间间隔 单位秒
	IsDebug     bool   `yaml:"isDebug" default:"false"`
}
