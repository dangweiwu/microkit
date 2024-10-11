package mysqlx

import (
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
)

func FakeDb(dbName string) (rhost string, se *server.Server, err error) {

	db := memory.NewDatabase(dbName)
	pro := memory.NewDBProvider(db)
	engine := sqle.NewDefault(pro)

	config := server.Config{
		Protocol: "tcp",
		Address:  ":0", //随机端口
	}
	s, err := server.NewServer(config, engine, memory.NewSessionBuilder(pro), nil)
	if err != nil {
		return "", nil, err
	}
	go s.Start()

	rhost = s.Listener.Addr().String()
	return rhost, s, nil
}
