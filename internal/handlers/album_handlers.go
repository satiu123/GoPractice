package handlers

import (
	"net/http"
	"strconv"

	"go-practice/internal/cache"
	"go-practice/internal/model"
	"go-practice/internal/repository"

	"github.com/gin-gonic/gin"
)

type AlbumHandler struct {
	repo *repository.AlbumRepository
}

func NewAlbumHandler(repo *repository.AlbumRepository) *AlbumHandler {
	return &AlbumHandler{repo}
}

// GetAlbums 响应一个包含所有专辑列表的JSON

func (h *AlbumHandler) GetAlbums(c *gin.Context) {
	albums, err := h.repo.GetAlbums(c.Request.Context())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(albums) == 0 {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Albums is empty!"})
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}

// PostAlbums 从请求的JSON body中添加一张新专辑
func (h *AlbumHandler) PostAlbums(c *gin.Context) {
	var newAlbum model.Album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	err := h.repo.Create(c.Request.Context(), &newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, newAlbum)
}

// GetAlbumByID 根据ID定位并返回一张专辑
func (h *AlbumHandler) GetAlbumByID(c *gin.Context) {
	id := c.Param("id")
	albumID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid album ID"})
		return
	}

	album, err := h.repo.GetAlbumByID(c.Request.Context(), uint(albumID))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

// ClearCache 清除所有缓存
func (h *AlbumHandler) ClearCache(c *gin.Context) {
	if cache.RedisClient == nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"error": "Redis is not available"})
		return
	}

	// 清除所有专辑相关缓存
	keys := []string{"albums:all"}

	// 获取所有专辑ID的缓存键
	if result, err := cache.RedisClient.Keys(c.Request.Context(), "album:*").Result(); err == nil {
		keys = append(keys, result...)
	}

	// 删除所有键
	for _, key := range keys {
		if err := cache.DeleteCache(key); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cache"})
			return
		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Cache cleared successfully"})
}

// GetCacheStatus 获取缓存状态
func (h *AlbumHandler) GetCacheStatus(c *gin.Context) {
	if cache.RedisClient == nil {
		c.IndentedJSON(http.StatusOK, gin.H{"redis_status": "disconnected"})
		return
	}

	// 检查Redis连接
	_, err := cache.RedisClient.Ping(c.Request.Context()).Result()
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"redis_status": "error",
			"error":        err.Error(),
		})
		return
	}

	// 获取缓存统计信息
	albumsExists, _ := cache.Exists("albums:all")
	albumKeys, _ := cache.RedisClient.Keys(c.Request.Context(), "album:*").Result()

	c.IndentedJSON(http.StatusOK, gin.H{
		"redis_status":             "connected",
		"albums_cached":            albumsExists,
		"individual_albums_cached": len(albumKeys),
		"cached_album_keys":        albumKeys,
	})
}
