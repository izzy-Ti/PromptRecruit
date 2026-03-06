package cvs

import (
	rag "github.com/izzy-Ti/PromptRecruit/internals/Rag"
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
	ok, err, _, content := s.repo.GetJobByID(JobId)
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
