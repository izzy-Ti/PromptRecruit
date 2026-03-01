package cvs

import (
	"fmt"

	models "github.com/izzy-Ti/PromptRecruit/internals/Models"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type CvRepo struct {
	db *gorm.DB
}

func NewCvRepo(db *gorm.DB) *CvRepo {
	return &CvRepo{db: db}
}

type ScoredCv struct {
	models.Cvs
	Score float64
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
func (r *CvRepo) GetUserCv(jobId uint) (bool, error, [][]float32) {
	var Apps []models.Application
	var dbCvs []models.Cvs
	var allUserCvs [][]float32
	var userIDs []uint

	err := r.db.Where("job_id = ?", jobId).Find(&Apps).Error
	if err != nil {
		return false, err, nil
	}
	if len(Apps) == 0 {
		return false, nil, nil
	}
	for _, app := range Apps {
		userIDs = append(userIDs, app.UserID)
	}
	err = r.db.Where("uploadby IN ?", userIDs).Find(&dbCvs).Error
	if err != nil {
		return false, err, nil
	}

	for _, vec := range dbCvs {
		allUserCvs = append(allUserCvs, vec.Vector.Slice())
	}
	return true, nil, allUserCvs
}
func (r *CvRepo) GetJobByID(jobID uint) (bool, error, [][]float32) {
	var job []models.Jobs
	var jobVec [][]float32

	err := r.db.Where("ID = ?", jobID).Find(&job).Error
	if err != nil {
		return false, err, nil
	}
	for _, vec := range job {
		jobVec = append(jobVec, vec.Vector.Slice())
	}
	return true, nil, jobVec
}
func (r *CvRepo) GetTopMatchingCvs(jobVector []float32, topK int) ([]ScoredCv, error) {
	var results []ScoredCv

	vec := pgvector.NewVector(jobVector)

	err := r.db.Model(&models.Cvs{}).
		Select("cvs.*, (cvs.vector <=> ?) AS score", vec).
		Order("score ASC").
		Limit(topK).
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query top candidates: %v", err)
	}

	return results, nil
}
