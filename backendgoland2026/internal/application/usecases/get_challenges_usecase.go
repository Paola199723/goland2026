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
	// Obtener cursor de la página
	cursor, err := uc.challengeRepo.GetNextPageCursor(pageNum)

	var challenges []dto.ChallengeDTO
	var nextPage string

	if err == nil && cursor != "" {
		// Si existe cursor, usarlo para obtener datos de la API
		apiResponse, err := uc.karenAIService.GetChallenges(cursor)
		if err != nil {
			// Si falla, obtener de la BD
			dbChallenges, _ := uc.challengeRepo.GetChallengesByPage(pageNum)
			challenges = convertToDTO(dbChallenges)
		} else {
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
		// Obtener de la BD
		dbChallenges, err := uc.challengeRepo.GetChallengesByPage(pageNum)
		if err != nil {
			return nil, err
		}

		// Si la BD está vacía, intentar obtener desde la API externa sin cursor
		if len(dbChallenges) == 0 {
			log.Printf("GetChallengesUseCase: DB empty for page %d, calling external API", pageNum)
			apiResponse, err := uc.karenAIService.GetChallenges("")
			if err == nil {
				log.Printf("GetChallengesUseCase: external API returned %d items, next_page=%s", len(apiResponse.Items), apiResponse.NextPage)
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

				uc.challengeRepo.SaveChallenges(challengesToSave)
				if nextPage != "" {
					uc.challengeRepo.SaveNextPageCursor(nextPage)
				}
			} else {
				challenges = convertToDTO(dbChallenges)
			}
		} else {
			challenges = convertToDTO(dbChallenges)
		}
	}

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
