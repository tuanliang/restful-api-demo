package impl_test

import (
	"context"
	"testing"

	"github.com/infraboard/mcube/logger/zap"
	"github.com/tuanliang/restful-api-demo/apps/host"
	"github.com/tuanliang/restful-api-demo/apps/host/impl"
)

var (
	// 定义对象必须是满足该接口的实例
	service host.Service
)

func TestCreate(t *testing.T) {
	ins := host.NewHost()
	ins.Name = "test"
	service.CreateHost(context.Background(), ins)
}

func init() {

	// 需要初始化全局Logger
	// 为什么不设置为默认打印，因为性能
	zap.DevelopmentSetup()

	// host service 的具体实现
	service = impl.NewHostServiceImpl()
}
