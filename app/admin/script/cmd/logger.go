package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"go-web-template/app/admin/internal/global"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"strings"
)

var level string

func init() {
	rootCmd.AddCommand(loggerCmd)
	flags := loggerCmd.PersistentFlags()
	flags.StringVar(&level, "level", "", "日志级别 [debug,info,warn,error]")
	loggerCmd.MarkPersistentFlagRequired("level")
}

var loggerCmd = &cobra.Command{
	Use:   "logger",
	Short: "打印日志",
	Long:  "打印一条日志",
	PreRun: func(cmd *cobra.Command, args []string) {
		config.Init("./configs/config_dev.toml", global.Configs)
		logger.Init(global.Configs.App.Name, global.Configs.LoggerConfig)
	},
	Run: func(cmd *cobra.Command, args []string) {
		switch strings.ToLower(level) {
		case "debug":
			logger.Debug(context.Background(), "debug")
		case "info":
			logger.Info(context.Background(), "debug")
		case "warn":
			logger.Warn(context.Background(), "debug")
		case "error":
			logger.Error(context.Background(), "debug")
		default:
			logger.Debug(context.Background(), "debug")
		}
	},
}
