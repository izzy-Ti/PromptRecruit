package models


type Jobs struct {
	ID       uint `gorm:"primaryKey"`
	Title    string
	Content  string          `gorm:"not null"`
	Uploadby uint
	User     User `gorm:"foreignKey:Uploadby"`
}