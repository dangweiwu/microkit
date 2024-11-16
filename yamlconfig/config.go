package yamlconfig

import (
	"errors"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"reflect"
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
	bytes, err := os.ReadFile(in)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(bytes, this.cfg); err != nil {
		return err
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
