package user

import (
	"context"
	"strings"

	"github.com/Je33/imperial_fleet/internal/domain"
	"github.com/Je33/imperial_fleet/internal/repository/db/mysql"
	"github.com/Je33/imperial_fleet/internal/service"
	"gorm.io/gorm"

	"github.com/pkg/errors"
)

var (
	// errors prefix
	userErrorPrefix = "[repository.db.mysql.user]"

	// test interface
	_ service.UserRepository = (*UserMysqlRepo)(nil)
)

type UserMysqlRepo struct {
	db *mysql.DB
}

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"index"`
	Password  string
	CreatedAt int
	UpdatedAt int
}

func NewUserRepo(db *mysql.DB) *UserMysqlRepo {
	return &UserMysqlRepo{db}
}

func (repo *UserMysqlRepo) GetById(ctx context.Context, id uint) (*domain.User, error) {
	userDb := User{ID: id}
	err := repo.db.First(&userDb).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(domain.ErrNotFound, "%s: get by id", userErrorPrefix)
		}
		return nil, errors.Wrapf(err, "%s: get by id", userErrorPrefix)
	}
	return &domain.User{
		ID:        userDb.ID,
		Email:     userDb.Email,
		Password:  userDb.Password,
		CreatedAt: userDb.CreatedAt,
	}, nil
}

func (repo *UserMysqlRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	userDb := User{}
	err := repo.db.Where("lower(email) = ?", strings.ToLower(email)).First(&userDb).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(domain.ErrNotFound, "%s: get by email", userErrorPrefix)
		}
		return nil, errors.Wrapf(err, "%s: get by email", userErrorPrefix)
	}
	return &domain.User{
		ID:        userDb.ID,
		Email:     userDb.Email,
		Password:  userDb.Password,
		CreatedAt: userDb.CreatedAt,
	}, nil
}

func (repo *UserMysqlRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	userDb := User{
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	err := repo.db.Create(&userDb).Error
	if err != nil {
		return nil, errors.Wrapf(err, "%s: create", userErrorPrefix)
	}
	userDomain := &domain.User{
		ID:        userDb.ID,
		Email:     userDb.Email,
		Password:  userDb.Password,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
	}
	return userDomain, nil
}

func (repo *UserMysqlRepo) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	userQuery := User{ID: user.ID}
	userDb := User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		UpdatedAt: user.UpdatedAt,
	}
	err := repo.db.First(&userQuery).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(domain.ErrNotFound, "%s: get by id", userErrorPrefix)
	}
	err = repo.db.Save(&userDb).Error
	if err != nil {
		return nil, errors.Wrapf(err, "%s: update", userErrorPrefix)
	}
	userDomain := &domain.User{
		ID:        userDb.ID,
		Email:     userDb.Email,
		Password:  userDb.Password,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
	}
	return userDomain, nil
}
