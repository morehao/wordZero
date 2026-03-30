// Package main WordZero HTTP服务启动入口
package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/zerx-lab/wordZero/api"
)

func main() {
	// 命令行参数
	configFile := flag.String("config", "", "配置文件路径（JSON格式）")
	flag.Parse()

	// 加载配置
	cfg := api.DefaultConfig()
	if *configFile != "" {
		if err := loadConfigFromFile(*configFile, cfg); err != nil {
			log.Fatalf("[WordZero] 加载配置文件失败: %v", err)
		}
		log.Printf("[WordZero] 已加载配置文件: %s", *configFile)
	} else {
		// 从环境变量加载配置
		loadConfigFromEnv(cfg)
		log.Printf("[WordZero] 从环境变量加载配置")
	}

	// 创建并启动服务器
	srv, err := api.NewServer(cfg)
	if err != nil {
		log.Fatalf("[WordZero] 创建服务器失败: %v", err)
	}

	// 优雅关闭处理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// 在后台启动服务器
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("[WordZero] 服务器已停止: %v", err)
		}
	}()

	// 等待退出信号
	<-quit
	log.Println("[WordZero] 正在关闭服务器...")

	// 给正在处理的请求最多30秒时间完成
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[WordZero] 关闭服务器时出错: %v", err)
	}
	log.Println("[WordZero] 服务器已关闭")
}

// loadConfigFromFile 从JSON文件加载配置
func loadConfigFromFile(path string, cfg *api.Config) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(cfg)
}

// loadConfigFromEnv 从环境变量加载配置
func loadConfigFromEnv(cfg *api.Config) {
	if v := os.Getenv("SERVER_HOST"); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Port = port
		}
	}
	if v := os.Getenv("S3_ENDPOINT"); v != "" {
		cfg.S3Config.Endpoint = v
	}
	if v := os.Getenv("S3_REGION"); v != "" {
		cfg.S3Config.Region = v
	}
	if v := os.Getenv("S3_BUCKET"); v != "" {
		cfg.S3Config.Bucket = v
	}
	if v := os.Getenv("S3_ACCESS_KEY_ID"); v != "" {
		cfg.S3Config.AccessKeyID = v
	}
	if v := os.Getenv("S3_SECRET_ACCESS_KEY"); v != "" {
		cfg.S3Config.SecretAccessKey = v
	}
	if v := os.Getenv("S3_KEY_PREFIX"); v != "" {
		cfg.S3Config.KeyPrefix = v
	}
	if v := os.Getenv("S3_USE_PATH_STYLE"); v != "" {
		cfg.S3Config.UsePathStyle = v == "true" || v == "1"
	}
	if v := os.Getenv("S3_PUBLIC_BASE_URL"); v != "" {
		cfg.S3Config.PublicBaseURL = v
	}
}
