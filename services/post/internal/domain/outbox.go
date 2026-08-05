package domain

import "time"

type OutboxPost struct {
	ID        int64
	Payload   string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOutboxPost(payload, status string) *OutboxPost {
	return &OutboxPost{
		ID:        0,
		Payload:   payload,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
