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
	jobVec := 
	cv, err := s.repo.GetTopMatchingCvs()
	return true, nil
}
func (s *CVservice) ApplicationService(){

}