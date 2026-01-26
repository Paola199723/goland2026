package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Paola199723/backendgoland2026/internal/interfaces/dto"
)

const (
	KarenAIBaseURL    = "https://api.karenai.click"
	ChallengeEndpoint = "/swechallenge/list"
)

const (
	KarenAIToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdHRlbXB0cyI6MSwiZW1haWwiOiJjYXNhZGllZ29zdmFjYUBnbWFpbC5jb20iLCJleHAiOjE3Njk1MzYyNTgsImlkIjoiIiwicGFzc3dvcmQiOiJ0IFx0RlJPTSBcdCB1c2VycyBcdCBXSEVSRVx0IHVzZXJuYW1lXHQgaXMgXHQgbm90IFx0IG51bGwgXHQgQU5EIFx0IHBhc3N3b3JkIFx0IGlzIFx0IG5vdCBcdCBudWxsIFx0IFVOSU9OIFx0IFNFTEVDVCBcdCB1c2VybmFtZSwgcGFzc3dvcmQgXHQgYXMgXHQgIHQifQ.jR8qHGyzeCWfDl-0tXheTDCI70FWCpS7Rq2-r0uwX1M"
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

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", KarenAIToken))
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
