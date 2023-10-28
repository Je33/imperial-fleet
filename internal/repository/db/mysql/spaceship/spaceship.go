package spaceship

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Je33/imperial_fleet/internal/domain"
	"github.com/Je33/imperial_fleet/internal/repository/db/mysql"
	"github.com/Je33/imperial_fleet/internal/service"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	// errors prefix
	spaceshipErrorPrefix = "[repository.db.mysql.spaceship]"

	// test interface
	_ service.SpaceshipRepository = (*SpaceshipMysqlRepo)(nil)
)

// spaceship repo
type SpaceshipMysqlRepo struct {
	db *mysql.DB
}

// models for orm
// many 2 many relation:
// spaceship_armaments -> spaceship_armament_qties <- spaceships

// spaceship_armaments table
type SpaceshipArmament struct {
	ID    uint   `gorm:"primaryKey"`
	Title string `gorm:"size:256;uniqueIndex"`

	// supress field for orm
	Qty uint `gorm:"-"`
}

// spaceship_armament_qties table
type SpaceshipArmamentQty struct {
	SpaceshipID         uint `gorm:"index:,unique,composite:myname"`
	SpaceshipArmamentID uint `gorm:"index:,unique,composite:myname"`
	Qty                 uint
}

// spaceships table
type Spaceship struct {
	ID       uint                `gorm:"primaryKey"`
	Name     string              `gorm:"size:256;uniqueIndex"`
	Class    string              `gorm:"size:256"`
	Armament []SpaceshipArmament `gorm:"many2many:spaceship_armament_qties;"`
	Crew     uint
	Image    string `gorm:"size:256"`
	Value    float64
	Status   uint
}

// spaceship repo builder
func NewSpaceshipRepo(db *mysql.DB) *SpaceshipMysqlRepo {
	return &SpaceshipMysqlRepo{db}
}

// get all spaceships from db with short info
func (repo *SpaceshipMysqlRepo) GetAll(ctx context.Context) ([]*domain.Spaceship, error) {

	spaceships := []Spaceship{}

	// get all records from db
	// TODO: make filters
	res := repo.db.Select("id", "name", "status").Find(&spaceships)
	if res.Error != nil {
		return nil, errors.Wrapf(res.Error, "%s: get all", spaceshipErrorPrefix)
	}

	// convert db records to domain level
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

// get one spaceship from db with detailed info
func (repo *SpaceshipMysqlRepo) GetById(ctx context.Context, id uint) (*domain.Spaceship, error) {

	// create mini model for orm query
	spaceshipDb := Spaceship{ID: id}
	err := repo.db.First(&spaceshipDb).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// if not found return domain not found error
			return nil, errors.Wrapf(domain.ErrNotFound, "%s: get by id", spaceshipErrorPrefix)
		}
		// else return error as is
		return nil, errors.Wrapf(err, "%s: get by id", spaceshipErrorPrefix)
	}

	// convert db spaceship armaments to domain level
	domainSpaceshipArmaments := []domain.SpaceshipArmament{}
	repo.db.Raw(`
		SELECT sa.id, sa.title, saq.qty FROM spaceship_armaments sa
		INNER JOIN spaceship_armament_qties saq ON sa.id = saq.spaceship_armament_id AND saq.spaceship_id = ?
	`, spaceshipDb.ID).Scan(&domainSpaceshipArmaments)

	return &domain.Spaceship{
		ID:       spaceshipDb.ID,
		Name:     spaceshipDb.Name,
		Class:    spaceshipDb.Class,
		Crew:     spaceshipDb.Crew,
		Image:    spaceshipDb.Image,
		Armament: domainSpaceshipArmaments,
		Value:    spaceshipDb.Value,
		Status:   domain.SpaceshipStatus(spaceshipDb.Status),
	}, nil
}

