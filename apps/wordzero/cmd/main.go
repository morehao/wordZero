package main

import (
	"net"
	"time"

	wordzero "github.com/zerx-lab/wordZero/apps/wordzero"
	"github.com/zerx-lab/wordZero/apps/wordzero/internal/config"
	"github.com/zerx-lab/wordZero/pkg/global"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/ygpkg/yg-go/apis/runtime/server"
	"github.com/ygpkg/yg-go/lifecycle"
	"github.com/ygpkg/yg-go/logs"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use: "wordzero",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadConfig(configFile)
			if err != nil {
				logs.Warnf("[main] load config failed: %s", err)
			}
			if cfg == nil {
				logs.Warnf("[main] use default config")
			}
			logs.ReloadConfig(cfg.MainConf.App, cfg.LogsConf)
			logs.Debugf("[main] config loaded, env: %s", cfg.MainConf.Env)
		},
	}
)

func main() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file path")
	rootCmd.Run = mainRun()
	rootCmd.Execute()
}

func mainRun() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		cfg := config.Conf()
		defer time.Sleep(time.Second)

		if err := initS3Uploader(cfg); err != nil {
			logs.Errorf("[main] init S3 failed: %s", err)
			return
		}

		svr := server.NewRouter(global.PrefixAPIV1)

		if cfg.MainConf.Env == "dev" || cfg.MainConf.Env == "test" {
			gin.SetMode(gin.DebugMode)
		}

		l, err := net.Listen("tcp", cfg.MainConf.HttpAddr)
		if err != nil {
			logs.Fatalf("[main] listen at %s failed: %s", cfg.MainConf.HttpAddr, err)
			return
		}

		wordzero.Routers(svr)
		logs.Infof("[main] start http server at %s", cfg.MainConf.HttpAddr)

		if err := svr.Run(l); err != nil {
			logs.Errorf("[main] run server failed: %s", err)
			return
		}
		lifecycle.Std().WaitExit()
	}
}
