package model

import "time"

// Status 表结构
type Status struct {
	ID       uint      `gorm:"primaryKey"`
	Url      string    `gorm:"default:null"`
	Host     string    `gorm:"default:null"`
	CrawDone int       `gorm:"type:tinyint(1);default:0"`
	CrawTime time.Time `gorm:"default:'2001-01-01 00:00:01'"`
}
