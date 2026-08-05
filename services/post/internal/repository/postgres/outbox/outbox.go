package outbox

import (
	"context"
	"time"

	"github.com/disdreamq/fantastic-telegram/services/post/internal/domain"
	"github.com/jmoiron/sqlx"
)

type dbOutboxPost struct {
	id        int64     `db:"id"`
	payload   string    `db:"payload"`
	status    string    `db:"status"`
	createdAt time.Time `db:"created_at"`
	updatedAt time.Time `db:"updated_at"`
}

func (d *dbOutboxPost) toDomain() *domain.OutboxPost {
	return domain.NewOutboxPost(d.payload, d.status)
}

type OutboxRepository struct {
	db *sqlx.DB
}

func NewOutboxRepository(db *sqlx.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

func (r *OutboxRepository) GetPosts(ctx context.Context) ([]*domain.OutboxPost, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return []*domain.OutboxPost{}, err
	}
	defer tx.Rollback()
	txCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
		SELECT * FROM posts_outbox
		WHERE processed_at IS NULL
		ORDER BY created_at
		LIMIT 10
		FOR UPDATE SKIP LOCKED
	`
	var oPosts []*dbOutboxPost
	if err = r.db.SelectContext(txCtx, &oPosts, query); err != nil {
		return []*domain.OutboxPost{}, err
	}
	res := make([]*domain.OutboxPost, len(oPosts))
	for _, post := range oPosts {
		res = append(res, post.toDomain())
	}
	return res, nil
}

func (r *OutboxRepository) UpdatePosts(ctx context.Context, ids []int64) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	txCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
	UPDATE posts_outbox
	SET processed_at = NOW()
	WHERE id IN ($1)
`
	if _, err = tx.ExecContext(txCtx, query, ids); err != nil {
		return err
	}
	return nil
}
