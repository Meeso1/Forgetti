package dto

import (
	"errors"
	"fmt"
	"time"
)

const maxExpiration time.Duration = 30 * 24 * time.Hour

type NewKeyRequest struct {
	Content    string    `json:"content" binding:"required,min=1,max=1000"`
	Expiration time.Time `json:"expiration" binding:"required"`
}

func (r NewKeyRequest) Validate() error {
	if r.Expiration.Before(time.Now()) {
		return errors.New("expiration must be in the future")
	}

	if r.Expiration.After(time.Now().Add(maxExpiration)) {
		return fmt.Errorf("expiration must be less than %s in the future", maxExpiration.String())
	}

	return nil
}
