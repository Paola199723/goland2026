package entities

import "time"

type Challenge struct {
	ID           int       `db:"id"`
	Ticker       string    `db:"ticker"`
	TargetFrom   string    `db:"target_from"`
	TargetTo     string    `db:"target_to"`
	Company      string    `db:"company"`
	Action       string    `db:"action"`
	Brokerage    string    `db:"brokerage"`
	RatingFrom   string    `db:"rating_from"`
	RatingTo     string    `db:"rating_to"`
	Time         time.Time `db:"time"`
	CreatedAt    time.Time `db:"created_at"`
	NextPageCursor string   `db:"next_page_cursor"`
}
