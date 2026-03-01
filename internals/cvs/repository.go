package cvs

import (
	models "github.com/izzy-Ti/PromptRecruit/internals/Models"
	"gorm.io/gorm"
)

type CvRepo struct {
	db *gorm.DB
}

func NewCvRepo(db *gorm.DB) *CvRepo {
	return &CvRepo{db: db}
}
func (r *CvRepo) ApplicationSaver(jobId, userId uint, score float32) (bool, error) {
	App := models.Application{
		JobID:  jobId,
		UserID: userId,
		Score:  score,
	}
	if err := r.db.Create(&App).Error; err != nil {
		return false, err
	}
	return true, nil
}
func (r *CvRepo) GetUserCv(userId uint) (bool, error, [][]float32) {
	var UserCv []models.Cvs
	var Cvec [][]float32
	err := r.db.Where("Uploadby = ?", userId).Find(&UserCv).Error
	if err != nil {
		return false, err, nil
	}
	for _, vec := range UserCv {
		Cvec = append(Cvec, vec.Vector.Slice())
	}
	return true, nil, Cvec
}
