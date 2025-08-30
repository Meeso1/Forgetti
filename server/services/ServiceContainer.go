package services

type ServiceContainer struct {
	Encryptor Encryptor
	KeyStore KeyStore
}

func CreateServiceContainer() *ServiceContainer {
	keyStore := CreateKeyStore()
	encryptor := CreateEncryptor(keyStore)

	return &ServiceContainer{
		Encryptor: encryptor,
		KeyStore: keyStore,
	}
}
