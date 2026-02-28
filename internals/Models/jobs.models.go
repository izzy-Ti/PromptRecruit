package models

import "github.com/pgvector/pgvector-go"

type Jobs struct {
	ID       uint `gorm:"primaryKey"`
	Title    string
	Content  string          `gorm:"not null"`
	Vector   pgvector.Vector `gorm:"type:vector(1024); not null"`
	Uploadby uint
	User     User `gorm:"foreignKey:Uploadby"`
}
