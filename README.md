<h1 align="center">Go Demo X</h1>
<p align="center">
<img alt="Static Badge" src="https://img.shields.io/badge/Go- 1.9-blue">

<img alt="Static Badge" src="https://img.shields.io/badge/license- MIT-blue">

</p>

_GRPC基于配置文件的地址注册与发现_

---
## 特点：
- 1 服务注册需要手动进行yaml文件的配置。
- 2 服务发现会根据文件的更新实时进行。
---
### 使用说明:
```
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
			break
		case <-time.After(time.Second * 5):
			r, e := Cli.Hello.SayHello(context.Background(), &gohello.Reqmsg{Name: "anan:" + time.Now().Format(time.TimeOnly)})
			log.Println(r, e)
		}
	}

```

---
### 其他
1. 因为是基于文件的服务注册与发现，具有简单可靠低资源性，适合小规模服务注册与发现尤其适合docker-compose项目。

---
### License
© Dangweiwu, 2024~time.Now

Released under the MIT License