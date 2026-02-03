package usecases

import (
	"log"
	"time"

	"github.com/Paola199723/backendgoland2026/internal/domain/entities"
	"github.com/Paola199723/backendgoland2026/internal/domain/repositories"
	"github.com/Paola199723/backendgoland2026/internal/infrastructure/services"
	"github.com/Paola199723/backendgoland2026/internal/interfaces/dto"
)

type GetChallengesUseCase struct {
	challengeRepo  repositories.ChallengeRepository
	karenAIService *services.KarenAIService
}

func NewGetChallengesUseCase(challengeRepo repositories.ChallengeRepository) *GetChallengesUseCase {
	return &GetChallengesUseCase{
		challengeRepo:  challengeRepo,
		karenAIService: services.NewKarenAIService(),
	}
}

func (uc *GetChallengesUseCase) ExecuteByPage(pageNum int) (*dto.LoginResponse, error) {
	log.Printf("🔍 GetChallengesUseCase.ExecuteByPage called with pageNum=%d", pageNum)
	
	// Obtener cursor de la página
	cursor, err := uc.challengeRepo.GetNextPageCursor(pageNum)
	log.Printf("🔍 GetNextPageCursor returned: cursor='%s', err=%v", cursor, err)

	var challenges []dto.ChallengeDTO
	var nextPage string

	if err == nil && cursor != "" {
		log.Printf("🔍 Using cursor from DB")
		// Si existe cursor, usarlo para obtener datos de la API
		apiResponse, err := uc.karenAIService.GetChallenges(cursor)
		if err != nil {
			log.Printf("⚠️ API call failed: %v", err)
			// Si falla, obtener de la BD
			dbChallenges, _ := uc.challengeRepo.GetChallengesByPage(pageNum)
			challenges = convertToDTO(dbChallenges)
		} else {
			log.Printf("✅ API call succeeded, got %d items", len(apiResponse.Items))
			// Convertir respuesta de API a DTOs
			nextPage = apiResponse.NextPage
			challenges = apiResponse.Items

			// Guardar challenges en BD
			var challengesToSave []entities.Challenge
			for _, c := range apiResponse.Items {
				// Parsear el time string a time.Time
				parsedTime, err := time.Parse(time.RFC3339Nano, c.Time)
				if err != nil {
					parsedTime = time.Now()
				}

				challengesToSave = append(challengesToSave, entities.Challenge{
					Ticker:     c.Ticker,
					TargetFrom: c.TargetFrom,
					TargetTo:   c.TargetTo,
					Company:    c.Company,
					Action:     c.Action,
					Brokerage:  c.Brokerage,
					RatingFrom: c.RatingFrom,
					RatingTo:   c.RatingTo,
					Time:       parsedTime,
				})
			}

			// Guardar en BD si no están ya guardadas
			uc.challengeRepo.SaveChallenges(challengesToSave)

			// Guardar cursor si existe
			if nextPage != "" {
				uc.challengeRepo.SaveNextPageCursor(nextPage)
			}
		}
	} else {
		log.Printf("🔍 No cursor from DB or cursor is empty. Trying DB...")
		// Obtener de la BD
		dbChallenges, err := uc.challengeRepo.GetChallengesByPage(pageNum)
		if err != nil {
			log.Printf("❌ DB error: %v", err)
			return nil, err
		}

		log.Printf("🔍 DB returned %d challenges", len(dbChallenges))

		// Si la BD está vacía, intentar obtener desde la API externa sin cursor
		if len(dbChallenges) == 0 {
			log.Printf("🔍 DB empty for page %d, calling external API without cursor", pageNum)
			apiResponse, err := uc.karenAIService.GetChallenges("")
			if err == nil && len(apiResponse.Items) > 0 {
				log.Printf("✅ External API succeeded: got %d items, next_page='%s'", len(apiResponse.Items), apiResponse.NextPage)
				// Convertir respuesta de API a DTOs y guardar en BD
				nextPage = apiResponse.NextPage
				challenges = apiResponse.Items

				var challengesToSave []entities.Challenge
				for _, c := range apiResponse.Items {
					parsedTime, err := time.Parse(time.RFC3339Nano, c.Time)
					if err != nil {
						parsedTime = time.Now()
					}

					challengesToSave = append(challengesToSave, entities.Challenge{
						Ticker:     c.Ticker,
						TargetFrom: c.TargetFrom,
						TargetTo:   c.TargetTo,
						Company:    c.Company,
						Action:     c.Action,
						Brokerage:  c.Brokerage,
						RatingFrom: c.RatingFrom,
						RatingTo:   c.RatingTo,
						Time:       parsedTime,
					})
				}

				saveErr := uc.challengeRepo.SaveChallenges(challengesToSave)
				log.Printf("🔍 SaveChallenges result: %v", saveErr)
				
				if nextPage != "" {
					uc.challengeRepo.SaveNextPageCursor(nextPage)
				}
			} else if err == nil && len(apiResponse.Items) == 0 {
				// Si API devuelve items vacíos, simplemente usar eso
				log.Printf("⚠️ External API returned empty list")
				challenges = apiResponse.Items
				nextPage = apiResponse.NextPage
			} else {
				// API fallo - usar datos de mock como fallback
				log.Printf("⚠️ External API failed: %v - using mock data", err)
				mockResponse := uc.karenAIService.GetMockChallenges("")
				challenges = mockResponse.Items
				nextPage = mockResponse.NextPage
			}
		} else {
			challenges = convertToDTO(dbChallenges)
		}
	}

	log.Printf("🔍 Returning %d challenges to handler", len(challenges))

	// Obtener total de páginas
	totalPages, err := uc.challengeRepo.GetTotalPages()
	if err != nil {
		totalPages = 1
	}

	return &dto.LoginResponse{
		Challenges: challenges,
		NextPage:   nextPage,
		TotalPages: totalPages,
	}, nil
}

