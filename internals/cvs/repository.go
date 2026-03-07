package cvs

import (
	"fmt"
	"math"
	"strings"

	models "github.com/izzy-Ti/PromptRecruit/internals/Models"
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
		Score:  float32(score),
	}
	if err := r.db.Create(&App).Error; err != nil {
		return false, err
	}
	return true, nil
}
func (r *CvRepo) GetUsersCv(jobId uint) (bool, error, [][]float32, string) {
	var Apps []models.Application
	var dbCvs []models.Cvs
	var allUserCvs [][]float32
	var userIDs []uint
	var cv []string

	err := r.db.Where("job_id = ?", jobId).Find(&Apps).Error
	if err != nil {
		return false, err, nil, ""
	}
	if len(Apps) == 0 {
		return false, nil, nil, ""
	}
	for _, app := range Apps {
		userIDs = append(userIDs, app.UserID)
	}
	err = r.db.Where("uploadby IN ?", userIDs).Find(&dbCvs).Error
	if err != nil {
		return false, err, nil, ""
	}

	for _, vec := range dbCvs {
		allUserCvs = append(allUserCvs, vec.Vector.Slice())
		cv = append(cv, vec.Content)
	}
	return true, nil, allUserCvs, strings.Join(cv, "")
}
func (r *CvRepo) GetJobByID(jobID uint) (bool, error, string, [][]float32) {
	var job models.Jobs
	var Jobs []models.JobChunk
	var jobvecs [][]float32
	err := r.db.First(&job, jobID).Error
	if err != nil {
		return false, err, "", nil
	}
	err = r.db.Where("job_id = ?", jobID).Find(&Jobs).Error
	if err != nil {
		return false, err, "", nil
	}
	for _, vecs := range Jobs {
		jobvecs = append(jobvecs, vecs.Vector.Slice())
	}
	fullContent := job.Content
	return true, nil, fullContent, jobvecs
}
func (r *CvRepo) GetMatchScore(jobVector, cvVector []float32) (float32, error) {
	if len(jobVector) == 0 || len(cvVector) == 0 {
		return 0, fmt.Errorf("empty vector")
	}

	if len(jobVector) != len(cvVector) {
		return 0, fmt.Errorf("vector dimensions do not match")
	}

	var dotProduct float32
	var jobNorm float32
	var cvNorm float32

	for i := 0; i < len(jobVector); i++ {
		dotProduct += jobVector[i] * cvVector[i]
		jobNorm += jobVector[i] * jobVector[i]
		cvNorm += cvVector[i] * cvVector[i]
	}
	if jobNorm == 0 || cvNorm == 0 {
		return 0, fmt.Errorf("zero vector found")
	}
	similarity := dotProduct / (float32(math.Sqrt(float64(jobNorm))) * float32(math.Sqrt(float64(cvNorm))))
	score := similarity * 100

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	fmt.Print(score)

	return score, nil
}
func (r *CvRepo) GetBestMatchScore(jobVectors, cvVectors [][]float32) (float32, error) {
	if len(jobVectors) == 0 || len(cvVectors) == 0 {
		return 0, fmt.Errorf("empty vectors")
	}
	fmt.Println("jobVectors:", len(jobVectors))
	fmt.Println("cvVectors:", len(cvVectors))

	var bestScore float32 = 0

	for _, jobVec := range jobVectors {
		for _, cvVec := range cvVectors {
			score, err := r.GetMatchScore(jobVec, cvVec)
			if err != nil {
				continue
			}
			if score > bestScore {
				bestScore = score
			}
		}
	}

	return bestScore, nil
}
func (r *CvRepo) GetUserFullCV(userID uint) (string, error, [][]float32) {
	var cvs []models.Cvs
	var cvec [][]float32

	err := r.db.Where("uploadby = ?", userID).Find(&cvs).Error
	if err != nil {
		return "", err, nil
	}

	var parts []string
	for _, cv := range cvs {
		parts = append(parts, cv.Content)
		cvec = append(cvec, cv.Vector.Slice())
	}

	fullCV := strings.Join(parts, " ")

	return fullCV, nil, cvec
}
func (r *CvRepo) JobAdder(job *models.Jobs) (bool, error) {
	if err := r.db.Create(job).Error; err != nil {
		return false, err
	}
	return true, nil
}
func (r *CvRepo) JobChunkSaver(job models.JobChunk) (bool, error) {
	if err := r.db.Create(&job).Error; err != nil {
		return false, err
	}
	return true, nil
}
