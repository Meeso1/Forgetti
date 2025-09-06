package routes

import (
	"ForgettiServer/services"
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	apiErrors "ForgettiServer/errors"
)

func createEndpoint[T any](
	container *services.ServiceContainer,
	f func(c *gin.Context, s *services.ServiceContainer) (*T, error),
) func(c *gin.Context) {
	return func(c *gin.Context) {
		result, err := f(c, container)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, *result)
	}
}

func handleError(c *gin.Context, err error) {
	var apiError *apiErrors.ApiError
	if errors.As(err, &apiError) {
		c.JSON(apiError.StatusCode, apiError.ToResponse())
		return
	}

	c.JSON(http.StatusInternalServerError, apiErrors.InternalServerError(err).ToResponse())
}