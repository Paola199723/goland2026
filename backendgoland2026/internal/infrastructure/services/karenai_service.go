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
	if karenAIToken == "" || karenAIToken == "your_karenai_token_here" {
		fmt.Println("⚠️ WARNING: TOKEN not configured. Returning mock data for testing.")
		return s.GetMockChallenges(nextPage), nil
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

// GetMockChallenges devuelve datos de prueba cuando el TOKEN no está configurado
// Soporta paginación con cursor para demostrar la funcionalidad
func (s *KarenAIService) GetMockChallenges(nextPage string) *dto.ChallengeAPIResponse {
	// Página 1 (sin cursor o cursor vacío)
	if nextPage == "" {
		return &dto.ChallengeAPIResponse{
			Items: []dto.ChallengeDTO{
				{
					Ticker:     "AAPL",
					TargetFrom: "$150",
					TargetTo:   "$160",
					Company:    "Apple Inc.",
					Action:     "BUY",
					Brokerage:  "Interactive Brokers",
					RatingFrom: "HOLD",
					RatingTo:   "BUY",
					Time:       "2026-02-02T10:30:00Z",
				},
				{
					Ticker:     "GOOGL",
					TargetFrom: "$140",
					TargetTo:   "$150",
					Company:    "Alphabet Inc.",
					Action:     "HOLD",
					Brokerage:  "Interactive Brokers",
					RatingFrom: "BUY",
					RatingTo:   "BUY",
					Time:       "2026-02-02T11:00:00Z",
				},
				{
					Ticker:     "MSFT",
					TargetFrom: "$380",
					TargetTo:   "$400",
					Company:    "Microsoft Corporation",
					Action:     "BUY",
					Brokerage:  "Interactive Brokers",
					RatingFrom: "HOLD",
					RatingTo:   "BUY",
					Time:       "2026-02-02T11:30:00Z",
				},
				{
					Ticker:     "TSLA",
					TargetFrom: "$200",
					TargetTo:   "$220",
					Company:    "Tesla Inc.",
					Action:     "SELL",
					Brokerage:  "Interactive Brokers",
					RatingFrom: "BUY",
					RatingTo:   "HOLD",
					Time:       "2026-02-02T12:00:00Z",
				},
			},
			NextPage: "page_2", // Cursor para la siguiente página
		}
	}

	// Página 2 (cursor = "page_2")
	if nextPage == "page_2" {
		return &dto.ChallengeAPIResponse{
			Items: []dto.ChallengeDTO{
				{
					Ticker:     "AMZN",
					TargetFrom: "$170",
					TargetTo:   "$190",
					Company:    "Amazon.com Inc.",
					Action:     "BUY",
					Brokerage:  "Morgan Stanley",
					RatingFrom: "EQUAL-WEIGHT",
					RatingTo:   "OVERWEIGHT",
					Time:       "2026-02-02T13:00:00Z",
				},
				{
					Ticker:     "NFLX",
					TargetFrom: "$380",
					TargetTo:   "$420",
					Company:    "Netflix Inc.",
					Action:     "BUY",
					Brokerage:  "Goldman Sachs",
					RatingFrom: "NEUTRAL",
					RatingTo:   "BUY",
					Time:       "2026-02-02T13:30:00Z",
				},
				{
					Ticker:     "META",
					TargetFrom: "$280",
					TargetTo:   "$310",
					Company:    "Meta Platforms Inc.",
					Action:     "BUY",
					Brokerage:  "JPMorgan",
					RatingFrom: "NEUTRAL",
					RatingTo:   "OVERWEIGHT",
					Time:       "2026-02-02T14:00:00Z",
				},
				{
					Ticker:     "NVIDIA",
					TargetFrom: "$710",
					TargetTo:   "$780",
					Company:    "NVIDIA Corporation",
					Action:     "BUY",
					Brokerage:  "Barclays",
					RatingFrom: "EQUAL-WEIGHT",
					RatingTo:   "OVERWEIGHT",
					Time:       "2026-02-02T14:30:00Z",
				},
			},
			NextPage: "page_3", // Cursor para la tercera página
		}
	}

	// Página 3 (cursor = "page_3")
	if nextPage == "page_3" {
		return &dto.ChallengeAPIResponse{
			Items: []dto.ChallengeDTO{
				{
					Ticker:     "INTEL",
					TargetFrom: "$85",
					TargetTo:   "$95",
					Company:    "Intel Corporation",
					Action:     "HOLD",
					Brokerage:  "Credit Suisse",
					RatingFrom: "UNDERPERFORM",
					RatingTo:   "NEUTRAL",
					Time:       "2026-02-02T15:00:00Z",
				},
				{
					Ticker:     "AMD",
					TargetFrom: "$130",
					TargetTo:   "$150",
					Company:    "Advanced Micro Devices Inc.",
					Action:     "BUY",
					Brokerage:  "Mizuho Securities",
					RatingFrom: "OUTPERFORM",
					RatingTo:   "OUTPERFORM",
					Time:       "2026-02-02T15:30:00Z",
				},
			},
			NextPage: "", // Última página, sin siguiente cursor
		}
	}

	// Fallback
	return &dto.ChallengeAPIResponse{
		Items:    []dto.ChallengeDTO{},
		NextPage: "",
	}
}
