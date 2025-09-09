package model

import "gorm.io/gorm"

type Song struct {
	gorm.Model
	Name      string   `json:"name"`
	Artist    string   `json:"artist"`
	BPM       int      `json:"bpm"`
	Length    int      `json:"length"`
	Tags      []string `json:"tags"`
	ChartList []uint   `json:"chartlist"`
}
