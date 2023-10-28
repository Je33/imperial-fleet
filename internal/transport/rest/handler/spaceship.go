package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Je33/imperial_fleet/internal/domain"
	"github.com/Je33/imperial_fleet/internal/service"
	"github.com/Je33/imperial_fleet/internal/transport/rest/model"

	"github.com/labstack/echo/v4"
)

var (
	// errors prefix
	spaceshipErrorPrefix = "[transport.rest.handler.spaceship]"

	// test interface
	_ SpaceshipService = (*service.SpaceshipService)(nil)
)

//go:generate mockery --dir . --name SpaceshipService --output ./mocks
type SpaceshipService interface {
	GetAll(context.Context) ([]*domain.Spaceship, error)
	GetById(context.Context, uint) (*domain.Spaceship, error)
	CreateSpaceship(context.Context, *domain.Spaceship) error
	UpdateSpaceship(context.Context, *domain.Spaceship) error
	DeleteSpaceship(context.Context, *domain.Spaceship) error
}

type SpaceshipHandler struct {
	service SpaceshipService
}

func NewSpaceshipHandler(service SpaceshipService) *SpaceshipHandler {
	return &SpaceshipHandler{service}
}

func (h *SpaceshipHandler) GetAll(ctx echo.Context) error {

	spaceships, err := h.service.GetAll(ctx.Request().Context())

	if err != nil {
		return err
	}

	restSpaceships := make([]model.SpaceshipShort, 0, len(spaceships))
	for _, s := range spaceships {
		restSpaceships = append(restSpaceships, model.SpaceshipShort{
			ID:     s.ID,
			Name:   s.Name,
			Status: s.Status.String(),
		})
	}
	res := model.SpaceshipsResponce{
		Data: restSpaceships,
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *SpaceshipHandler) GetById(ctx echo.Context) error {

	idString := ctx.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		return err
	}

	if idInt < 0 {
		return domain.ErrNotFound
	}

	spaceship, err := h.service.GetById(ctx.Request().Context(), uint(idInt))
	if err != nil {
		return err
	}

	modelSpaceshipArmament := make([]model.SpaceshipArmament, 0, len(spaceship.Armament))
	for _, a := range spaceship.Armament {
		modelSpaceshipArmament = append(modelSpaceshipArmament, model.SpaceshipArmament{
			Title: a.Title,
			Qty:   a.Qty,
		})
	}

	restSpaceship := model.SpaceshipFull{
		ID:       spaceship.ID,
		Name:     spaceship.Name,
		Class:    spaceship.Class,
		Crew:     spaceship.Crew,
		Image:    spaceship.Image,
		Value:    spaceship.Value,
		Status:   spaceship.Status.String(),
		Armament: modelSpaceshipArmament,
	}
	return ctx.JSON(http.StatusOK, restSpaceship)
}

func (h *SpaceshipHandler) CreateSpaceship(ctx echo.Context) error {

	spaceship := new(model.SpaceshipFull)
	err := ctx.Bind(spaceship)
	if err != nil {
		return err
	}

	domainSpaceshipArmament := make([]domain.SpaceshipArmament, 0, len(spaceship.Armament))
	for _, a := range spaceship.Armament {
		domainSpaceshipArmament = append(domainSpaceshipArmament, domain.SpaceshipArmament{
			Title: a.Title,
			Qty:   a.Qty,
		})
	}

	domainSpaceship := &domain.Spaceship{
		Name:     spaceship.Name,
		Class:    spaceship.Class,
		Crew:     spaceship.Crew,
		Status:   domain.SpaceshipStatusFromString(spaceship.Status),
		Image:    spaceship.Image,
		Value:    spaceship.Value,
		Armament: domainSpaceshipArmament,
	}

	err = h.service.CreateSpaceship(ctx.Request().Context(), domainSpaceship)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, model.PostResponce{Success: true})
}

func (h *SpaceshipHandler) UpdateSpaceship(ctx echo.Context) error {

	spaceship := new(model.SpaceshipFull)
	err := ctx.Bind(spaceship)
	if err != nil {
		return err
	}

	idString := ctx.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		return err
	}

	if idInt < 0 {
		return domain.ErrNotFound
	}

	spaceship.ID = uint(idInt)

	domainSpaceshipArmament := make([]domain.SpaceshipArmament, 0, len(spaceship.Armament))
	for _, a := range spaceship.Armament {
		domainSpaceshipArmament = append(domainSpaceshipArmament, domain.SpaceshipArmament{
			Title: a.Title,
			Qty:   a.Qty,
		})
	}

	domainSpaceship := &domain.Spaceship{
		ID:       spaceship.ID,
		Name:     spaceship.Name,
		Class:    spaceship.Class,
		Crew:     spaceship.Crew,
		Status:   domain.SpaceshipStatusFromString(spaceship.Status),
		Image:    spaceship.Image,
		Value:    spaceship.Value,
		Armament: domainSpaceshipArmament,
	}

	err = h.service.UpdateSpaceship(ctx.Request().Context(), domainSpaceship)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, model.PostResponce{Success: true})
}

func (h *SpaceshipHandler) DeleteSpaceship(ctx echo.Context) error {

	idString := ctx.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		return err
	}

	if idInt < 0 {
		return domain.ErrNotFound
	}

	domainSpaceship := &domain.Spaceship{
		ID: uint(idInt),
	}

	fmt.Println(domainSpaceship)

	err = h.service.DeleteSpaceship(ctx.Request().Context(), domainSpaceship)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, model.PostResponce{Success: true})
}
