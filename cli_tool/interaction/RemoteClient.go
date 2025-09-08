package interaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forgetti-common/constants"
	"forgetti-common/dto"
	"forgetti-common/logging"
	"net/http"
	"time"
)

type RemoteClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewRemoteClient(baseURL string) *RemoteClient {
	logger := logging.MakeLogger("RemoteClient.NewRemoteClient")
	logger.Verbose("Creating new remote client for server: %s", baseURL)
	return &RemoteClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (r *RemoteClient) NewKey(content string, expiration time.Time) (*dto.NewKeyResponse, error) {
	logger := logging.MakeLogger("RemoteClient.NewKey")
	request := dto.NewKeyRequest{
		Content:    content,
		Expiration: expiration,
	}

	logger.Verbose("Validating new key request")
	if err := request.Validate(); err != nil {
		logger.Error("New key request validation failed: %v", err)
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	logger.Verbose("New key request validation successful")

	logger.Verbose("Marshaling new key request to JSON")
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := r.baseURL + constants.NewKeyRoute
	logger.Verbose("Making HTTP POST request to: %s", url)
	resp, err := r.httpClient.Post(
		url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		logger.Error("HTTP POST request failed: %v", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	logger.Verbose("HTTP response received with status code: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		logger.Error("HTTP request failed with status code: %d", resp.StatusCode)
		return nil, handleApiError(resp)
	}

	logger.Verbose("Decoding new key response from JSON")
	var response dto.NewKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logger.Error("Failed to decode new key response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	logger.Info("Successfully created new key. KeyId: %s", response.Metadata.KeyId)
	return &response, nil
}

func (r *RemoteClient) Encrypt(content string, keyId string) (*dto.EncryptResponse, error) {
	logger := logging.MakeLogger("RemoteClient.Encrypt")
	request := dto.EncryptRequest{
		Content: content,
		KeyId:   keyId,
	}

	logger.Verbose("Marshaling encrypt request to JSON for KeyId: %s", keyId)
	jsonData, err := json.Marshal(request)
	if err != nil {
		logger.Error("Failed to marshal encrypt request: %v", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := r.baseURL + constants.EncryptRoute
	logger.Verbose("Making HTTP POST request to: %s", url)
	resp, err := r.httpClient.Post(
		url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		logger.Error("HTTP POST request failed for encrypt: %v", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	logger.Verbose("HTTP response received with status code: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		logger.Error("HTTP encrypt request failed with status code: %d", resp.StatusCode)
		return nil, handleApiError(resp)
	}

	logger.Verbose("Decoding encrypt response from JSON")
	var response dto.EncryptResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logger.Error("Failed to decode encrypt response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	logger.Info("Successfully completed encrypt request for KeyId: %s", keyId)
	return &response, nil
}

func handleApiError(resp *http.Response) error {
	var response dto.ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if err := makePrettyError(response); err != nil {
		return err
	}

	return fmt.Errorf("[%d] %s: %s", resp.StatusCode, response.ErrorCode, response.Message)
}

func makePrettyError(response dto.ErrorResponse) error {
	switch response.ErrorCode {
	case "key-not-found":
		return fmt.Errorf("key %s does not exist on server - it could have expired, or another server was used to generate it", response.Data["key_id"])
	case "key-expired":
		return fmt.Errorf("key %s expired at %s", response.Data["key_id"], response.Data["expiration"])
	case "bad-request":
		return fmt.Errorf("request failed: %s", response.Data["error"])
	case "internal-server-error":
		return fmt.Errorf("server error: %s", response.Message)
	default:
		return nil
	}
}
