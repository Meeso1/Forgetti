package routes

import (
	"ForgettiServer/services"
	"fmt"
	"forgetti-common/constants"
	"forgetti-common/crypto"
	"forgetti-common/dto"
	apiErrors "ForgettiServer/errors"

	"github.com/gin-gonic/gin"
)

func newKeyRoute(c *gin.Context, s *services.ServiceContainer) (*dto.NewKeyResponse, error) {
	var request dto.NewKeyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		return nil, apiErrors.BadRequestError(err)
	}

	if err := request.Validate(); err != nil {
		return nil, apiErrors.BadRequestError(err)
	}

	newKey, err := s.Encryptor.CreateNewKeyAndEncrypt(request.Content, request.Expiration)
	if err != nil {
		return nil, err
	}

	verificationKey, err := crypto.SerializePrivateKey(newKey.VerificationKey)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize verification key: %w", err)
	}

	response := dto.NewKeyResponse{
		EncryptedContent: newKey.EncryptedContent,
		Metadata: dto.Metadata{
			KeyId:           newKey.KeyId,
			Expiration:      newKey.Expiration,
			VerificationKey: verificationKey,
		},
	}

	return &response, nil
}

func encryptRoute(c *gin.Context, s *services.ServiceContainer) (*dto.EncryptResponse, error) {
	var request dto.EncryptRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		return nil, apiErrors.BadRequestError(err)
	}

	encryptedContent, err := s.Encryptor.EncryptWithExistingKey(request.Content, request.KeyId)
	if err != nil {
		return nil, err
	}

	response := dto.EncryptResponse{
		EncryptedContent: encryptedContent,
	}

	return &response, nil
}

func AddEncRoutes(router *gin.Engine, serviceContainer *services.ServiceContainer) {
	router.POST(constants.NewKeyRoute, createEndpoint(serviceContainer, newKeyRoute))
	router.POST(constants.EncryptRoute, createEndpoint(serviceContainer, encryptRoute))
}
