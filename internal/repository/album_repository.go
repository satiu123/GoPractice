package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"go-practice/internal/cache"
	"go-practice/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AlbumRepository struct {
	db *gorm.DB
}

func NewAlbumRepository(db *gorm.DB) *AlbumRepository {
	return &AlbumRepository{db: db}
}

func (r *AlbumRepository) Create(ctx context.Context, album *model.Album) error {
	err := gorm.G[model.Album](r.db).Create(ctx, album)
	if err != nil {
		return err
	}

	// 清除相关缓存
	r.clearAlbumsCache()

	// 缓存新创建的专辑
	albumKey := fmt.Sprintf("album:%d", album.ID)
	if cacheErr := cache.SetCacheWithContext(ctx, albumKey, album, time.Hour); cacheErr != nil {
		log.Printf("Failed to cache album: %v", cacheErr)
	}

	return nil
}

func (r *AlbumRepository) GetAlbumByID(ctx context.Context, id uint) (model.Album, error) {
	var album model.Album
	albumKey := fmt.Sprintf("album:%d", id)

	// 尝试从缓存获取
	if cache.RedisClient != nil {
		if err := cache.GetCacheJSON(albumKey, &album); err == nil {
			log.Printf("Album %d retrieved from cache", id)
			return album, nil
		} else if err != redis.Nil {
			log.Printf("Cache error: %v", err)
		}
	}

	// 从数据库获取
	album, err := gorm.G[model.Album](r.db).Where("id=?", id).First(ctx)
	if err != nil {
		return album, err
	}

	// 缓存结果
	if cache.RedisClient != nil {
		if cacheErr := cache.SetCacheWithContext(ctx, albumKey, album, time.Hour); cacheErr != nil {
			log.Printf("Failed to cache album: %v", cacheErr)
		}
	}

	return album, nil
}

func (r *AlbumRepository) GetAlbums(ctx context.Context) ([]model.Album, error) {
	var albums []model.Album
	albumsKey := "albums:all"

	// 尝试从缓存获取
	if cache.RedisClient != nil {
		if err := cache.GetCacheJSON(albumsKey, &albums); err == nil {
			log.Println("Albums retrieved from cache")
			return albums, nil
		} else if err != redis.Nil {
			log.Printf("Cache error: %v", err)
		}
	}

	// 从数据库获取
	albums, err := gorm.G[model.Album](r.db).Find(ctx)
	if err != nil {
		return albums, err
	}

	// 缓存结果
	if cache.RedisClient != nil {
		if cacheErr := cache.SetCacheWithContext(ctx, albumsKey, albums, 30*time.Minute); cacheErr != nil {
			log.Printf("Failed to cache albums: %v", cacheErr)
		}
	}

	return albums, nil
}

// clearAlbumsCache 清除专辑列表缓存
func (r *AlbumRepository) clearAlbumsCache() {
	if cache.RedisClient != nil {
		if err := cache.DeleteCache("albums:all"); err != nil {
			log.Printf("Failed to clear albums cache: %v", err)
		}
	}
}
