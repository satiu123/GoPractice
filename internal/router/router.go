package router

import (
	"go-practice/internal/handlers"

	"github.com/gin-gonic/gin"
)

// Register 注册所有路由
func Register(r *gin.Engine, ah *handlers.AlbumHandler) {
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(
			200,
			"pong",
		)
	})

	// 专辑相关路由
	r.GET("/albums", ah.GetAlbums)
	r.GET("/albums/:id", ah.GetAlbumByID)
	r.POST("/albums", ah.PostAlbums)

	// 缓存管理路由
	r.DELETE("/cache", ah.ClearCache)
	r.GET("/cache/status", ah.GetCacheStatus)
}
