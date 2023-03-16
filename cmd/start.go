package cmd

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/tuanliang/restful-api-demo/apps"
	_ "github.com/tuanliang/restful-api-demo/apps/all"
	"github.com/tuanliang/restful-api-demo/apps/host/http"
	"github.com/tuanliang/restful-api-demo/conf"
)

var (
	confType string
	confFile string
	confETCD string
)

// 程序的启动时，组装都在这里进行
// 1.
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "启动 demo 后端API",
	Long:  "启动 demo 后端API",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载程序配置
		err := conf.LoadConfigFromToml(confFile)
		if err != nil {
			panic(err)
		}

		// 加载Host Service实体类
		// service := impl.NewHostServiceImpl()

		// 注册HostService的实例到IOC
		// apps.HostService = impl.NewHostServiceImpl()
		// 采用	_ "github.com/tuanliang/restful-api-demo/apps/host/impl" 完成注册

		apps.Init()

		// 通过Host api Handler提供 HTTP RestFul接口
		api := http.NewHostHTTPHandler()
		api.Config()

		// 提供一个Gin Router
		g := gin.Default()
		api.Registry(g)

		if err := g.Run(conf.C().App.HttpAddr()); err != nil {
			return err
		}

		return errors.New("no flags find")
	},
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&confFile, "config", "f", "etc/demo.toml", "demo api 配置文件路径")
	RootCmd.AddCommand(StartCmd)
}
