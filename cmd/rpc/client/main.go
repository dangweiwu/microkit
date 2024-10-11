package main

import "C"
import (
	"context"
	"log"
	"microkit/connect"
	"microkit/proto/gohello"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

/*
*
连接管理器
tag配置

	name 对应配置文件中name
	fun 选择Dial函数
*/
type Conns struct {
	Hello *grpc.ClientConn `name:"Hello"`
}

type Client struct {
	Hello gohello.HelloClient
}

func main() {
	//日志
	c := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	lg := zap.New(c)
	//配置文件connect管理
	fileConnect, err := connect.NewFileConnect(lg, "./config.yaml")
	if err != nil {
		log.Fatal("initconnect error", err)
	}
	defer fileConnect.Close()

	//连接管理器
	connmanager := connect.NewConnectManager()
	defer connmanager.Close()
	conn := &Conns{}
	if err := connmanager.InitConnect(conn); err != nil {
		log.Fatal(err)
	}

	//客户端管理
	Cli := Client{
		Hello: gohello.NewHelloClient(conn.Hello),
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	log.Println("准备开始")
	for {
		select {
		case <-interrupt:
			return
		case <-time.After(time.Second * 5):
			r, e := Cli.Hello.SayHello(context.Background(), &gohello.Reqmsg{Name: "anan:" + time.Now().Format(time.TimeOnly)})
			log.Println(r, e)
		}
	}

}
