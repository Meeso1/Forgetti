package routes

import (
	apiErrors "ForgettiServer/errors"
	"ForgettiServer/services"
	"fmt"
	"forgetti-common/constants"
	"forgetti-common/crypto"
	"forgetti-common/dto"
	"forgetti-common/logging"

	"github.com/gin-gonic/gin"
)

func newKeyRoute(c *gin.Context, s *services.ServiceContainer) (*dto.NewKeyResponse, error) {
	logger := logging.MakeLogger("routes.newKeyRoute")
	logger.Verbose("Received new key request from %s", c.ClientIP())

	var request dto.NewKeyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("Failed to bind JSON request: %v", err)
		return nil, apiErrors.BadRequestError(err)
	}
	logger.Verbose("Request bound successfully, expiration: %s", request.Expiration.Format("2006-01-02 15:04:05"))

	logger.Verbose("Validating new key request")
	if err := request.Validate(); err != nil {
		logger.Error("Request validation failed: %v", err)
		return nil, apiErrors.BadRequestError(err)
	}
	logger.Verbose("Request validation successful")

	logger.Verbose("Calling Encryptor to create new key and encrypt")
	newKey, err := s.Encryptor.CreateNewKeyAndEncrypt(request.Content, request.Expiration)
	if err != nil {
		logger.Error("Failed to create new key and encrypt: %v", err)
		return nil, err
	}
	logger.Verbose("New key created and content encrypted successfully")

	logger.Verbose("Serializing verification key")
	verificationKey, err := crypto.SerializePrivateKey(newKey.VerificationKey)
	if err != nil {
		logger.Error("Failed to serialize verification key: %v", err)
		return nil, fmt.Errorf("failed to serialize verification key: %w", err)
	}
	logger.Verbose("Verification key serialized successfully")

	response := dto.NewKeyResponse{
		EncryptedContent: newKey.EncryptedContent,
		Metadata: dto.Metadata{
			KeyId:           newKey.KeyId,
			Expiration:      newKey.Expiration,
			VerificationKey: verificationKey,
		},
	}

	logger.Info("New key request completed successfully. KeyId: %s, Client: %s", newKey.KeyId, c.ClientIP())
	return &response, nil
}

func encryptRoute(c *gin.Context, s *services.ServiceContainer) (*dto.EncryptResponse, error) {
	logger := logging.MakeLogger("routes.encryptRoute")
	logger.Verbose("Received encrypt request from %s", c.ClientIP())

	var request dto.EncryptRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("Failed to bind JSON request: %v", err)
		return nil, apiErrors.BadRequestError(err)
	}
	logger.Verbose("Request bound successfully, KeyId: %s", request.KeyId)

	logger.Verbose("Calling Encryptor to encrypt with existing key")
	encryptedContent, err := s.Encryptor.EncryptWithExistingKey(request.Content, request.KeyId)
	if err != nil {
		logger.Error("Failed to encrypt with existing key: %v", err)
		return nil, err
	}
	logger.Verbose("Content encrypted successfully with existing key")

	response := dto.EncryptResponse{
		EncryptedContent: encryptedContent,
	}

	logger.Info("Encrypt request completed successfully. KeyId: %s, Client: %s", request.KeyId, c.ClientIP())
	return &response, nil
}

func AddEncRoutes(router *gin.Engine, serviceContainer *services.ServiceContainer) {
	logger := logging.MakeLogger("routes.AddEncRoutes")
	logger.Verbose("Adding route: POST %s", constants.NewKeyRoute)
	router.POST(constants.NewKeyRoute, createEndpoint(serviceContainer, newKeyRoute))
	logger.Verbose("Adding route: POST %s", constants.EncryptRoute)
	router.POST(constants.EncryptRoute, createEndpoint(serviceContainer, encryptRoute))
	logger.Verbose("Encryption routes added successfully")
}
