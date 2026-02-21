package user

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	models "github.com/izzy-Ti/PromptRecruit/internals/Models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo       *UserRepository
	jwtSecrete []byte
}

func NewUserService(repo *UserRepository, secret string) *UserService {
	return &UserService{repo: repo, jwtSecrete: []byte(secret)}
}

func (s *UserService) Login(email, password string) (*models.User, string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}
	if !user.IsAccVerified {
		return nil, "", errors.New("please verify your account")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("incorrect credentials")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(s.jwtSecrete)
	if err != nil {
		return nil, "", err
	}
	return user, tokenString, nil
}
