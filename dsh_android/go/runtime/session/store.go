package session

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Event struct {
	ID        int64  `json:"id"`
	SessionID string `json:"session_id"`
	Type      string `json:"type"`
	Payload   string `json:"payload"`
	Timestamp int64  `json:"timestamp"`
}

type Store interface {
	CreateSession() (SessionID, error)
	ListSessions() ([]Session, error)
	GetSession(id string) (*Session, error)
	DeleteSession(id string) error
	UpdateSessionTitle(id, title string) error
	ArchiveSession(id string) error
	AddEvent(event *Event) error
	GetEvents(sessionID string, afterID int64, limit int) ([]Event, error)
	Close() error
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSessionStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			title TEXT DEFAULT 'Untitled',
			created_at INTEGER,
			updated_at INTEGER,
			archived INTEGER DEFAULT 0
		)
	`); err != nil {
		return nil, err
	}
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id TEXT NOT NULL,
			type TEXT NOT NULL,
			payload TEXT,
			timestamp INTEGER
		)
	`); err != nil {
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) CreateSession() (SessionID, error) {
	id := fmt.Sprintf("sess_%d", time.Now().UnixNano())
	now := time.Now().Unix()
	_, err := s.db.Exec(
		"INSERT INTO sessions (id, title, created_at, updated_at) VALUES (?, ?, ?, ?)",
		id, "Untitled", now, now,
	)
	return SessionID(id), err
}

func (s *SQLiteStore) ListSessions() ([]Session, error) {
	rows, err := s.db.Query("SELECT id, title, created_at, updated_at, archived FROM sessions ORDER BY updated_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sessions []Session
	for rows.Next() {
		var sess Session
		var archived int
		if err := rows.Scan(&sess.ID, &sess.Title, &sess.CreatedAt, &sess.UpdatedAt, &archived); err != nil {
			return nil, err
		}
		sess.Archived = archived != 0
		sessions = append(sessions, sess)
	}
	return sessions, nil
}

func (s *SQLiteStore) GetSession(id string) (*Session, error) {
	var sess Session
	var archived int
	err := s.db.QueryRow(
		"SELECT id, title, created_at, updated_at, archived FROM sessions WHERE id = ?", id,
	).Scan(&sess.ID, &sess.Title, &sess.CreatedAt, &sess.UpdatedAt, &archived)
	if err != nil {
		return nil, err
	}
	sess.Archived = archived != 0
	return &sess, nil
}

func (s *SQLiteStore) DeleteSession(id string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE id = ?", id)
	return err
}

func (s *SQLiteStore) UpdateSessionTitle(id, title string) error {
	_, err := s.db.Exec("UPDATE sessions SET title = ?, updated_at = ? WHERE id = ?", title, time.Now().Unix(), id)
	return err
}

func (s *SQLiteStore) ArchiveSession(id string) error {
	_, err := s.db.Exec("UPDATE sessions SET archived = 1, updated_at = ? WHERE id = ?", time.Now().Unix(), id)
	return err
}

func (s *SQLiteStore) AddEvent(event *Event) error {
	event.Timestamp = time.Now().Unix()
	data, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(
		"INSERT INTO events (session_id, type, payload, timestamp) VALUES (?, ?, ?, ?)",
		event.SessionID, event.Type, string(data), event.Timestamp,
	)
	return err
}

func (s *SQLiteStore) GetEvents(sessionID string, afterID int64, limit int) ([]Event, error) {
	rows, err := s.db.Query(
		"SELECT id, session_id, type, payload, timestamp FROM events WHERE session_id = ? AND id > ? ORDER BY id ASC LIMIT ?",
		sessionID, afterID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event
	for rows.Next() {
		var ev Event
		var payloadBytes []byte
		if err := rows.Scan(&ev.ID, &ev.SessionID, &ev.Type, &payloadBytes, &ev.Timestamp); err != nil {
			return nil, err
		}
		ev.Payload = string(payloadBytes)
		events = append(events, ev)
	}
	return events, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
