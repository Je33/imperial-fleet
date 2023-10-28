package service

import (
	"context"
	"time"

	"github.com/Je33/imperial_fleet/internal/domain"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	userErrorPrefix = "[service.user]"
)

//go:generate mockery --dir . --name UserRepository --output ./mocks
type UserRepository interface {
	GetByEmail(context.Context, string) (*domain.User, error)
	Create(context.Context, *domain.User) (*domain.User, error)
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{repository}
}

func (s *UserService) Register(ctx context.Context, req *domain.UserRegisterReq) error {

	_, err := s.repository.GetByEmail(ctx, req.Email)

	// if user exists
	if !errors.Is(err, domain.ErrNotFound) {
		return domain.ErrUserExists
	}

	// if repeated password not match
	if req.Password != req.RePassword {
		return domain.ErrRePasswordWrong
	}

	// create new user and encode password
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return errors.Wrapf(err, "%s: encode password error", userErrorPrefix)
	}
	newUser := &domain.User{
		Email:    req.Email,
		Password: string(passwordBytes),
		CreatedAt: time.Now().UTC().Second(),
		UpdatedAt: time.Now().UTC().Second(),
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

	// if user not exists
	if errors.Is(err, domain.ErrNotFound) {
		return domain.ErrNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return domain.ErrPasswordWrong
	}

	return nil
}
