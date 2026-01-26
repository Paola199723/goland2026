package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := "postgresql://root@localhost:26257/market_data?sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	// Crear tablas si no existen
	createTables()

	log.Println("✅ CockroachDB conectado")
}

func GetConnection() *sql.DB {
	return DB
}

func createTables() {
	// Crear tabla de challenges
	challengesTable := `
	CREATE TABLE IF NOT EXISTS challenges (
		id SERIAL PRIMARY KEY,
		ticker VARCHAR(10) NOT NULL,
		target_from VARCHAR(20),
		target_to VARCHAR(20),
		company VARCHAR(100) NOT NULL,
		action TEXT,
		brokerage VARCHAR(100),
		rating_from VARCHAR(50),
		rating_to VARCHAR(50),
		time TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(ticker, DATE(time))
	);
	`

	// Crear tabla de cursores de paginación
	cursorTable := `
	CREATE TABLE IF NOT EXISTS page_cursors (
		id SERIAL PRIMARY KEY,
		cursor VARCHAR(255) UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	// Crear índices
	indexTicker := `CREATE INDEX IF NOT EXISTS idx_challenges_ticker ON challenges(ticker);`
	indexTime := `CREATE INDEX IF NOT EXISTS idx_challenges_time ON challenges(time DESC);`
	indexDate := `CREATE INDEX IF NOT EXISTS idx_challenges_date ON challenges(DATE(time));`

	queries := []string{
		challengesTable,
		cursorTable,
		indexTicker,
		indexTime,
		indexDate,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			log.Printf("Error creating table/index: %v", err)
		}
	}
}
