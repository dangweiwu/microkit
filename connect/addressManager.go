package connect

import (
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
	"log"
)

/*
地址管理器
全量更新地址
实现grpc builder接口

*/

// 获取地址列表接口
type IAddressManager interface {
	GetAddr(servicename string) []string
	IsReload() <-chan struct{} //地址是否变更
}

// 跟踪上游地址 通知下游更新
type AddressManager struct {
	iaddrs IAddressManager
	conn   map[string]resolver.ClientConn
	log    *zap.Logger
	close  chan struct{}
}

func NewAddressManager(log *zap.Logger, iaddrs IAddressManager) *AddressManager {
	if iaddrs == nil {
		panic("iaddrs is nil")
	}
	a := &AddressManager{
		iaddrs: iaddrs,
		log:    log,
		close:  make(chan struct{}),
		conn:   map[string]resolver.ClientConn{},
	}
	go a.watch()
	return a
}

// dail时候进行调用一次
func (this *AddressManager) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ips := this.iaddrs.GetAddr(target.Endpoint())
	this.conn[target.Endpoint()] = cc
	_address := []resolver.Address{}
	for _, v := range ips {
		_address = append(_address, resolver.Address{Addr: v})
	}
	cc.UpdateState(resolver.State{Addresses: _address})
	this.log.Info("ADDRES注册成功", zap.String("service", target.Endpoint()), zap.Any("address", _address))
	return this, nil
}
func (this *AddressManager) Scheme() string {
	return ""
}

// 多次调用
func (this *AddressManager) ResolveNow(rn resolver.ResolveNowOptions) {

}

// 遍历当前注册的client端
func (this *AddressManager) updateAll() {
	for k, v := range this.conn {
		ips := this.iaddrs.GetAddr(k)
		_address := []resolver.Address{}
		if len(ips) != 0 {
			for _, v := range ips {
				_address = append(_address, resolver.Address{Addr: v})
			}
		}
		v.UpdateState(resolver.State{Addresses: _address})
	}
}

// 监听地址变化
func (this *AddressManager) watch() {
	for {
		select {
		case <-this.iaddrs.IsReload():
			this.updateAll()
		case <-this.close:
			break
		}
	}
}

func (this *AddressManager) Close() {
	log.Println("close address Manager")
	this.close <- struct{}{}
}
