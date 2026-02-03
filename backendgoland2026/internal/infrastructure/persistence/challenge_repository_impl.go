package persistence

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Paola199723/backendgoland2026/internal/domain/entities"
	"github.com/Paola199723/backendgoland2026/internal/infrastructure/db"
)

type ChallengeRepositoryImpl struct {
	conn *sql.DB
}

func NewChallengeRepository() *ChallengeRepositoryImpl {
	return &ChallengeRepositoryImpl{
		conn: db.GetConnection(),
	}
}

func (r *ChallengeRepositoryImpl) SaveChallenges(challenges []entities.Challenge) error {
	if len(challenges) == 0 {
		return nil
	}

	for _, challenge := range challenges {
		// Primero verificar si ya existe un challenge con el mismo ticker
		checkQuery := `
			SELECT id FROM challenges 
			WHERE ticker = $1
			LIMIT 1
		`
		
		var existingID int64
		err := r.conn.QueryRow(checkQuery, challenge.Ticker).Scan(&existingID)
		
		// Si ya existe, no insertar (evitar duplicados)
		if err == nil {
			continue
		}
		
		// Si no existe, insertarlo
		insertQuery := `
			INSERT INTO challenges (ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`
		
		// Limpiar y normalizar strings para evitar caracteres especiales
		ticker := strings.TrimSpace(challenge.Ticker)
		targetFrom := strings.TrimSpace(challenge.TargetFrom)
		targetTo := strings.TrimSpace(challenge.TargetTo)
		company := strings.TrimSpace(challenge.Company)
		action := strings.TrimSpace(challenge.Action)
		brokerage := strings.TrimSpace(challenge.Brokerage)
		ratingFrom := strings.TrimSpace(challenge.RatingFrom)
		ratingTo := strings.TrimSpace(challenge.RatingTo)
		
		_, err = r.conn.Exec(
			insertQuery,
			ticker,
			targetFrom,
			targetTo,
			company,
			action,
			brokerage,
			ratingFrom,
			ratingTo,
			challenge.Time,
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("error saving challenge: %w", err)
		}
	}

	return nil
}

func (r *ChallengeRepositoryImpl) GetChallengesByDate(date string) ([]entities.Challenge, error) {
	query := `
		SELECT ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time
		FROM challenges
		WHERE DATE(time) = $1
		ORDER BY time DESC
	`

	rows, err := r.conn.Query(query, date)
	if err != nil {
		return nil, fmt.Errorf("error querying challenges: %w", err)
	}
	defer rows.Close()

	var challenges []entities.Challenge
	for rows.Next() {
		var c entities.Challenge
		err := rows.Scan(
			&c.Ticker,
			&c.TargetFrom,
			&c.TargetTo,
			&c.Company,
			&c.Action,
			&c.Brokerage,
			&c.RatingFrom,
			&c.RatingTo,
			&c.Time,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning challenge: %w", err)
		}
		challenges = append(challenges, c)
	}

	return challenges, nil
}

func (r *ChallengeRepositoryImpl) GetChallengesByPage(pageNum int) ([]entities.Challenge, error) {
	pageSize := 10
	offset := (pageNum - 1) * pageSize

	query := `
		SELECT ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time
		FROM challenges
		ORDER BY time DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.conn.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying challenges: %w", err)
	}
	defer rows.Close()

	var challenges []entities.Challenge
	for rows.Next() {
		var c entities.Challenge
		err := rows.Scan(
			&c.Ticker,
			&c.TargetFrom,
			&c.TargetTo,
			&c.Company,
			&c.Action,
			&c.Brokerage,
			&c.RatingFrom,
			&c.RatingTo,
			&c.Time,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning challenge: %w", err)
		}
		challenges = append(challenges, c)
	}

	return challenges, nil
}

func (r *ChallengeRepositoryImpl) SaveNextPageCursor(cursor string) error {
	query := `
		INSERT INTO page_cursors (cursor, created_at) VALUES ($1, $2)
		ON CONFLICT (cursor) DO NOTHING
	`
	_, err := r.conn.Exec(query, cursor, time.Now())
	return err
}

func (r *ChallengeRepositoryImpl) GetNextPageCursor(pageNum int) (string, error) {
	query := `
		SELECT cursor FROM page_cursors
		ORDER BY created_at DESC
		LIMIT 1 OFFSET $1
	`
	var cursor string
	err := r.conn.QueryRow(query, pageNum-1).Scan(&cursor)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return cursor, nil
}

func (r *ChallengeRepositoryImpl) GetTotalPages() (int, error) {
	query := `SELECT COUNT(*) / 10 + 1 FROM challenges`
	var total int
	err := r.conn.QueryRow(query).Scan(&total)
	return total, err
}

// GetChallengesFromToday retorna los desafíos guardados de hoy
func (r *ChallengeRepositoryImpl) GetChallengesFromToday() ([]entities.Challenge, error) {
	query := `
		SELECT ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time
		FROM challenges
		WHERE DATE(time) = CURRENT_DATE
		ORDER BY time DESC
	`

	rows, err := r.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying today's challenges: %w", err)
	}
	defer rows.Close()

	var challenges []entities.Challenge
	for rows.Next() {
		var c entities.Challenge
		err := rows.Scan(
			&c.Ticker,
			&c.TargetFrom,
			&c.TargetTo,
			&c.Company,
			&c.Action,
			&c.Brokerage,
			&c.RatingFrom,
			&c.RatingTo,
			&c.Time,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning challenge: %w", err)
		}
		challenges = append(challenges, c)
	}

	return challenges, nil
}

// SaveCursorForDate guarda el cursor de la siguiente página para una fecha específica
func (r *ChallengeRepositoryImpl) SaveCursorForDate(cursor string, date string) error {
	query := `
		INSERT INTO page_cursors (cursor, created_at) 
		VALUES ($1, $2)
		ON CONFLICT (cursor) DO NOTHING
	`
	_, err := r.conn.Exec(query, cursor, time.Now())
	return err
}

// GetCursorForToday retorna el cursor de la siguiente página más reciente
func (r *ChallengeRepositoryImpl) GetCursorForToday() (string, error) {
	query := `
		SELECT cursor FROM page_cursors
		ORDER BY created_at DESC
		LIMIT 1
	`
	var cursor string
	err := r.conn.QueryRow(query).Scan(&cursor)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("error getting cursor: %w", err)
	}
	return cursor, nil
}
