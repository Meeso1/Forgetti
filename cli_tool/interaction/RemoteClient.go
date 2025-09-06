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
		return nil, handleApiError(resp)
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
		return nil, handleApiError(resp)
	}

	var response dto.EncryptResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

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
