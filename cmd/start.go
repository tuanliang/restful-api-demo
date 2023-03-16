package cmd

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/spf13/cobra"
	"github.com/tuanliang/restful-api-demo/apps"
	_ "github.com/tuanliang/restful-api-demo/apps/all"
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
			return (err)
		}

		// 初始化全局日志Logger
		if err := loadGlobalLogger(); err != nil {
			return err
		}
		// 加载Host Service实体类
		// service := impl.NewHostServiceImpl()

		// 注册HostService的实例到IOC
		// apps.HostService = impl.NewHostServiceImpl()
		// 采用	_ "github.com/tuanliang/restful-api-demo/apps/host/impl" 完成注册

		// apps.Init()
		apps.InitImpl()

		// 提供一个Gin Router
		g := gin.Default()
		// 注册IOC的所有http handler
		apps.InitGin(g)

		if err := g.Run(conf.C().App.HttpAddr()); err != nil {
			return err
		}

		return errors.New("no flags find")
	},
}

// log 为全局变量, 只需要load 即可全局可用户, 依赖全局配置先初始化
func loadGlobalLogger() error {
	var (
		logInitMsg string
		level      zap.Level
	)
	// 根据Config里面的日志配置，来配置全局Logger对象
	lc := conf.C().Log
	// 设置日志级别
	lv, err := zap.NewLevel(lc.Level)
	if err != nil {
		logInitMsg = fmt.Sprintf("%s, use default level INFO", err)
		level = zap.InfoLevel
	} else {
		level = lv
		logInitMsg = fmt.Sprintf("log level: %s", lv)
	}
	// 使用默认配置初始化Logger的全局配置
	zapConfig := zap.DefaultConfig()
	// 配置日志的Level级别
	zapConfig.Level = level
	// 程序没启动一次，不必都生成一个新日志文件
	zapConfig.Files.RotateOnStartup = false
	// 配置日志的输出方式
	switch lc.To {
	case conf.ToStdout:
		zapConfig.ToStderr = true
		zapConfig.ToFiles = false
	case conf.ToFile:
		zapConfig.Files.Name = "api.log"
		zapConfig.Files.Path = lc.PathDir
	}
	// 配置日志的输出格式
	switch lc.Format {
	case conf.JSONFormat:
		zapConfig.JSON = true
	}
	// 把配置运用到全局Logger
	if err := zap.Configure(zapConfig); err != nil {
		return err
	}
	zap.L().Named("INIT").Info(logInitMsg)
	return nil
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&confFile, "config", "f", "etc/demo.toml", "demo api 配置文件路径")
	RootCmd.AddCommand(StartCmd)
}
