package protocol

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/tuanliang/restful-api-demo/apps"
	"github.com/tuanliang/restful-api-demo/conf"
)

func NewRestfulService() *RestfulService {
	r := restful.DefaultContainer
	server := &http.Server{
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1M
		Addr:              conf.C().App.RestfulAddr(),
		Handler:           r,
	}
	return &RestfulService{
		server: server,
		l:      zap.L().Named("HTTP Service"),
		r:      r,
	}
}

type RestfulService struct {
	server *http.Server
	l      logger.Logger
	r      *restful.Container
}

func (s *RestfulService) Start() error {
	// 加载handler,把所有模块的Handler注册给了Restful Router
	apps.InitRestful(s.r)

	// 已加载App的日志信息
	apps := apps.LoadedRestApps()
	s.l.Info("loaded rest apps :%v", apps)

	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed { // 正常退出
			s.l.Info("service stopped success")
			return nil
		}
		return fmt.Errorf("start service error,%s", err.Error())
	}

	return nil
}
func (s *RestfulService) Stop() {
	s.l.Info("start graceful shundown")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.l.Warnf("shut down http service error,%s", err)
	}
}
