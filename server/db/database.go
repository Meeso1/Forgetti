package db

import (
	"ForgettiServer/config"
	"ForgettiServer/db/models"
	"fmt"
	"forgetti-common/io"
	"reflect"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

type DatabaseService struct {
	db *gorm.DB
}

// TODO: Add method for encryption at rest (separate service, use encryption key from config)
func CreateDb(cfg *config.Config) (*gorm.DB, error) {
	path, err := io.GetRelativePathFromBin(cfg.Database.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path from bin: %w", err)
	}

	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        path,
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	for _, model := range models.ModelsToMigrate {
		if err := db.AutoMigrate(model); err != nil {
			return nil, fmt.Errorf("failed to run database migrations for model %T: %w", reflect.TypeOf(model), err)
		}
	}

	return db, nil
}

func NewDatabaseService(db *gorm.DB) *DatabaseService {
	return &DatabaseService{
		db: db,
	}
}

func (ds *DatabaseService) Close() error {
	sqlDB, err := ds.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (ds *DatabaseService) HealthCheck() error {
	sqlDB, err := ds.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
