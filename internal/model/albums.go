package model

import (
	"gorm.io/gorm"
)

// Album 表示一张唱片
type Album struct {
	gorm.Model
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var DB *gorm.DB
