// Package api 提供WordZero的HTTP服务接口
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Server HTTP服务器
type Server struct {
	cfg    *Config
	mux    *http.ServeMux
	server *http.Server
}

// NewServer 创建新的HTTP服务器
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

	mux := http.NewServeMux()

	// 注册路由
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/v1/documents/content", h.handleGenerateFromContent)
	mux.HandleFunc("/api/v1/documents/template", h.handleGenerateFromTemplate)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeoutSeconds) * time.Second,
	}

	return &Server{
		cfg:    cfg,
		mux:    mux,
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
func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// loggingMiddleware 请求日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		log.Printf("[WordZero] %s %s %d %v", r.Method, r.URL.Path, lrw.statusCode, time.Since(start))
	})
}

// loggingResponseWriter 包装ResponseWriter以记录状态码
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// writeJSON 将数据以JSON格式写入响应
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[WordZero] 写入JSON响应失败: %v", err)
	}
}

// writeError 写入错误响应
func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, ErrorResponse{Error: message})
}
