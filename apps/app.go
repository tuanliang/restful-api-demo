package apps

import "github.com/tuanliang/restful-api-demo/apps/host"

// IOC 容器层：管理所有的服务的实例
// 1.HostService的实例必须注册过来，HostService才会有具体的实例，服务启动时注册
// 2.HTTP 暴露模块，依赖IOC里面的HostService

var (
	HostService host.Service
)
