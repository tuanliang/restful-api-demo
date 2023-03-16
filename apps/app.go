package apps

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tuanliang/restful-api-demo/apps/host"
)

// IOC 容器层：管理所有的服务的实例
// 1.HostService的实例必须注册过来，HostService才会有具体的实例，服务启动时注册
// 2.HTTP 暴露模块，依赖IOC里面的HostService

var (
	HostService host.Service
	implAPPs    = map[string]ImplService{}
	ginAPPs     = map[string]GinService{}
)

func GetImpl(name string) interface{} {
	for k, v := range implAPPs {
		if k == name {
			return v
		}
	}

	return nil
}

func RegistryImpl(svc ImplService) {
	// 服务实例注册到svcs map当中
	if _, ok := implAPPs[svc.Name()]; ok {
		panic(fmt.Sprintf("service %s has registried", svc.Name()))
	}
	implAPPs[svc.Name()] = svc

	if v, ok := svc.(host.Service); ok {
		HostService = v
	}
}
func RegistryGin(svc GinService) {
	// 服务实例注册到svcs map当中
	if _, ok := ginAPPs[svc.Name()]; ok {
		panic(fmt.Sprintf("service %s has registried", svc.Name()))
	}
	ginAPPs[svc.Name()] = svc
}

// 用于初始化，注册到IOC容器里面的所有服务
func InitImpl() {
	for _, v := range implAPPs {
		v.Config()
	}
}

// 已经加载完成的Gin App有哪些
func LoadedGinApps() (names []string) {
	for k := range ginAPPs {
		names = append(names, k)
	}
	return
}

func InitGin(r gin.IRouter) {
	// 先初始化好所有对象
	for _, v := range ginAPPs {
		v.Config()
	}
	// 完成Http Handler的注册
	for _, v := range ginAPPs {
		v.Registry(r)
	}
}

type ImplService interface {
	Config()
	Name() string
}

// 注册由gin编写的handler
type GinService interface {
	Registry(r gin.IRouter)
	Config()
	Name() string
}
