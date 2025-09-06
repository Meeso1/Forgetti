package interaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forgetti-common/constants"
	"forgetti-common/dto"
	"net/http"
	"time"
)

type RemoteClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewRemoteClient(baseURL string) *RemoteClient {
	return &RemoteClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (r *RemoteClient) NewKey(content string, expiration time.Time) (*dto.NewKeyResponse, error) {
	request := dto.NewKeyRequest{
		Content:    content,
		Expiration: expiration,
	}

	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := r.httpClient.Post(
		r.baseURL + constants.NewKeyRoute,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var response dto.NewKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (r *RemoteClient) Encrypt(content string, keyId string) (*dto.EncryptResponse, error) {
	request := dto.EncryptRequest{
		Content: content,
		KeyId:   keyId,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := r.httpClient.Post(
		r.baseURL + constants.EncryptRoute,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var response dto.EncryptResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
