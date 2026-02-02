package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Paola199723/backendgoland2026/internal/interfaces/dto"
)

const (
	KarenAIBaseURL    = "https://api.karenai.click"
	ChallengeEndpoint = "/swechallenge/list"
)

type KarenAIService struct {
	client *http.Client
}

func NewKarenAIService() *KarenAIService {
	return &KarenAIService{
		client: &http.Client{},
	}
}

func (s *KarenAIService) GetChallenges(nextPage string) (*dto.ChallengeAPIResponse, error) {
	url := KarenAIBaseURL + ChallengeEndpoint

	if nextPage != "" {
		url = fmt.Sprintf("%s?next_page=%s", url, nextPage)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Obtener token desde variable de entorno
	karenAIToken := os.Getenv("TOKEN")
	if karenAIToken == "" {
		return nil, fmt.Errorf("TOKEN environment variable not set")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", karenAIToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d: %s", resp.StatusCode, string(body))
	}

	var response dto.ChallengeAPIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &response, nil
}
