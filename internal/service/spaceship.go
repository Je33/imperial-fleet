package service

import (
	"context"

	"github.com/Je33/imperial_fleet/internal/domain"
	"github.com/pkg/errors"
)

var (
	// prefix for wrap errors
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

// spaceship service
type SpaceshipService struct {
	repository SpaceshipRepository
}

// spaceship service builder
func NewSpaceshipService(repository SpaceshipRepository) *SpaceshipService {
	return &SpaceshipService{repository}
}

// get list of all spaceships
func (s *SpaceshipService) GetAll(ctx context.Context) ([]*domain.Spaceship, error) {

	// get all spaceships
	spaceships, err := s.repository.GetAll(ctx)

	if err != nil {
		return nil, errors.Wrapf(err, "%s: get all spaceships error", spaceshipErrorPrefix)
	}

	return spaceships, nil
}

func (s *SpaceshipService) GetById(ctx context.Context, id uint) (*domain.Spaceship, error) {

	spaceship, err := s.repository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return spaceship, nil
}

// create spaceship record
func (s *SpaceshipService) CreateSpaceship(ctx context.Context, spaceship *domain.Spaceship) error {

	// name of spaceship is required
	if spaceship.Name == "" {
		return domain.ErrNameRequired
	}

	// create spaceship record in repo db
	err := s.repository.Create(ctx, spaceship)
	if err != nil {
		return err
	}

	return nil
}


// update spaceship record
func (s *SpaceshipService) UpdateSpaceship(ctx context.Context, spaceship *domain.Spaceship) error {

	// name of spaceship is required
	if spaceship.Name == "" {
		return domain.ErrNameRequired
	}

	// update spaceship record in repo db
	err := s.repository.Update(ctx, spaceship)
	if err != nil {
		return err
	}

	return nil
}

// delete spaceship record
func (s *SpaceshipService) DeleteSpaceship(ctx context.Context, spaceship *domain.Spaceship) error {

	// delete spaceship record and all related records in repo db
	err := s.repository.Delete(ctx, spaceship)
	if err != nil {
		return err
	}

	return nil
}
