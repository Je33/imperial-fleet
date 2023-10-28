package service

import (
	"context"

	"github.com/Je33/imperial_fleet/internal/domain"
	"github.com/pkg/errors"
)

var (
	spaceshipErrorPrefix = "[service.spaceship]"
)

//go:generate mockery --dir . --name SpaceshipRepository --output ./mocks
type SpaceshipRepository interface {
	GetAll(context.Context) ([]*domain.Spaceship, error)
	GetById(context.Context, uint) (*domain.Spaceship, error)
	Create(context.Context, *domain.Spaceship) error
	Update(context.Context, *domain.Spaceship) error
	Delete(context.Context, *domain.Spaceship) error
}

type SpaceshipService struct {
	repository SpaceshipRepository
}

func NewSpaceshipService(repository SpaceshipRepository) *SpaceshipService {
	return &SpaceshipService{repository}
}

func (s *SpaceshipService) GetAll(ctx context.Context) ([]*domain.Spaceship, error) {

	// get all spaceships
	spaceships, err := s.repository.GetAll(ctx)

	if err != nil {
		return nil, errors.Wrapf(err, "%s: get all spaceships error", spaceshipErrorPrefix)
	}

	return spaceships, nil
}

func (s *SpaceshipService) GetById(ctx context.Context, id uint) (*domain.Spaceship, error) {

	// get spaceship by id
	spaceship, err := s.repository.GetById(ctx, id)

	// if spaceship not exists
	if errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrNotFound
	}

	return spaceship, nil
}

func (s *SpaceshipService) CreateSpaceship(ctx context.Context, spaceship *domain.Spaceship) error {
	err := s.repository.Create(ctx, spaceship)
	if err != nil {
		return err
	}
	return nil
}

func (s *SpaceshipService) UpdateSpaceship(ctx context.Context, spaceship *domain.Spaceship) error {
	err := s.repository.Update(ctx, spaceship)
	if err != nil {
		return err
	}
	return nil
}

func (s *SpaceshipService) DeleteSpaceship(ctx context.Context, spaceship *domain.Spaceship) error {
	err := s.repository.Delete(ctx, spaceship)
	if err != nil {
		return err
	}
	return nil
}
