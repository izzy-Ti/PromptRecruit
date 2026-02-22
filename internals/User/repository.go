package user

import (
	"errors"

	models "github.com/izzy-Ti/PromptRecruit/internals/Models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) CheckForEmail(email string) bool {
	var userModel models.User
	UserEmail := r.db.Where("email = ?", email).First(&userModel).Error
	if UserEmail == nil {
		return false
	}
	return true
}
func (r *UserRepository) RegisterUser(email, name, password string) error {
	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(password),
	}
	if err := r.db.Create(&user).Error; err != nil {
		return errors.New("Error saving user")
	}
	return nil
}
func (r *UserRepository) UpdateUserOTP(email, VerifyOTP string, OTPExpireAt int64) (*models.User, error) {
	user, err := r.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	user.VerifyOTP = VerifyOTP
	user.OTPExpireAt = OTPExpireAt
	r.db.Save(user)
	return user, nil
}
func (r *UserRepository) VerifyOTPRepo(email string) (*models.User, error) {
	user, err := r.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	user.VerifyOTP = ""
	user.OTPExpireAt = 0
	user.IsAccVerified = true
	r.db.Save(user)
	return user, nil
}
func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) SaveResetOTP(email, ResetOTP string, OTPExpireAt int64) (*models.User, error) {
	user, err := r.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	user.VerifyOTP = ResetOTP
	user.OTPExpireAt = OTPExpireAt
	r.db.Save(user)
	return user, nil
}
func (r *UserRepository) VerifyResetOTPRepo(email, password string) (*models.User, error) {
	user, err := r.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	user.ResetOTP = ""
	user.ResetOTPExpireAt = 0
	user.Password = password
	r.db.Save(user)
	return user, nil
}
