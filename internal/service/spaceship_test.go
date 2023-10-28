package service

import (
	"context"
	"testing"

	"github.com/Je33/imperial_fleet/internal/domain"
	"github.com/Je33/imperial_fleet/internal/service/mocks"

	"github.com/stretchr/testify/assert"
)

func TestSpaceshipService_GetAll(t *testing.T) {

	spaceships := []*domain.Spaceship{
		{
			ID:     1,
			Name:   "Devastator",
			Status: domain.SpaceshipStatusOperational,
		},
	}

	testCases := []struct {
		name         string
		expectations func(context.Context, *mocks.SpaceshipRepository)
		err          error
	}{
		{
			name: "success get all",
			expectations: func(ctx context.Context, spaceshipRepo *mocks.SpaceshipRepository) {
				spaceshipRepo.On("GetAll", ctx).Return(spaceships, nil)
			},
			err: nil,
		},
		{
			name: "failed get all",
			expectations: func(ctx context.Context, spaceshipRepo *mocks.SpaceshipRepository) {
				spaceshipRepo.On("GetAll", ctx).Return(nil, domain.ErrNotFound)
			},
			err: domain.ErrNotFound,
		},
	}

	for _, test := range testCases {
		t.Logf("testing %s", test.name)

		ctx := context.Background()

		spaceshipRepo := mocks.NewSpaceshipRepository(t)
		spaceshipService := NewSpaceshipService(spaceshipRepo)

		test.expectations(ctx, spaceshipRepo)

		_, err := spaceshipService.GetAll(ctx)

		if err != nil {
			if test.err != nil {
				assert.Error(t, test.err, err)
			} else {
				assert.NoError(t, err)
			}
		}

		spaceshipRepo.AssertExpectations(t)

	}
}

func TestSpaceshipService_GetById(t *testing.T) {

	spaceship := &domain.Spaceship{
		ID:     1,
		Name:   "Devastator",
		Status: domain.SpaceshipStatusOperational,
	}
	var id uint = 1

	testCases := []struct {
		name         string
		expectations func(context.Context, *mocks.SpaceshipRepository)
		err          error
	}{
		{
			name: "success get by id",
			expectations: func(ctx context.Context, spaceshipRepo *mocks.SpaceshipRepository) {
				spaceshipRepo.On("GetById", ctx, id).Return(spaceship, nil)
			},
			err: nil,
		},
		{
			name: "failed get by id",
			expectations: func(ctx context.Context, spaceshipRepo *mocks.SpaceshipRepository) {
				spaceshipRepo.On("GetById", ctx, id).Return(nil, domain.ErrNotFound)
			},
			err: domain.ErrNotFound,
		},
	}

	for _, test := range testCases {
		t.Logf("testing %s", test.name)

		ctx := context.Background()

		spaceshipRepo := mocks.NewSpaceshipRepository(t)
		spaceshipService := NewSpaceshipService(spaceshipRepo)

		test.expectations(ctx, spaceshipRepo)

		_, err := spaceshipService.GetById(ctx, id)

		if err != nil {
			if test.err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.NoError(t, err)
			}
		}

		spaceshipRepo.AssertExpectations(t)

	}
}

func TestSpaceshipService_CreateSpaceship(t *testing.T) {

	spaceship := &domain.Spaceship{
		ID:     1,
		Name:   "Devastator",
		Status: domain.SpaceshipStatusOperational,
	}

	testCases := []struct {
		name         string
		expectations func(context.Context, *mocks.SpaceshipRepository)
		err          error
	}{
		{
			name: "success create spaceship",
			expectations: func(ctx context.Context, spaceshipRepo *mocks.SpaceshipRepository) {
				spaceshipRepo.On("Create", ctx, spaceship).Return(nil)
			},
			err: nil,
		},
		{
			name: "failed create spaceship",
			expectations: func(ctx context.Context, spaceshipRepo *mocks.SpaceshipRepository) {
				spaceshipRepo.On("Create", ctx, spaceship).Return(domain.ErrNotFound)
			},
			err: domain.ErrNotFound,
		},
	}

	for _, test := range testCases {
		t.Logf("testing %s", test.name)

		ctx := context.Background()

		spaceshipRepo := mocks.NewSpaceshipRepository(t)
		spaceshipService := NewSpaceshipService(spaceshipRepo)

		test.expectations(ctx, spaceshipRepo)

		err := spaceshipService.CreateSpaceship(ctx, spaceship)

		if err != nil {
			if test.err != nil {
				assert.Error(t, test.err, err)
			} else {
				assert.NoError(t, err)
			}
		}

		spaceshipRepo.AssertExpectations(t)

	}
}

func TestSpaceshipService_UpdateSpaceship(t *testing.T) {

	spaceship := &domain.Spaceship{
		ID:     1,
		Name:   "Devastator",
		Status: domain.SpaceshipStatusOperational,
	}

	testCases := []struct {
		name         string
		expectations func(context.Context, *mocks.SpaceshipRepository)
		err          error
	}{
		{
			name: "success update spaceship",
			expectations: func(ctx context.Context, spaceshipRepo *mocks.SpaceshipRepository) {
				spaceshipRepo.On("Update", ctx, spaceship).Return(nil)
			},
			err: nil,
		},
		{
			name: "failed update spaceship",
			expectations: func(ctx context.Context, spaceshipRepo *mocks.SpaceshipRepository) {
				spaceshipRepo.On("Update", ctx, spaceship).Return(domain.ErrNotFound)
			},
			err: domain.ErrNotFound,
		},
	}

	for _, test := range testCases {
		t.Logf("testing %s", test.name)

		ctx := context.Background()

		spaceshipRepo := mocks.NewSpaceshipRepository(t)
		spaceshipService := NewSpaceshipService(spaceshipRepo)

		test.expectations(ctx, spaceshipRepo)

		err := spaceshipService.UpdateSpaceship(ctx, spaceship)

		if err != nil {
			if test.err != nil {
				assert.Error(t, test.err, err)
			} else {
				assert.NoError(t, err)
			}
		}

		spaceshipRepo.AssertExpectations(t)

	}
}

func TestSpaceshipService_DeleteSpaceship(t *testing.T) {

	spaceship := &domain.Spaceship{
		ID:     1,
		Name:   "Devastator",
		Status: domain.SpaceshipStatusOperational,
	}

	testCases := []struct {
		name         string
		expectations func(context.Context, *mocks.SpaceshipRepository)
		err          error
	}{
		{
			name: "success delete spaceship",
			expectations: func(ctx context.Context, spaceshipRepo *mocks.SpaceshipRepository) {
				spaceshipRepo.On("Delete", ctx, spaceship).Return(nil)
			},
			err: nil,
		},
		{
			name: "failed delete spaceship",
			expectations: func(ctx context.Context, spaceshipRepo *mocks.SpaceshipRepository) {
				spaceshipRepo.On("Delete", ctx, spaceship).Return(domain.ErrNotFound)
			},
			err: domain.ErrNotFound,
		},
	}

	for _, test := range testCases {
		t.Logf("testing %s", test.name)

		ctx := context.Background()

		spaceshipRepo := mocks.NewSpaceshipRepository(t)
		spaceshipService := NewSpaceshipService(spaceshipRepo)

		test.expectations(ctx, spaceshipRepo)

		err := spaceshipService.DeleteSpaceship(ctx, spaceship)

		if err != nil {
			if test.err != nil {
				assert.Error(t, test.err, err)
			} else {
				assert.NoError(t, err)
			}
		}

		spaceshipRepo.AssertExpectations(t)

	}
}
