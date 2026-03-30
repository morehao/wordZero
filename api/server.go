// Package api 提供WordZero的HTTP服务接口，使用 gin 框架实现
package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server HTTP服务器
type Server struct {
	cfg    *Config
	engine *gin.Engine
	server *http.Server
}

// NewServer 创建新的基于 gin 框架的HTTP服务器
func NewServer(cfg *Config) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("服务器配置不能为空")
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	h, err := newHandler(cfg)
	if err != nil {
		return nil, fmt.Errorf("创建请求处理器失败: %w", err)
	}

	engine := gin.Default()

	// 注册路由
	engine.GET("/health", handleHealth)
	v1 := engine.Group("/api/v1")
	{
		docs := v1.Group("/documents")
		{
			docs.POST("/content", h.handleGenerateFromContent)
			docs.POST("/template", h.handleGenerateFromTemplate)
		}
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeoutSeconds) * time.Second,
	}

	return &Server{
		cfg:    cfg,
		engine: engine,
		server: srv,
	}, nil
}

// Start 启动HTTP服务器
func (s *Server) Start() error {
	log.Printf("[WordZero] HTTP服务启动，监听地址: %s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown 优雅关闭HTTP服务器，等待正在处理的请求完成
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// handleHealth 健康检查接口
// GET /health
func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