// Función auxiliar para convertir entities a DTOs
func convertToDTO(challenges []entities.Challenge) []dto.ChallengeDTO {
	dtos := make([]dto.ChallengeDTO, len(challenges))
	for i, c := range challenges {
		dtos[i] = dto.ChallengeDTO{
			Ticker:     c.Ticker,
			TargetFrom: c.TargetFrom,
			TargetTo:   c.TargetTo,
			Company:    c.Company,
			Action:     c.Action,
			Brokerage:  c.Brokerage,
			RatingFrom: c.RatingFrom,
			RatingTo:   c.RatingTo,
			Time:       c.Time.Format(time.RFC3339Nano),
		}
	}
	return dtos
}
// ExecuteByPageCursor obtiene desafíos usando un cursor de paginación
func (uc *GetChallengesUseCase) ExecuteByPageCursor(nextPageCursor string) (*dto.LoginResponse, error) {
	log.Printf("🔍 GetChallengesUseCase.ExecuteByPageCursor called with cursor='%s'", nextPageCursor)
	
	// Obtener datos de la API usando el cursor
	apiResponse, err := uc.karenAIService.GetChallenges(nextPageCursor)
	if err != nil {
		log.Printf("❌ API call failed: %v", err)
		return &dto.LoginResponse{
			Challenges: []dto.ChallengeDTO{},
			NextPage:   "",
			TotalPages: 1,
		}, nil
	}

	log.Printf("✅ API call succeeded: got %d items, next_page='%s'", len(apiResponse.Items), apiResponse.NextPage)

	// Convertir respuesta de API a DTOs
	challenges := apiResponse.Items
	nextPage := apiResponse.NextPage

	// Guardar challenges en BD para caché
	var challengesToSave []entities.Challenge
	for _, c := range apiResponse.Items {
		parsedTime, err := time.Parse(time.RFC3339Nano, c.Time)
		if err != nil {
			parsedTime = time.Now()
		}

		challengesToSave = append(challengesToSave, entities.Challenge{
			Ticker:     c.Ticker,
			TargetFrom: c.TargetFrom,
			TargetTo:   c.TargetTo,
			Company:    c.Company,
			Action:     c.Action,
			Brokerage:  c.Brokerage,
			RatingFrom: c.RatingFrom,
			RatingTo:   c.RatingTo,
			Time:       parsedTime,
		})
	}

	uc.challengeRepo.SaveChallenges(challengesToSave)

	// Obtener total de páginas
	totalPages, err := uc.challengeRepo.GetTotalPages()
	if err != nil {
		totalPages = 1
	}

	return &dto.LoginResponse{
		Challenges: challenges,
		NextPage:   nextPage,
		TotalPages: totalPages,
	}, nil
}