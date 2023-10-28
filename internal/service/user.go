package service

import (
	"context"
	"time"

	"github.com/Je33/imperial_fleet/internal/domain"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	// prefix for wrap errors
	userErrorPrefix = "[service.user]"
)

//go:generate mockery --dir . --name UserRepository --output ./mocks
type UserRepository interface {
	GetByEmail(context.Context, string) (*domain.User, error)
	Create(context.Context, *domain.User) (*domain.User, error)
}

// user service
type UserService struct {
	repository UserRepository
}

// user service builder
func NewUserService(repository UserRepository) *UserService {
	return &UserService{repository}
}

// user registration
func (s *UserService) Register(ctx context.Context, req *domain.UserRegisterReq) error {

	// required fields
	if req.Email == "" || req.Password == "" {
		return domain.ErrRegRequiredFields
	}

	// if repassword and password are not match
	if req.Password != req.RePassword {
		return domain.ErrRePasswordWrong
	}

	_, err := s.repository.GetByEmail(ctx, req.Email)

	// if user exists
	if !errors.Is(err, domain.ErrNotFound) {
		return domain.ErrUserExists
	}

	// encode password
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return errors.Wrapf(err, "%s: encode password error", userErrorPrefix)
	}

	// save user
	newUser := &domain.User{
		Email:     req.Email,
		Password:  string(passwordBytes),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	_, err = s.repository.Create(ctx, newUser)
	if err != nil {
		return errors.Wrapf(err, "%s: repo save error", userErrorPrefix)
	}

	return nil
}

func (s *UserService) Auth(ctx context.Context, req *domain.UserAuthReq) error {

	// find user and compare password
	user, err := s.repository.GetByEmail(ctx, req.Email)

	// if user not exists or other errors
	if err != nil {
		return err
	}

	// compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return domain.ErrPasswordWrong
	}

	return nil
}
