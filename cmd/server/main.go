package main

import (
	"log"

	"go-practice/data"
	"go-practice/internal/cache"
	"go-practice/internal/config"
	"go-practice/internal/handlers"
	"go-practice/internal/repository"
	"go-practice/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化Redis
	redisConfig := cache.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}
	// 初始化DB
	dbConfig := data.DBConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
	}
	if err := cache.InitRedis(redisConfig); err != nil {
		log.Printf("Redis initialization failed: %v", err)
		log.Println("Application will continue without Redis cache")
	}

	// 初始化数据库
	db := data.InitDB(dbConfig)
	// 初始化repo
	albumRepo := repository.NewAlbumRepository(db)

	// 初始化handler
	albumHandler := handlers.NewAlbumHandler(albumRepo)

	// 初始化 Gin 引擎（默认包含 Logger 与 Recovery 中间件）
	// r := gin.Default()
	r := gin.New()
	// 显式设置可信代理为 nil，避免警告；生产环境按需配置
	_ = r.SetTrustedProxies(nil)
	// 注册路由
	router.Register(r, albumHandler)

	// 启动 HTTP 服务（默认 :8080）
	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
