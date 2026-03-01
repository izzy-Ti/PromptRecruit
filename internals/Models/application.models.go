package models

type Application struct {
	ID     uint `gorm:"primaryKey"`
	JobID  uint
	Job    Jobs `gorm:"foreignKey:JobID"`
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
	Score  float32
}
