package repositories

import (
	"github.com/Paola199723/backendgoland2026/internal/domain/entities"
)

type ChallengeRepository interface {
	SaveChallenges(challenges []entities.Challenge) error
	GetChallengesByDate(date string) ([]entities.Challenge, error)
	GetChallengesByPage(pageNum int) ([]entities.Challenge, error)
	SaveNextPageCursor(cursor string) error
	GetNextPageCursor(pageNum int) (string, error)
	GetTotalPages() (int, error)
}
