package http

import (
	"github.com/gin-gonic/gin"
	"github.com/tuanliang/restful-api-demo/apps"
	"github.com/tuanliang/restful-api-demo/apps/host"
)

var handler = &Handler{}

// 通过写一个实体类，把内部的接口通过hTTp协议暴露出去
// 需要依赖内部接口的实现
// 该实体类，会实现Gin的Http Handler
type Handler struct {
	svc host.Service
}

func (h *Handler) Config() {

	h.svc = apps.GetImpl(host.AppName).(host.Service)
}
func (h *Handler) Registry(r gin.IRouter) {
	r.POST("/hosts", h.createHost)
	r.GET("/hosts", h.queryHost)
	r.GET("/hosts/:id", h.describeHost)
	r.PUT("/hosts/:id", h.putHost)
	r.PATCH("/hosts/:id", h.patchHost)
}
func (h *Handler) Name() string {
	return host.AppName
}

// 完成Http Handler注册
func init() {
	apps.RegistryGin(handler)
}
