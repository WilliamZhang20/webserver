package data

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

type PostgresPageStore struct {
	db *sql.DB
}

func NewPostgresPageStore(dataSourceName string) (*PostgresPageStore, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("Connected to PostgreSQL")
	return &PostgresPageStore{db: db}, nil
}

func (p *PostgresPageStore) SavePage(page *Page) error {
	_, err := p.db.Exec("INSERT INTO pages (title, body) VALUES ($1, $2) ON CONFLICT (title) DO UPDATE SET body = $2", page.Title, page.Body)
	if err != nil {
		return fmt.Errorf("failed to save page: %w", err)
	}
	return nil
}

func (p *PostgresPageStore) LoadPage(title string) (*Page, error) {
	row := p.db.QueryRow("SELECT title, body FROM pages WHERE title = $1", title)
	page := &Page{}
	err := row.Scan(&page.Title, &page.Body)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("page '%s' not found", title)
		}
		return nil, fmt.Errorf("failed to load page '%s': %w", title, err)
	}
	return page, nil
}
