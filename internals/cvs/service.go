package cvs

import (
	models "github.com/izzy-Ti/PromptRecruit/internals/Models"
	rag "github.com/izzy-Ti/PromptRecruit/internals/Rag"
	"github.com/pgvector/pgvector-go"
)

type CVservice struct {
	repo *CvRepo
}

func NewUserService(repo *CvRepo) *CVservice {
	return &CVservice{repo: repo}
}

func CvUploadSvc() {

}

func (s *CVservice) ApplicationService(userId, JobId uint) (float32, bool, error) {
	ok, err, content, jobVecs := s.repo.GetJobByID(JobId)
	if !ok {
		return 0, false, err
	}
	Cv, err, cvec := s.repo.GetUserFullCV(userId)
	if err != nil {
		return 0, false, err
	}
	LLMscore, err := rag.UserScore(content, Cv)
	score, err := s.repo.GetBestMatchScore(jobVecs, cvec)
	if err != nil {
		return 0, false, err
	}
	finalscore := score*0.7 + float32(LLMscore)*0.3
	s.repo.ApplicationSaver(JobId, userId, float32(finalscore))
	return finalscore, true, nil
}
func (s *CVservice) jobAddService(Title, content string, userId uint) (bool, error) {
	//var jobChunk []models.JobChunk

	chunks := rag.ChunkText(content, 500)
	vecs, err := rag.EmbedText(content)

	if err != nil {
		return false, err
	}

	job := &models.Jobs{
		Title:    Title,
		Content:  content,
		Uploadby: userId,
	}
	if ok, err := s.repo.JobAdder(job); !ok || err != nil {
		return false, err
	}

	for i, vec := range vecs {
		if len(vec) == 0 {
			continue
		}
		jobChunk := models.JobChunk{
			JobID:   job.ID,
			Vector:  pgvector.NewVector(vec),
			Content: chunks[i],
		}
		if ok, err := s.repo.JobChunkSaver(jobChunk); !ok || err != nil {
			return false, err
		}
	}
	// for _, chunk := range jobChunk {
	// 	s.repo.JobChunkSaver(chunk)
	// }

	return true, nil
}
