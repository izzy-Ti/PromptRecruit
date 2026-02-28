package models

type Application struct {
	ID     uint `gorm:"primaryKey"`
	JobID  uint
	UserID uint
	Score  float32
}
