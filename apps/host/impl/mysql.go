package impl

import (
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/tuanliang/restful-api-demo/apps/host"
)

// 接口实现的静态检查
var _ host.Service = (*HostServiceImpl)(nil)

func NewHostServiceImpl() *HostServiceImpl {
	return &HostServiceImpl{
		// Host service 服务的子Logger
		// 封装zap让其满足Logger接口，为什么封装：
		// 1. Logger全局实例 2. Logger Level的动态调整，Logrus不支持Level共同调整 3.加入日志轮转功能的集合
		l: zap.L().Named("Host"),
	}
}

type HostServiceImpl struct {
	l logger.Logger
}
