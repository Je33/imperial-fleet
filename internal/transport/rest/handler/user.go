package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/Je33/imperial_fleet/internal/config"
	"github.com/Je33/imperial_fleet/internal/domain"
	"github.com/Je33/imperial_fleet/internal/service"
	"github.com/Je33/imperial_fleet/internal/transport/rest/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var (
	// errors prefix
	userErrorPrefix = "[transport.rest.handler.user]"

	// test interface
	_ UserService = (*service.UserService)(nil)
)

// jwt token struct
type jwtCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface {
	Auth(context.Context, *domain.UserAuthReq) error
	Register(context.Context, *domain.UserRegisterReq) error
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) Auth(ctx echo.Context) error {
	cfg := config.Get()

	restUserAuthReq := new(model.UserAuthReq)
	err := ctx.Bind(restUserAuthReq)
	if err != nil {
		return err
	}

	domainUserAuthReq := &domain.UserAuthReq{
		Email:    restUserAuthReq.Email,
		Password: restUserAuthReq.Password,
	}

	err = h.service.Auth(ctx.Request().Context(), domainUserAuthReq)
	if err != nil {
		return err
	}

	claims := &jwtCustomClaims{
		restUserAuthReq.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 168)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenSign, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return err
	}

	// TODO: Implement refresh strategy

	restUserAuthRes := &model.UserAuthRes{
		AuthToken:    tokenSign,
		RefreshToken: "",
	}

	return ctx.JSON(http.StatusOK, restUserAuthRes)
}

func (h *UserHandler) Register(ctx echo.Context) error {
	cfg := config.Get()

	restUserAuthReq := new(model.UserRegisterReq)
	err := ctx.Bind(restUserAuthReq)
	if err != nil {
		return err
	}

	if restUserAuthReq.Password != restUserAuthReq.RePassword {
		return domain.ErrRePasswordWrong
	}

	domainUserRegisterReq := &domain.UserRegisterReq{
		Email:      restUserAuthReq.Email,
		Password:   restUserAuthReq.Password,
		RePassword: restUserAuthReq.RePassword,
	}

	err = h.service.Register(ctx.Request().Context(), domainUserRegisterReq)
	if err != nil {
		return err
	}

	claims := &jwtCustomClaims{
		restUserAuthReq.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 168)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenSign, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return err
	}

	// TODO: Implement refresh strategy

	restUserAuthRes := &model.UserAuthRes{
		AuthToken:    tokenSign,
		RefreshToken: "",
	}

	return ctx.JSON(http.StatusOK, restUserAuthRes)
}
