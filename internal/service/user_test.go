package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Je33/imperial_fleet/internal/domain"
	"github.com/Je33/imperial_fleet/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Register(t *testing.T) {

	testCases := []struct {
		name         string
		expectations func(context.Context, *mocks.UserRepository)
		input        *domain.UserRegisterReq
		err          error
	}{
		{
			name: "success registration",
			input: &domain.UserRegisterReq{
				Email:      "test@test.com",
				Password:   "123123",
				RePassword: "123123",
			},
			expectations: func(ctx context.Context, userRepo *mocks.UserRepository) {
				userRepo.On("GetByEmail", ctx, "test@test.com").Return(nil, domain.ErrNotFound)
				userRepo.On("Create", ctx, mock.Anything).Return(&domain.User{}, nil)
			},
			err: nil,
		},
		{
			name: "failed registration requred fields",
			input: &domain.UserRegisterReq{
				Email:      "",
				Password:   "",
				RePassword: "123",
			},
			expectations: func(ctx context.Context, userRepo *mocks.UserRepository) {
				//
			},
			err: domain.ErrRegRequiredFields,
		},
		{
			name: "failed registration repassword",
			input: &domain.UserRegisterReq{
				Email:      "test@test.com",
				Password:   "123123",
				RePassword: "123",
			},
			expectations: func(ctx context.Context, userRepo *mocks.UserRepository) {
				//
			},
			err: domain.ErrRePasswordWrong,
		},
		{
			name: "failed registration user exists",
			input: &domain.UserRegisterReq{
				Email:      "test@test.com",
				Password:   "123123",
				RePassword: "123123",
			},
			expectations: func(ctx context.Context, userRepo *mocks.UserRepository) {
				userRepo.On("GetByEmail", ctx, "test@test.com").Return(nil, nil)
			},
			err: domain.ErrUserExists,
		},
		{
			name: "failed registration save error",
			input: &domain.UserRegisterReq{
				Email:      "test@test.com",
				Password:   "123123",
				RePassword: "123123",
			},
			expectations: func(ctx context.Context, userRepo *mocks.UserRepository) {
				userRepo.On("GetByEmail", ctx, "test@test.com").Return(nil, domain.ErrNotFound)
				userRepo.On("Create", ctx, mock.Anything).Return(nil, errors.New("error"))
			},
			err: errors.New("error"),
		},
	}

	for _, test := range testCases {
		t.Logf("testing %s", test.name)

		ctx := context.Background()

		userRepo := mocks.NewUserRepository(t)
		userService := NewUserService(userRepo)

		test.expectations(ctx, userRepo)

		err := userService.Register(ctx, test.input)

		if err != nil {
			if test.err != nil {
				assert.Error(t, test.err, err)
			} else {
				assert.NoError(t, err)
			}
		}

		userRepo.AssertExpectations(t)

	}
}

func TestUserService_Auth(t *testing.T) {

	userAuthReq := &domain.UserAuthReq{
		Email:    "test@test.com",
		Password: "123123",
	}

	testCases := []struct {
		name         string
		expectations func(context.Context, *mocks.UserRepository)
		input        *domain.UserAuthReq
		err          error
	}{
		{
			name: "success auth",
			input: &domain.UserAuthReq{
				Email:    "test@test.com",
				Password: "123123",
			},
			expectations: func(ctx context.Context, userRepo *mocks.UserRepository) {
				userRepo.On("GetByEmail", ctx, userAuthReq.Email).Return(&domain.User{
					Email:     "test@test.com",
					Password:  "$2a$14$yar6jY8fAfR85i0.KmVH1OFkkfuJLGr2o5uMPu3p7Iae2xspPIpAu",
					CreatedAt: 1,
					UpdatedAt: 1,
				}, nil)
			},
			err: nil,
		},
		{
			name: "failed auth password wrong",
			input: &domain.UserAuthReq{
				Email:    "test@test.com",
				Password: "123123",
			},
			expectations: func(ctx context.Context, userRepo *mocks.UserRepository) {
				userRepo.On("GetByEmail", ctx, userAuthReq.Email).Return(&domain.User{
					Email:     "test@test.com",
					Password:  "$2a$$yar6jY8fAfR85i0.KmVH1OFkkfuJLGr2o5uMPu3p7Iae2xspPIpAu",
					CreatedAt: 1,
					UpdatedAt: 1,
				}, nil)
			},
			err: domain.ErrPasswordWrong,
		},
		{
			name: "failed auth user not found",
			input: &domain.UserAuthReq{
				Email:    "test@test.com",
				Password: "123123",
			},
			expectations: func(ctx context.Context, userRepo *mocks.UserRepository) {
				userRepo.On("GetByEmail", ctx, userAuthReq.Email).Return(nil, domain.ErrNotFound)
			},
			err: domain.ErrNotFound,
		},
	}

	for _, test := range testCases {
		t.Logf("testing %s", test.name)

		ctx := context.Background()

		userRepo := mocks.NewUserRepository(t)
		userService := NewUserService(userRepo)

		test.expectations(ctx, userRepo)

		err := userService.Auth(ctx, test.input)

		if err != nil {
			if test.err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.NoError(t, err)
			}
		}

		userRepo.AssertExpectations(t)

	}
}
