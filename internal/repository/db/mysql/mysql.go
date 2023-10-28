package mysql

import (
	"context"

	"github.com/Je33/imperial_fleet/internal/config"
	"github.com/Je33/imperial_fleet/internal/domain"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

var (
	mysqlErrorPrefix = "[repository.db.mysql]"
)

func Connect(ctx context.Context) (*DB, error) {
	cfg := config.Get()
	if cfg.MysqlDSN == "" {
		return nil, errors.Wrapf(domain.ErrConfig, "%s: connection url not provided", mysqlErrorPrefix)
	}

	// Create a new client and connect to the server
	client, err := gorm.Open(mysql.Open(cfg.MysqlDSN), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrapf(err, "%s: disconnected", mysqlErrorPrefix)
	}

	res := &DB{client}

	// TODO: Send a ping to confirm a successful connection
	// err = res.Ping(ctx)
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "%s: disconnected", mysqlErrorPrefix)
	// }

	return res, nil
}

func (db *DB) Ping(ctx context.Context) error {
	// Send a ping to check connection
	err := db.Ping(ctx)
	if err != nil {
		return errors.Wrapf(err, "%s: disconnected", mysqlErrorPrefix)
	}
	return nil
}
