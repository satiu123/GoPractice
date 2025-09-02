package repository

import (
	"context"
	"go-practice/internal/model"

	"gorm.io/gorm"
)

type AlbumRepository struct {
	db *gorm.DB
}

func NewAlbumRepository(db *gorm.DB) *AlbumRepository {
	return &AlbumRepository{db: db}
}

func (r *AlbumRepository) Create(ctx context.Context, album *model.Album) error {
	return gorm.G[model.Album](r.db).Create(ctx, album)
}

func (r *AlbumRepository) GetAlbumByID(ctx context.Context, id uint) (model.Album, error) {
	return gorm.G[model.Album](r.db).Where("id=?", id).First(ctx)
}

func (r *AlbumRepository) GetAlbums(ctx context.Context) ([]model.Album, error) {
	return gorm.G[model.Album](r.db).Find(ctx)
}
