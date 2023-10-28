package spaceship

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Je33/imperial_fleet/internal/domain"
	"github.com/Je33/imperial_fleet/internal/repository/db/mysql"
	"github.com/Je33/imperial_fleet/internal/service"
	"gorm.io/gorm"
)

var (
	// errors prefix
	spaceshipErrorPrefix = "[repository.db.mysql.spaceship]"

	// test interface
	_ service.SpaceshipRepository = (*SpaceshipMysqlRepo)(nil)
)

type SpaceshipMysqlRepo struct {
	db *mysql.DB
}

type SpaceshipArmament struct {
	ID    uint
	Title string
	Qty   uint `gorm:"-"`
}

type SpaceshipArmamentQty struct {
	SpaceshipID         uint
	SpaceshipArmamentID uint
	Qty                 uint
}

type Spaceship struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Class    string
	Armament []SpaceshipArmament `gorm:"many2many:spaceship_armament_qty;"`
	Crew     uint
	Image    string
	Value    float64
	Status   uint
}

func NewSpaceshipRepo(db *mysql.DB) *SpaceshipMysqlRepo {
	return &SpaceshipMysqlRepo{db}
}

func (repo *SpaceshipMysqlRepo) GetAll(ctx context.Context) ([]*domain.Spaceship, error) {
	spaceships := []Spaceship{}
	res := repo.db.Find(&spaceships)
	if res.Error != nil {
		return nil, errors.Wrapf(res.Error, "%s: get all", spaceshipErrorPrefix)
	}
	domainSpaceships := make([]*domain.Spaceship, 0, res.RowsAffected)
	for _, ss := range spaceships {
		domainSpaceships = append(domainSpaceships, &domain.Spaceship{
			ID:     ss.ID,
			Name:   ss.Name,
			Status: domain.SpaceshipStatus(ss.Status),
		})
	}
	return domainSpaceships, nil
}

func (repo *SpaceshipMysqlRepo) GetById(ctx context.Context, id uint) (*domain.Spaceship, error) {
	spaceshipDb := Spaceship{ID: id}
	err := repo.db.Preload("SpaceshipArmament").First(&spaceshipDb).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(domain.ErrNotFound, "%s: get by id", spaceshipErrorPrefix)
		}
		return nil, errors.Wrapf(err, "%s: get by id", spaceshipErrorPrefix)
	}
	domainSpaceshipArmament := make([]domain.SpaceshipArmament, 0, len(spaceshipDb.Armament))
	for _, sa := range spaceshipDb.Armament {
		domainSpaceshipArmament = append(domainSpaceshipArmament, domain.SpaceshipArmament{
			Title: sa.Title,
			Qty:   sa.Qty,
		})
	}
	return &domain.Spaceship{
		ID:       spaceshipDb.ID,
		Name:     spaceshipDb.Name,
		Class:    spaceshipDb.Class,
		Crew:     spaceshipDb.Crew,
		Image:    spaceshipDb.Image,
		Armament: domainSpaceshipArmament,
		Value:    spaceshipDb.Value,
		Status:   domain.SpaceshipStatus(spaceshipDb.Status),
	}, nil
}

func (repo *SpaceshipMysqlRepo) Create(ctx context.Context, spaceship *domain.Spaceship) error {
	spaceshipArmamentDb := make([]SpaceshipArmament, 0, len(spaceship.Armament))
	for _, a := range spaceship.Armament {
		spaceshipArmamentDb = append(spaceshipArmamentDb, SpaceshipArmament{
			Title: a.Title,
			Qty:   a.Qty,
		})
	}
	spaceshipDb := Spaceship{
		Name:     spaceship.Name,
		Class:    spaceship.Class,
		Crew:     spaceship.Crew,
		Status:   uint(spaceship.Status),
		Image:    spaceship.Image,
		Value:    spaceship.Value,
		Armament: spaceshipArmamentDb,
	}
	err := repo.db.Create(&spaceshipDb).Error
	if err != nil {
		return errors.Wrapf(err, "%s: create", spaceshipErrorPrefix)
	}
	return nil
}

func (repo *SpaceshipMysqlRepo) Update(ctx context.Context, spaceship *domain.Spaceship) error {
	spaceshipQuery := Spaceship{ID: spaceship.ID}
	err := repo.db.First(&spaceshipQuery).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrapf(domain.ErrNotFound, "%s: get by id", spaceshipErrorPrefix)
	}
	spaceshipArmamentDb := make([]SpaceshipArmament, 0, len(spaceship.Armament))
	for _, a := range spaceship.Armament {
		spaceshipArmamentDb = append(spaceshipArmamentDb, SpaceshipArmament{
			Title: a.Title,
			Qty:   a.Qty,
		})
	}
	spaceshipDb := Spaceship{
		ID:       spaceship.ID,
		Name:     spaceship.Name,
		Class:    spaceship.Class,
		Crew:     spaceship.Crew,
		Status:   uint(spaceship.Status),
		Image:    spaceship.Image,
		Value:    spaceship.Value,
		Armament: spaceshipArmamentDb,
	}
	err = repo.db.Save(&spaceshipDb).Error
	if err != nil {
		return errors.Wrapf(err, "%s: update", spaceshipErrorPrefix)
	}
	return nil
}

func (repo *SpaceshipMysqlRepo) Delete(ctx context.Context, spaceship *domain.Spaceship) error {
	spaceshipQuery := Spaceship{ID: spaceship.ID}
	err := repo.db.First(&spaceshipQuery).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrapf(domain.ErrNotFound, "%s: get by id", spaceshipErrorPrefix)
	}
	err = repo.db.Delete(&spaceshipQuery).Error
	if err != nil {
		return errors.Wrapf(err, "%s: delete", spaceshipErrorPrefix)
	}
	return nil
}
