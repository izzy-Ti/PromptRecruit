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

func (s *CVservice) ApplicationService(userId, JobId uint) (bool, error) {
	ok, err, content := s.repo.GetJobByID(JobId)
	if !ok {
		return false, err
	}
	Cv, err := s.repo.GetUserFullCV(userId)
	if err != nil {
		return false, err
	}
	score, err := rag.UserScore(content, Cv)
	if err != nil {
		return false, err
	}
	s.repo.ApplicationSaver(JobId, userId, score)
	return true, nil
}
func (s *CVservice) jobAddService(Title, content string, userId uint) (bool, error) {
	var jobChunk []models.JobChunk

	chunks := rag.ChunkText(content, 500)
	vecs, err := rag.EmbedText(content)

	if err != nil {
		return false, err
	}

	job := models.Jobs{
		Title:    Title,
		Content:  content,
		Uploadby: userId,
	}
	s.repo.JobAdder(job)

	for i, vec := range vecs {
		jobChunk = append(jobChunk, models.JobChunk{
			JobID:   job.ID,
			Vector:  pgvector.NewVector(vec),
			Content: chunks[i],
		})
	}
	for _, chunk := range jobChunk {
		s.repo.JobChunkSaver(chunk)
	}

	return true, nil
}
