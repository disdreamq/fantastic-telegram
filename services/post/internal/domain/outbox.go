package domain

import "time"

type OutboxPost struct {
	ID        int64
	Payload   string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OutboxPayload struct {
	TraceID   string
	ID        int64
	UserID    int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOutboxPayload(post *Post, traceID string) *OutboxPayload {
	return &OutboxPayload{
		TraceID:   traceID,
		ID:        post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
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
