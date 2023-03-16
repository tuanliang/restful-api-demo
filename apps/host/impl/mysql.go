package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/tuanliang/restful-api-demo/apps"
	"github.com/tuanliang/restful-api-demo/apps/host"
	"github.com/tuanliang/restful-api-demo/conf"
)

// 接口实现的静态检查
// var _ host.Service = (*HostServiceImpl)(nil)
var impl = &HostServiceImpl{}

func NewHostServiceImpl() *HostServiceImpl {
	return &HostServiceImpl{
		// Host service 服务的子Logger
		// 封装zap让其满足Logger接口，为什么封装：
		// 1. Logger全局实例 2. Logger Level的动态调整，Logrus不支持Level共同调整 3.加入日志轮转功能的集合
		l:  zap.L().Named("Host"),
		db: conf.C().MySQL.GetDB(),
	}
}

type HostServiceImpl struct {
	l  logger.Logger
	db *sql.DB
}

// 之前都是在start的时候，手动把服务的实现注册到IOC层
// 注册HostService的实例到IOC
// apps.HostService = impl.NewHostServiceImpl()

// 自动执行注册逻辑
func init() {
	// apps.HostService = impl
	apps.Registry(impl)
}
func (i *HostServiceImpl) Config() {
	i.l = zap.L().Named("Host")
	i.db = conf.C().MySQL.GetDB()
}

// 返回服务名称
func (i *HostServiceImpl) Name() string {
	return host.AppName
}
