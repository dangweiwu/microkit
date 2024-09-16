package connect

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
	"time"
)

/**
 * 连接配置文件
监听配置文件配置文件变动
*/

type ConfigConnect struct {
	sync.RWMutex
	*AddressManager
	filePath string
	log      *zap.Logger
	addrsMap map[string][]string
	close    chan struct{}
	flash    chan struct{}
}

const debounceDelay = time.Second

func NewFileConnect(log *zap.Logger, configPath string) (*ConfigConnect, error) {

	a := &ConfigConnect{
		filePath: configPath,
		log:      log,
		close:    make(chan struct{}),
		flash:    make(chan struct{}),
	}
	a.AddressManager = NewAddressManager(log, a)
	if err := a.loadConfig(); err != nil {

		return nil, err
	}
	go a.watchFile()
	resolver.Register(a)
	return a, nil
}

func (this *ConfigConnect) IsReload() <-chan struct{} {
	return this.flash
}
func (this *ConfigConnect) loadConfig() error {
	data, err := os.ReadFile(this.filePath)
	if err != nil {
		return err
	}
	this.Lock()
	defer this.Unlock()
	if err := yaml.Unmarshal(data, &this.addrsMap); err != nil {
		return err
	}
	this.log.Info("加载配置文件完成", zap.Any("data", this.addrsMap))
	return nil
}
func (this *ConfigConnect) watchFile() {
	var lastEventTime time.Time
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	if err := w.Add(this.filePath); err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				currentTime := time.Now()
				if currentTime.Sub(lastEventTime) > debounceDelay {
					lastEventTime = currentTime
					this.log.Info("配置文件变动",
						zap.Time("time", currentTime),
						zap.Time("last", lastEventTime), zap.Any("dur", currentTime.Sub(lastEventTime)))
					this.loadConfig()
					this.flash <- struct{}{}
				}

			}
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		case <-this.close:
			this.AddressManager.Close()
			log.Println("watch config close")
			break
		}
	}

}

func (this *ConfigConnect) Close() {
	fmt.Println("")
	this.close <- struct{}{}
	this.AddressManager.Close()
}

func (this *ConfigConnect) GetAddr(servicename string) []string {
	this.RLock()
	defer this.RUnlock()
	return this.addrsMap[servicename]
}
