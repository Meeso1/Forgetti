package routes

import (
	"ForgettiServer/services"
	"fmt"
	"forgetti-common/constants"
	"forgetti-common/crypto"
	"forgetti-common/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewKeyRoute(c *gin.Context, s *services.ServiceContainer) {
	var request dto.NewKeyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newKey, err := s.Encryptor.CreateNewKeyAndEncrypt(request.Content, request.Expiration)
	if err != nil {
		// TODO: Improve
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	verificationKey, err := crypto.SerializePrivateKey(newKey.VerificationKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("failed to serialize verification key: %w", err).Error()})
		return
	}

	response := dto.NewKeyResponse{
		EncryptedContent: newKey.EncryptedContent,
		Metadata: dto.Metadata{
			KeyId:           newKey.KeyId,
			Expiration:      newKey.Expiration,
			VerificationKey: verificationKey,
		},
	}

	c.JSON(http.StatusOK, response)
}

func EncryptRoute(c *gin.Context, s *services.ServiceContainer) {
	var request dto.EncryptRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encryptedContent, err := s.Encryptor.EncryptWithExistingKey(request.Content, request.KeyId)
	if err != nil {
		// TODO: Improve
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.EncryptResponse{
		EncryptedContent: encryptedContent,
	}

	c.JSON(http.StatusOK, response)
}

func DecryptRoute(c *gin.Context, s *services.ServiceContainer) {
	var request dto.DecryptRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decryptedContent, err := s.Encryptor.Decrypt(request.EncryptedContent, request.VerificationKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"decrypted_content": decryptedContent})
}

func AddEncRoutes(router *gin.Engine, serviceContainer *services.ServiceContainer) {
	router.POST(constants.NewKeyRoute, func(c *gin.Context) { NewKeyRoute(c, serviceContainer) })
	router.POST(constants.EncryptRoute, func(c *gin.Context) { EncryptRoute(c, serviceContainer) })
	router.POST(constants.DecryptRoute, func(c *gin.Context) { DecryptRoute(c, serviceContainer) })
}
