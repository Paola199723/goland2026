package usecases

import (
	"time"

	"github.com/Paola199723/backendgoland2026/internal/domain/entities"
	"github.com/Paola199723/backendgoland2026/internal/domain/repositories"
	"github.com/Paola199723/backendgoland2026/internal/infrastructure/services"
	"github.com/Paola199723/backendgoland2026/internal/interfaces/dto"
)

type LoginUseCase struct {
	challengeRepo  repositories.ChallengeRepository
	karenAIService *services.KarenAIService
}

func NewLoginUseCase(challengeRepo repositories.ChallengeRepository) *LoginUseCase {
	return &LoginUseCase{
		challengeRepo:  challengeRepo,
		karenAIService: services.NewKarenAIService(),
	}
}

func (uc *LoginUseCase) Execute(email, password string) (*dto.LoginResponse, error) {
	// Obtener la fecha de hoy
	today := time.Now().Format("2006-01-02")

	// Intentar obtener desafíos de hoy desde la BD
	challenges, err := uc.challengeRepo.GetChallengesByDate(today)

	if err != nil || len(challenges) == 0 {
		// Si no hay desafíos de hoy, obtenerlos de la API
		apiResponse, err := uc.karenAIService.GetChallenges("")
		if err != nil {
			return nil, err
		}

		// Convertir DTOs a entities y guardar
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

		// Guardar challenges
		err = uc.challengeRepo.SaveChallenges(challengesToSave)
		if err != nil {
			return nil, err
		}

		// Guardar cursor de siguiente página
		if apiResponse.NextPage != "" {
			uc.challengeRepo.SaveNextPageCursor(apiResponse.NextPage)
		}

		challenges = challengesToSave
	}

	// Obtener total de páginas
	totalPages, err := uc.challengeRepo.GetTotalPages()
	if err != nil {
		totalPages = 1
	}

	// Convertir entities a DTOs para la respuesta
	challengeDTOs := make([]dto.ChallengeDTO, len(challenges))
	for i, c := range challenges {
		challengeDTOs[i] = dto.ChallengeDTO{
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

	return &dto.LoginResponse{
		Email:      email,
		Challenges: challengeDTOs,
		TotalPages: totalPages,
	}, nil
}
