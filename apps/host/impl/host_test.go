package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/infraboard/mcube/logger/zap"
	"github.com/stretchr/testify/assert"
	"github.com/tuanliang/restful-api-demo/apps/host"
	"github.com/tuanliang/restful-api-demo/apps/host/impl"
	"github.com/tuanliang/restful-api-demo/conf"
)

var (
	// 定义对象必须是满足该接口的实例
	service host.Service
)

func TestCreate(t *testing.T) {
	should := assert.New(t)
	ins := host.NewHost()
	ins.Id = "ins-02"
	ins.Name = "test"
	ins.Region = "hanghzou"
	ins.Type = "sm1"
	ins.CPU = 1
	ins.Memory = 2048

	ins, err := service.CreateHost(context.Background(), ins)
	if should.NoError(err) {
		fmt.Println(ins)
	}

}
func TestQuery(t *testing.T) {
	should := assert.New(t)
	req := host.NewQueryHostRequest()
	req.Keywords = "接口测试"
	set, err := service.QueryHost(context.Background(), req)
	if should.NoError(err) {
		for i := range set.Items {
			fmt.Println(set.Items[i].Id)
		}
	}
}
func TestDescribe(t *testing.T) {
	should := assert.New(t)
	req := host.NewDescribeHostRequestWithId("ins-05")

	ins, err := service.DescribeHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(ins.Id)
	}
}
func TestUpdate(t *testing.T) {
	should := assert.New(t)
	// req := host.NewPatchUpdateHostRequest("ins-05")
	req := host.NewPutUpdateHostRequest("ins-05")
	req.Name = "修改测试01"
	req.Region = "rg 02"
	req.Type = "small"
	req.CPU = 1
	req.Memory = 2048
	req.Description = "测试更新"
	ins, err := service.UpdateHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(ins.Id)
	}
}
func init() {
	// 测试用例的配置文件
	err := conf.LoadConfigFromToml("../../../etc/demo.toml")
	if err != nil {
		fmt.Println("配置文件错误")
		//panic(err)
	}

	// 需要初始化全局Logger
	// 为什么不设置为默认打印，因为性能
	zap.DevelopmentSetup()

	// host service 的具体实现
	service = impl.NewHostServiceImpl()

}
