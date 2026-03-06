package models

import "github.com/pgvector/pgvector-go"

type JobChunk struct {
	ID      uint            `gorm:"primaryKey"`
	JobID   uint            `gorm:"not null"`
	Content string          `gorm:"not null"`
	Vector  pgvector.Vector `gorm:"type:vector(1024);not null"`
}
