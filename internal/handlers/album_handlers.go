package handlers

import (
	"context"
	"net/http"
	"strconv"

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
	albums, err := h.repo.GetAlbums(context.Background())
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
	err := h.repo.Create(context.Background(), &newAlbum)
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

	album, err := h.repo.GetAlbumByID(context.Background(), uint(albumID))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}
