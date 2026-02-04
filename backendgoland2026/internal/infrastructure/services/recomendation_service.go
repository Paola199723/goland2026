package services

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Paola199723/backendgoland2026/internal/infrastructure/db"
)

type StockRecommendation struct {
	Ticker     string   `json:"ticker"`
	Company    string   `json:"company"`
	Action     string   `json:"action"` // BUY, HOLD, WATCH, SELL
	Score      float64  `json:"score"`
	Confidence string   `json:"confidence"` // HIGH, MEDIUM, LOW
	Reasons    []string `json:"reasons"`
}

type StockRanking struct {
	Ticker string  `json:"ticker"`
	Score  float64 `json:"score"`
	Action string  `json:"action"`
}

type RecommendationResponse struct {
	Date             string              `json:"date"`
	RecommendedStock StockRecommendation `json:"recommended_stock"`
	Ranking          []StockRanking      `json:"ranking"`
	Disclaimer       string              `json:"disclaimer"`
}

// --- Utilidades ---
func parsePrice(s string) float64 {
	re := regexp.MustCompile(`[\\$]([0-9]+\\.?[0-9]*)`)
	match := re.FindStringSubmatch(s)
	if len(match) < 2 {
		return 0
	}
	val, _ := strconv.ParseFloat(match[1], 64)
	return val
}

func ratingScore(rating string) float64 {
	switch strings.ToLower(strings.TrimSpace(rating)) {
	case "strong buy":
		return 2
	case "buy":
		return 1
	case "overweight":
		return 1
	case "speculative buy":
		return 0.5
	case "neutral":
		return 0
	case "hold":
		return 0
	case "underweight":
		return -1
	case "sell":
		return -2
	default:
		return 0
	}
}

func actionScore(action string) float64 {
	switch strings.ToLower(strings.TrimSpace(action)) {
	case "target raised by":
		return 1
	case "target lowered by":
		return -1
	default:
		return 0
	}
}

func freshnessScore(signalTime time.Time, today time.Time) float64 {
	days := int(today.Sub(signalTime).Hours() / 24)
	switch {
	case days <= 3:
		return 1
	case days <= 7:
		return 0.5
	default:
		return 0
	}
}

func actionLabel(score float64) string {
	switch {
	case score >= 5:
		return "COMPRAR"
	case score >= 2:
		return "SOSTENER"
	case score >= 0:
		return "SOSTENER"
	default:
		return "VENDER"
	}
}

func confidenceLabel(score float64) string {
	switch {
	case score >= 5:
		return "HIGH"
	case score >= 2:
		return "MEDIUM"
	default:
		return "LOW"
	}
}

// --- Lógica principal ---
func RecommendBestActionOfDayV2(today time.Time) (*RecommendationResponse, error) {
	dbConn := db.GetConnection()

	// Buscar la fecha más reciente con datos
	var dateStr string
	row := dbConn.QueryRow(`SELECT DATE(time) FROM challenges WHERE DATE(time) <= $1 ORDER BY DATE(time) DESC LIMIT 1`, today.Format("2006-01-02"))
	if err := row.Scan(&dateStr); err != nil {
		return nil, fmt.Errorf("no hay acciones para hoy ni días anteriores")
	}

	// Obtener todas las acciones de esa fecha
	rows, err := dbConn.Query(`
        SELECT ticker, company, target_from, target_to, action, rating_to, time
        FROM challenges
        WHERE DATE(time) = $1
    `, dateStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type rowData struct {
		Ticker     string
		Company    string
		TargetFrom string
		TargetTo   string
		Action     string
		RatingTo   string
		Time       time.Time
	}

	var all []rowData
	for rows.Next() {
		var r rowData
		if err := rows.Scan(&r.Ticker, &r.Company, &r.TargetFrom, &r.TargetTo, &r.Action, &r.RatingTo, &r.Time); err == nil {
			all = append(all, r)
		}
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("no hay acciones para el día más reciente")
	}

	// Calcular score y razones para cada acción
	type scored struct {
		rowData
		Score      float64
		Reasons    []string
		ActionText string
		Confidence string
	}
	var scoredList []scored

	for _, a := range all {
		from := parsePrice(a.TargetFrom)
		to := parsePrice(a.TargetTo)
		expectedChange := 0.0
		if from > 0 {
			expectedChange = (to - from) / from
		}
		rScore := ratingScore(a.RatingTo)
		aScore := actionScore(a.Action)
		fScore := freshnessScore(a.Time, today)

		score := expectedChange*10 + rScore + aScore + fScore

		// Razones
		var reasons []string
		if aScore > 0 {
			reasons = append(reasons, "Precio objetivo planteado por las analistas")
		} else if aScore < 0 {
			reasons = append(reasons, "Precio objetivo bajado por las analistas")
		}
		if expectedChange != 0 {
			reasons = append(reasons, fmt.Sprintf("Se espera un aumento de %.1f%%", expectedChange*100))
		}
		if rScore > 0 {
			reasons = append(reasons, fmt.Sprintf("Calificación positiva de los analistas (%s)", a.RatingTo))
		} else if rScore < 0 {
			reasons = append(reasons, fmt.Sprintf("Calificación negativa de los analistas (%s)", a.RatingTo))
		} else {
			reasons = append(reasons, fmt.Sprintf("Calificación neutral de los analistas (%s)", a.RatingTo))
		}
		if fScore == 1 {
			reasons = append(reasons, "Actualización reciente de analistas")
		} else if fScore == 0.5 {
			reasons = append(reasons, "Actualización dentro de la última semana")
		} else {
			reasons = append(reasons, "Actualización antigua de analistas")
		}

		scoredList = append(scoredList, scored{
			rowData:    a,
			Score:      score,
			Reasons:    reasons,
			ActionText: actionLabel(score),
			Confidence: confidenceLabel(score),
		})
	}

	// Ordenar por score descendente
	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].Score > scoredList[j].Score
	})

	// Armar ranking
	var ranking []StockRanking
	for _, s := range scoredList {
		ranking = append(ranking, StockRanking{
			Ticker: s.Ticker,
			Score:  s.Score,
			Action: s.ActionText,
		})
	}

	// Mejor acción recomendada
	best := scoredList[0]
	resp := &RecommendationResponse{
		Date: dateStr,
		RecommendedStock: StockRecommendation{
			Ticker:     best.Ticker,
			Company:    best.Company,
			Action:     best.ActionText,
			Score:      best.Score,
			Confidence: best.Confidence,
			Reasons:    best.Reasons,
		},
		Ranking:    ranking,
		Disclaimer: "Esta recomendación se basa en objetivos de analistas y no garantiza el rendimiento futuro.",
	}
	return resp, nil
}