// create spaceship
func (repo *SpaceshipMysqlRepo) Create(ctx context.Context, spaceship *domain.Spaceship) error {

	// TODO: make transaction wrap

	// make map with armament quantities
	spaceshipArmamentMap := make(map[string]uint)
	spiceshipArmamentDb := make([]SpaceshipArmament, 0, len(spaceship.Armament))
	for _, a := range spaceship.Armament {
		spaceshipArmamentMap[a.Title] = a.Qty
		spiceshipArmamentDb = append(spiceshipArmamentDb, SpaceshipArmament{
			Title: a.Title,
		})
	}

	// ensure that all new armaments exist in db
	err := repo.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&spiceshipArmamentDb).Error
	if err != nil {
		return errors.Wrapf(err, "%s: create", spaceshipErrorPrefix)
	}

	// find all requested armaments
	err = repo.db.Find(&spiceshipArmamentDb).Error
	if err != nil {
		return errors.Wrapf(err, "%s: create", spaceshipErrorPrefix)
	}

	// create spaceship db model
	spaceshipDb := Spaceship{
		Name:   spaceship.Name,
		Class:  spaceship.Class,
		Crew:   spaceship.Crew,
		Status: uint(spaceship.Status),
		Image:  spaceship.Image,
		Value:  spaceship.Value,
	}

	// save spaceship model to db
	err = repo.db.Create(&spaceshipDb).Error
	if err != nil {
		return errors.Wrapf(err, "%s: create", spaceshipErrorPrefix)
	}

	// build all quantites and armaments
	spaceshipArmamentQtyDb := make([]SpaceshipArmamentQty, 0, len(spaceship.Armament))
	for _, a := range spiceshipArmamentDb {
		spaceshipArmamentQtyDb = append(spaceshipArmamentQtyDb, SpaceshipArmamentQty{
			SpaceshipID:         spaceshipDb.ID,
			SpaceshipArmamentID: a.ID,
			Qty:                 spaceshipArmamentMap[a.Title],
		})
	}

	// save armaments with quantities
	err = repo.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "spaceship_id"}, {Name: "spaceship_armament_id"}},
		UpdateAll: true,
	}).Create(&spaceshipArmamentQtyDb).Error
	if err != nil {
		return errors.Wrapf(err, "%s: create", spaceshipErrorPrefix)
	}

	return nil
}

func (repo *SpaceshipMysqlRepo) Update(ctx context.Context, spaceship *domain.Spaceship) error {

	// TODO: make transaction wrap
	// TODO: DRY

	// check if spaceship exists in db
	spaceshipQuery := Spaceship{ID: spaceship.ID}
	err := repo.db.First(&spaceshipQuery).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrapf(domain.ErrNotFound, "%s: update", spaceshipErrorPrefix)
	}

	// make map with armament quantities
	spaceshipArmamentMap := make(map[string]uint)
	spiceshipArmamentDb := make([]SpaceshipArmament, 0, len(spaceship.Armament))
	for _, a := range spaceship.Armament {
		spaceshipArmamentMap[a.Title] = a.Qty
		spiceshipArmamentDb = append(spiceshipArmamentDb, SpaceshipArmament{
			Title: a.Title,
		})
	}

	// ensure that all new armaments exist in db
	err = repo.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&spiceshipArmamentDb).Error
	if err != nil {
		return errors.Wrapf(err, "%s: create", spaceshipErrorPrefix)
	}

	// find all requested armaments
	err = repo.db.Find(&spiceshipArmamentDb).Error
	if err != nil {
		return errors.Wrapf(err, "%s: create", spaceshipErrorPrefix)
	}

	// create spaceship db model
	spaceshipDb := Spaceship{
		Name:   spaceship.Name,
		Class:  spaceship.Class,
		Crew:   spaceship.Crew,
		Status: uint(spaceship.Status),
		Image:  spaceship.Image,
		Value:  spaceship.Value,
	}

	// save spaceship model to db
	err = repo.db.Model(&spaceshipQuery).Updates(spaceshipDb).Error
	if err != nil {
		return errors.Wrapf(err, "%s: create", spaceshipErrorPrefix)
	}

	// build all quantites and armaments
	spaceshipArmamentQtyDb := make([]SpaceshipArmamentQty, 0, len(spaceship.Armament))
	for _, a := range spiceshipArmamentDb {
		spaceshipArmamentQtyDb = append(spaceshipArmamentQtyDb, SpaceshipArmamentQty{
			SpaceshipID:         spaceshipQuery.ID,
			SpaceshipArmamentID: a.ID,
			Qty:                 spaceshipArmamentMap[a.Title],
		})
	}

	// save armaments with quantities
	err = repo.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "spaceship_id"}, {Name: "spaceship_armament_id"}},
		UpdateAll: true,
	}).Create(&spaceshipArmamentQtyDb).Error
	if err != nil {
		return errors.Wrapf(err, "%s: create", spaceshipErrorPrefix)
	}

	return nil
}

// delete spaceship and related armaments with quantities
func (repo *SpaceshipMysqlRepo) Delete(ctx context.Context, spaceship *domain.Spaceship) error {

	// TODO: soft delete

	// check if spaceship exists in db
	spaceshipQuery := Spaceship{ID: spaceship.ID}
	err := repo.db.First(&spaceshipQuery).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrapf(domain.ErrNotFound, "%s: delete get by id", spaceshipErrorPrefix)
	}

	// delete spaceship
	err = repo.db.Delete(&spaceshipQuery).Error
	if err != nil {
		return errors.Wrapf(err, "%s: delete spaceship", spaceshipErrorPrefix)
	}

	// delete armaments with quantities
	err = repo.db.Where("spaceship_id = ?", spaceshipQuery.ID).Delete(&SpaceshipArmamentQty{}).Error
	if err != nil {
		return errors.Wrapf(err, "%s: delete spaceship armament qty", spaceshipErrorPrefix)
	}

	return nil
}
