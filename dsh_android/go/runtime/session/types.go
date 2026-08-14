package session

import "time"

type SessionID string

type Session struct {
	ID        SessionID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
	Archived  bool      `json:"archived"`
}

func (s *Session) UpdateTimestamp() {
	now := time.Now().Unix()
	s.UpdatedAt = now
	if s.CreatedAt == 0 {
		s.CreatedAt = now
	}
}
