package main

import (
	"log"

	"go-practice/data"
	"go-practice/internal/handlers"
	"go-practice/internal/repository"
	"go-practice/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {

	// 初始化数据库
	db := data.InitDB()
	// 初始化repo
	albumRepo := repository.NewAlbumRepository(db)

	// 初始化handler
	albumHandler := handlers.NewAlbumHandler(albumRepo)

	// 初始化 Gin 引擎（默认包含 Logger 与 Recovery 中间件）
	r := gin.Default()
	// 显式设置可信代理为 nil，避免警告；生产环境按需配置
	_ = r.SetTrustedProxies(nil)
	// 注册路由
	router.Register(r, albumHandler)

	// 启动 HTTP 服务（默认 :8080）
	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
