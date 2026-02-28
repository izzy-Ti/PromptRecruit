package models

import "github.com/pgvector/pgvector-go"

type Cvs struct {
	ID        uint            `gorm:"primaryKey"`
	Content   string          `gorm:"not null"`
	Vector    pgvector.Vector `gorm:"type:vector(1024);not null"`
	SourceURL string
	Uploadby  uint
	User      User `gorm:"foreignKey:Uploadby"`
}
