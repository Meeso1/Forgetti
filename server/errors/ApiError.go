package errors

import (
	"fmt"
	"net/http"
	"forgetti-common/dto"
	"time"
)

type ApiError struct {
	Message string
	ErrorCode string
	StatusCode int
	Data map[string]string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode, e.Message)
}

func (e *ApiError) ToResponse() *dto.ErrorResponse {
	return &dto.ErrorResponse{
		Message: e.Message,
		ErrorCode: e.ErrorCode,
		Data: e.Data,
	}
}

func KeyNotFoundError(keyId string) *ApiError {
	return &ApiError{
		Message: fmt.Sprintf("key not found: %s", keyId),
		ErrorCode: "key-not-found",
		StatusCode: http.StatusNotFound,
		Data: map[string]string{
			"key_id": keyId,
		},
	}
}

func KeyExpiredError(keyId string, expiration time.Time) *ApiError {
	return &ApiError{
		Message: fmt.Sprintf("key %s expired at %s", keyId, expiration.Format(time.RFC3339)),
		ErrorCode: "key-expired",
		StatusCode: http.StatusNotFound,
		Data: map[string]string{
			"key_id": keyId,
			"expiration": expiration.Format(time.RFC3339),
		},
	}
}

func BadRequestError(err error) *ApiError {
	return &ApiError{
		Message: fmt.Sprintf("failed to parse request: %s", err.Error()),
		ErrorCode: "bad-request",
		StatusCode: http.StatusBadRequest,
		Data: map[string]string{
			"error": err.Error(),
		},
	}
}

func InternalServerError(err error) *ApiError {
	return &ApiError{
		Message: fmt.Sprintf("internal server error: %s", err.Error()),
		ErrorCode: "internal-server-error",
		StatusCode: http.StatusInternalServerError,
		Data: map[string]string{
			"error": err.Error(),
		},
	}
}