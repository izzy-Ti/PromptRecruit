package cvs

type CVservice struct {
	repo *CvRepo
}

func NewUserService(repo *CvRepo) *CVservice {
	return &CVservice{repo: repo}
}

func CvUploadSvc() {

}
func (s *CVservice) ApplicationService(userId, JobId uint) (bool, error) {
	ok, err, jobVec := s.repo.GetJobByID(JobId)
	if !ok {
		return false, err
	}
	ok, err, Cvec := s.repo.GetUserCv(userId)
	if !ok {
		return false, err
	}

	return true, nil
}
