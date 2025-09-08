package services

import (
	"ForgettiServer/config"
	"ForgettiServer/db"
	"ForgettiServer/db/repositories"
)

type ServiceContainer struct {
	Config              *config.Config
	DatabaseService     *db.DatabaseService
	KeyRepo             *repositories.KeyRepo
	RecentlyExpiredRepo *repositories.RecentlyExpiredRepo
	Encryptor           Encryptor
	KeyStore            KeyStore
	DataProtection      DataProtection
}

func CreateServiceContainer(cfg *config.Config) (*ServiceContainer, error) {
	database, err := db.CreateDb(cfg)
	if err != nil {
		return nil, err
	}

	databaseService := db.NewDatabaseService(database)
	keyRepo := repositories.NewKeyRepo(database)
	recentlyExpiredRepo := repositories.NewRecentlyExpiredRepo(database)
	
	dataProtection, err := NewDataProtection(cfg)
	if err != nil {
		return nil, err
	}

	keyStore := NewKeyStore(keyRepo, recentlyExpiredRepo, dataProtection, cfg)
	encryptor := CreateEncryptor(keyStore)


	return &ServiceContainer{
		Config:              cfg,
		DatabaseService:     databaseService,
		KeyRepo:             keyRepo,
		RecentlyExpiredRepo: recentlyExpiredRepo,
		DataProtection:      dataProtection,
		KeyStore:            keyStore,
		Encryptor:           encryptor,
	}, nil
}
