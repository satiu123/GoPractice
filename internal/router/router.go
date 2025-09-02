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
	r.GET("/albums", ah.GetAlbums)
	r.GET("/albums/:id", ah.GetAlbumByID)
	r.POST("/albums", ah.PostAlbums)
}
