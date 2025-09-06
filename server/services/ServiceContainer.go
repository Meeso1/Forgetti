package services

import (
	"ForgettiServer/config"
)

type ServiceContainer struct {
	Config    *config.Config
	Encryptor Encryptor
	KeyStore  KeyStore
}

func CreateServiceContainer(cfg *config.Config) *ServiceContainer {
	keyStore := CreateKeyStore(cfg)
	encryptor := CreateEncryptor(keyStore)

	return &ServiceContainer{
		Config:    cfg,
		Encryptor: encryptor,
		KeyStore:  keyStore,
	}
}
