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
	log.Printf(" GetChallengesUseCase.ExecuteByPage called with pageNum=%d", pageNum)

	// PASO 1: Verificar si hay desafíos en la BD (caché)
	allChallenges, err := uc.challengeRepo.GetChallengesByPage(1)
	if err != nil {
		log.Printf("Error querying all challenges: %v", err)
	}

	var challenges []dto.ChallengeDTO
	var nextPage string

	if len(allChallenges) > 0 {
		// PASO 2a: Si hay datos en caché, usarlos
		log.Printf(" Found %d challenges in cache", len(allChallenges))

		challenges = convertToDTO(allChallenges)

		// Obtener el cursor guardado
		cursor, err := uc.challengeRepo.GetCursorForToday()
		if err != nil {
			log.Printf(" Error getting cursor: %v", err)
		}
		nextPage = cursor

	} else {
		// PASO 2b: Si no hay datos en caché, consultar /swechallenge/list
		log.Printf(" No challenges in cache, calling /swechallenge/list")

		apiResponse, err := uc.karenAIService.GetChallenges("")
		if err != nil {
			log.Printf(" API call failed: %v", err)
			return &dto.LoginResponse{
				Challenges: []dto.ChallengeDTO{},
				NextPage:   "",
				TotalPages: 1,
			}, nil
		}

		log.Printf("API call succeeded: got %d items, next_page='%s'", len(apiResponse.Items), apiResponse.NextPage)

		// PASO 3: Guardar lo que retorna /swechallenge/list en la tabla challenges
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
		log.Printf("SaveChallenges result: %v", saveErr)

		// PASO 4: Guardar el next_page en page_cursors
		if apiResponse.NextPage != "" {
			cursorErr := uc.challengeRepo.SaveCursorForDate(apiResponse.NextPage, time.Now().Format("2006-01-02"))
			log.Printf("SaveCursorForDate result: %v", cursorErr)
		}

		// PASO 5: Retornar al frontend
		challenges = apiResponse.Items
		nextPage = apiResponse.NextPage
	}

	log.Printf(" Returning %d challenges to handler with next_page='%s'", len(challenges), nextPage)

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
// Guarda los resultados en la BD y actualiza el next_page en page_cursors
func (uc *GetChallengesUseCase) ExecuteByPageCursor(nextPageCursor string) (*dto.LoginResponse, error) {
	log.Printf(" GetChallengesUseCase.ExecuteByPageCursor called with cursor='%s'", nextPageCursor)

	// Obtener datos de la API usando el cursor
	apiResponse, err := uc.karenAIService.GetChallenges(nextPageCursor)
	if err != nil {
		log.Printf("API call failed: %v", err)
		return &dto.LoginResponse{
			Challenges: []dto.ChallengeDTO{},
			NextPage:   "",
			TotalPages: 1,
		}, nil
	}

	log.Printf("API call succeeded: got %d items, next_page='%s'", len(apiResponse.Items), apiResponse.NextPage)

	// Convertir respuesta de API a DTOs
	challenges := apiResponse.Items
	nextPage := apiResponse.NextPage

	// GUARDAR: Guardar challenges en BD para caché
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
	log.Printf("SaveChallenges result: %v", saveErr)

	// GUARDAR: Guardar el next_page en page_cursors
	if nextPage != "" {
		cursorErr := uc.challengeRepo.SaveCursorForDate(nextPage, time.Now().Format("2006-01-02"))
		log.Printf("SaveCursorForDate result: %v", cursorErr)
	}

	// Obtener total de páginas
	totalPages, err := uc.challengeRepo.GetTotalPages()
	if err != nil {
		totalPages = 1
	}

	log.Printf("Returning %d challenges with next_page='%s'", len(challenges), nextPage)

	return &dto.LoginResponse{
		Challenges: challenges,
		NextPage:   nextPage,
		TotalPages: totalPages,
	}, nil
}
