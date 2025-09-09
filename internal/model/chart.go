package model

import "gorm.io/gorm"

type StatusCode uint8

const (
	Stable StatusCode = iota
	Alpha
	Beta
)

type Mode uint8

const (
	Key Mode = iota
	Taiko
	Ring
)

func (s StatusCode) String() string {
	switch s {
	case Stable:
		return "Stable"
	case Alpha:
		return "Alpha"
	case Beta:
		return "Beta"
	default:
		return "Unknown"
	}
}
func (m Mode) String() string {
	switch m {
	case Key:
		return "Key"
	case Taiko:
		return "Taiko"
	case Ring:
		return "Ring"
	default:
		return "Unknown"
	}
}

type Chart struct {
	gorm.Model
	Version   string     `json:"version"`
	Length    int        `json:"length"`
	ChartMode Mode       `json:"chartmode"`
	Tags      []string   `json:"tags"`
	Status    StatusCode `json:"status"`
	Creator   string     `json:"creator"`
	SongID    uint       `json:"songid"`
}
