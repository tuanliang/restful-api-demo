package all

// 所有模块的注册
import (
	_ "github.com/tuanliang/restful-api-demo/apps/book/api"
	_ "github.com/tuanliang/restful-api-demo/apps/book/impl"
	_ "github.com/tuanliang/restful-api-demo/apps/host/http"
	_ "github.com/tuanliang/restful-api-demo/apps/host/impl"
)
