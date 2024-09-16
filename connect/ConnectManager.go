package connect

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"reflect"
)

/*
*
多connect管理器
*/
type DialConnect func(target string) (conn *grpc.ClientConn, err error)

type ConnectManager struct {
	funcMap map[string]DialConnect
	conn    map[string]*grpc.ClientConn
}

func NewConnectManager() *ConnectManager {
	return &ConnectManager{
		funcMap: make(map[string]DialConnect),
		conn:    make(map[string]*grpc.ClientConn),
	}
}
func (this *ConnectManager) RegFunc(fname string, conn DialConnect) {
	this.funcMap[fname] = conn
}

func (this *ConnectManager) InitConnect(t interface{}) error {
	// 获取 t 的反射类型
	tVal := reflect.ValueOf(t)
	tType := tVal.Type()

	// 确保 t 是一个指针
	if tType.Kind() != reflect.Ptr {
		return errors.New("must be a pointer")

	}
	// 获取基础类型
	elemType := tType.Elem()
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		fieldVal := tVal.Elem().Field(i)
		var conn *grpc.ClientConn
		var err error
		// 获取字段上的 Host 标签
		name := field.Tag.Get("name")
		if name == "" {
			continue
		}

		fun := field.Tag.Get("fun")
		if len(fun) == 0 {
			conn, err = this.defaultFunc(name)
		} else {
			if f, ok := this.funcMap[fun]; ok {
				conn, err = f(name)
			} else {
				conn, err = this.defaultFunc(name)
			}
		}

		if err != nil {
			return err
		}
		if fieldVal.Kind() == reflect.Ptr {
			fieldVal.Set(reflect.ValueOf(conn))
			this.conn[name] = conn
		} else {
			return fmt.Errorf("unsupported type for field %s", field.Name)
		}
	}
	return nil
}
func (this *ConnectManager) defaultFunc(name string) (conn *grpc.ClientConn, err error) {
	return grpc.NewClient(name, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
}

func (this *ConnectManager) Close() {
	for _, v := range this.conn {
		v.Close()
	}
}
