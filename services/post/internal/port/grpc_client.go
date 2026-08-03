package port

import (
	"context"

	"github.com/disdreamq/fantastic-telegram/services/post/internal/domain"
)

type GRPCClient interface {
	ValidateToken(ctx context.Context, token string) (*domain.Claims, error)
}
