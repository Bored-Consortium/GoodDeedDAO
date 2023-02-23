package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"GoodDeedDAO/storage"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage.
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

// Save saves page to storage.
func (s *Storage) Save(ctx context.Context, p *storage.User) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

// PickRandom picks random page from storage.
func (s *Storage) AddKarma(ctx context.Context, userName string) (*storage.User, error) {
	q := `SELECT url FROM pages WHERE user_name = ?`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &storage.User{
		URL:      url,
		UserName: userName,
	}, nil
}

// Remove removes page from storage.
func (s *Storage) Remove(ctx context.Context, page *storage.User) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`
	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return nil
}

// IsExists checks if page exists in storage.
func (s *Storage) IsExists(ctx context.Context, page *storage.User) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q0 := `CREATE TABLE IF NOT EXISTS USERS (id_user 	INTEGER, 
											user_name	TEXT, 
											karma 		INTEGER, 
											deeds 		INTEGER, 
											validations INTEGER)`

	q1 := `CREATE TABLE IF NOT EXISTS DEEDS (id_deed 	 INTEGER, 
											upvote	 	 INTEGER, 
											downvote 	 INTEGER, 
											is_validated INTEGER,
											type 		 TEXT)`

	q2 := `CREATE TABLE IF NOT EXISTS DEED_BY_USER (id_deed 	INTEGER, 
													id_user 	INTEGER)`
	_, err := s.db.ExecContext(ctx, q0)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	_, err1 := s.db.ExecContext(ctx, q1)
	if err1 != nil {
		return fmt.Errorf("can't create table: %w", err1)
	}

	_, err2 := s.db.ExecContext(ctx, q2)
	if err2 != nil {
		return fmt.Errorf("can't create table: %w", err2)
	}
	return nil
}

// GetUserInfo return user's info
func (s *Storage) GetUserInfo(ctx context.Context, username string) (*storage.User, error) {
	q := `SELECT karma, deeds, validations FROM USERS WHERE user_name = ?`

	var k, d, v int

	err := s.db.QueryRowContext(ctx, q, username).Scan(&k, &d, &v)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick user's info: %w", err)
	}

	return &storage.User{
		UserName:    username,
		Karma:       k,
		Deeds:       d,
		Validations: v,
	}, nil
}
