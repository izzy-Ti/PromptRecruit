package models

type User struct {
	ID               uint   `gorm:"primaryKey"`
	Name             string `gorm:"not null"`
	Email            string `gorm:"uniqueIndex;not null"`
	Password         string `gorm:""`
	VerifyOTP        string `gorm:"default:''"`
	OTPExpireAt      int64  `gorm:"default:0"`
	IsAccVerified    bool   `gorm:"default:false"`
	ResetOTP         string `gorm:"default:''"`
	ResetOTPExpireAt int64  `gorm:"default:0"`
	Avater           string `gorm:"default:''"`
	Role             string `gorm:"default: 'USER'"`
	AuthType         string `gorm:"size:20;not null"`
	GoogleId         string `gorm:""`
}
