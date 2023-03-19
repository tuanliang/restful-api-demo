package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/tuanliang/restful-api-demo/apps"
	"github.com/tuanliang/restful-api-demo/apps/book"
	"github.com/tuanliang/restful-api-demo/conf"
	"google.golang.org/grpc"
)

var (
	// Service 服务实例
	svr = &service{}
)

// 这个就是Grpc接口的实现类
type service struct {
	db  *sql.DB
	log logger.Logger
	book.UnimplementedServiceServer
}

func (s *service) Config() {

	s.log = zap.L().Named(s.Name())
	s.db = conf.C().MySQL.GetDB()
	return

}

func (s *service) Name() string {
	return book.AppName
}
func (s *service) Registry(server *grpc.Server) {
	book.RegisterServiceServer(server, svr)
}

func init() {
	apps.RegistryGrpc(svr)
}
