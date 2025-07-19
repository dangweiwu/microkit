package yamlconfig

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type YamlConfig struct {
	cfg interface{} //指针类型
}

func NewYamlConfig(cfg interface{}) (*YamlConfig, error) {
	rfv := reflect.ValueOf(cfg)
	if err := ValidatePtr(&rfv); err != nil {
		return nil, err
	}

	return &YamlConfig{cfg: cfg}, nil

}

func (this *YamlConfig) read(in string) error {
	// 读取原始配置文件内容
	content, err := os.ReadFile(in)
	if err != nil {
		return fmt.Errorf("read config file error: %w", err)
	}

	// 创建模板并解析配置内容
	tpl, err := template.New("config").Parse(string(content))
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	// 收集所有环境变量到map
	envData := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envData[parts[0]] = parts[1]
		}
	}

	// 执行模板替换
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, envData); err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	// 解析处理后的YAML内容
	if err := yaml.Unmarshal(buf.Bytes(), this.cfg); err != nil {
		return fmt.Errorf("yaml unmarshal error: %w", err)
	}
	return nil
}

func (this *YamlConfig) doDefault() error {
	return defaults.Set(this.cfg)
}

func (this *YamlConfig) doValide() error {
	return validator.New().Struct(this.cfg)
}

func Load(configFile string, in interface{}) error {
	cfg, err := NewYamlConfig(in)
	if err != nil {
		return err
	}

	if err := cfg.read(configFile); err != nil {
		return err
	}

	if err := cfg.doDefault(); err != nil {
		return err
	}

	if err := cfg.doValide(); err != nil {
		return err
	}
	return nil
}

func MustLoad(configFile string, in interface{}) {
	if err := Load(configFile, in); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Panicf("配置文件不存在 %v\n", err.Error())
		} else {
			log.Panicf("配置文件错误 %v\n", err.Error())
		}
	}
}

func ValidatePtr(v *reflect.Value) error {
	if !v.IsValid() || v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("not a valid pointer: %v", v)
	}
	return nil
}
