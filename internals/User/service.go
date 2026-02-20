package user

import (
	"errors"

	models "github.com/izzy-Ti/PromptRecruit/internals/Models"
)

type UserService struct {
	repo       *UserRepository
	jwtSecrete []byte
}

func NewUserService(repo *UserRepository, secret string) *UserService {
	return &UserService{repo: repo, jwtSecrete: []byte(secret)}
}

func (s *UserService) Authenticate(email, password string) (*models.User, string, error) {
	if user, err := s.repo.GetByEmail(email); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

}
