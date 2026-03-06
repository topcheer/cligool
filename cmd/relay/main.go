package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cligool/cligool/internal/relay"
	"github.com/cligool/cligool/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// 初始化数据库
	db, err := database.NewPostgresDB(
		os.Getenv("DATABASE_URL"),
		os.Getenv("REDIS_URL"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// 运行数据库迁移
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 创建中继服务
	relayService := relay.NewService(relay.Config{
		DB:   db,
		Host: getEnv("RELAY_HOST", "0.0.0.0"),
		Port: getEnv("RELAY_PORT", "8080"),
	})

	// 创建HTTP服务器
	router := gin.Default()

	// 静态文件服务（Web界面）
	router.LoadHTMLGlob("./web/*.html")

	// 静态JavaScript和CSS库
	router.Static("/lib", "./web/lib")

	// 下载文件目录
	router.Static("/downloads", "./web/downloads")

	// 首页 - 着陆页面（包含下载链接）
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "landing.html", gin.H{
			"title": "CliGool - 跨平台远程终端",
		})
	})

	// 会话页面 - 直接通过URL访问
	router.GET("/session/:session_id", func(c *gin.Context) {
		sessionID := c.Param("session_id")
		c.HTML(http.StatusOK, "terminal.html", gin.H{
			"title":      "CliGool - 远程终端",
			"session_id": sessionID,
		})
	})

	// API路由
	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"time":   time.Now().Unix(),
			})
		})

		// WebSocket终端连接
		api.GET("/terminal/:session_id", relayService.HandleTerminalConnection)

		// 会话管理
		api.POST("/sessions", relayService.CreateSession)
		api.GET("/sessions/:id", relayService.GetSession)
		api.DELETE("/sessions/:id", relayService.DeleteSession)
		api.GET("/sessions", relayService.ListSessions)
	}

	// 创建服务器
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", relayService.Config.Host, relayService.Config.Port),
		Handler: router,
	}

	// 启动服务器
	go func() {
		log.Printf("Starting relay server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}