package cvs

import (
	"gorm.io/gorm"
)

type CvRepo struct {
	db *gorm.DB
}

func NewCvRepo(db *gorm.DB) *CvRepo {
	return &CvRepo{db: db}
}
