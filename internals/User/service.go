package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	models "github.com/izzy-Ti/PromptRecruit/internals/Models"
	"github.com/izzy-Ti/PromptRecruit/internals/Utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *UserRepository
	jwtSecret []byte
}

func NewUserService(repo *UserRepository, secret string) *UserService {
	return &UserService{repo: repo, jwtSecret: []byte(secret)}
}

func (s *UserService) RegisterService(email, password, name string) (bool, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return false, errors.New("invalid email or password")
	}
	if user != nil {
		return false, errors.New("email already exist")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err := s.repo.RegisterUser(email, name, string(hashedPassword)); err != nil {
		return false, err
	}
	verificationUrl := ""
	subject := "Registration successful"
	html := fmt.Sprintf(`
        <p>Hi %s,</p>
        <p>Welcome! Please verify your email by clicking the link below:</p>
        <a href="%s">Verify Email</a>
        <p>If this wasn’t you, please ignore this email.</p>
    `, name, verificationUrl)
	err = Utils.Sendemail(email, name, subject, html)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (s *UserService) LoginService(email, password string) (*models.User, string, error) {
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
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, "", err
	}
	return user, tokenString, nil
}
func (s *UserService) SendVerifyOTPService(email string) (bool, error) {
	otp := Utils.GenerateOTP()
	expiresAt := time.Now().Add(24 * time.Hour).UnixMilli()

	user, err := s.repo.UpdateUserOTP(email, otp, int64(expiresAt))
	if err != nil {
		return false, err
	}
	subject := "OTP verfication"
	html := fmt.Sprintf(`
		<p>Hi %s,</p>
		<p>Your one-time verification code is:</p>
		<h2 style="letter-spacing:2px;">%s</h2>
		<p>This code will expire soon. Do not share it with anyone.</p>
		<p>If you didn’t request this, you can ignore this email.</p>
	`, user.Name, user.VerifyOTP)

	err = Utils.Sendemail(user.Email, user.Name, subject, html)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (s *UserService) VerifyOTPService(email, otp string) (bool, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return false, err
	}
	if user.VerifyOTP != otp || user.VerifyOTP == "" {
		return false, errors.New("Invalid otp")
	}
	if user.OTPExpireAt < time.Now().UnixMilli() {
		return false, errors.New("OTP expired")
	}
	subject := "Welcome! Your account has been verified"
	html := fmt.Sprintf(`
		<p>Hi %s,</p>
		<p>Thank you. Your account has been successfully verified.</p>
		<p>You can now sign in and start using your account.</p>
		<p>If you did not perform this action, you can reply to this email directly.</p>
	`, user.Name)

	err = Utils.Sendemail(user.Email, user.Name, subject, html)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (s *UserService) ValidateToken(tokenstr string) (*models.User, error) {
	token, err := jwt.Parse(tokenstr, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("sub not found")
	}
	user, err := s.repo.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	return user, nil
}
