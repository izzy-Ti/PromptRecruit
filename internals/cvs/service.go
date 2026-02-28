package cvs

type CVservice struct {
	repo *CvRepo
}

func NewUserService(repo *CvRepo) *CVservice {
	return &CVservice{repo: repo}
}

func CvUploadSvc() {

}
