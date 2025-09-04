package data

import (
	"go-practice/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// 使用 docker-compose 中配置的 MariaDB 连接信息 (端口3307)
	dsn := "app_user:app_password@tcp(127.0.0.1:3307)/album_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_ = db.AutoMigrate(&model.Album{})
	model.DB = db
	return db
}
