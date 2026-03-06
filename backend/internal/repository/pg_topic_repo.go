package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/devprep/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGTopicRepo struct {
	db *pgxpool.Pool
}

func NewPGTopicRepo(db *pgxpool.Pool) *PGTopicRepo {
	return &PGTopicRepo{db: db}
}

func (r *PGTopicRepo) ListTopics(ctx context.Context) ([]model.Topic, error) {
	const q = `
		SELECT id, slug, name, description, icon, sort_order
		FROM topics
		ORDER BY sort_order ASC`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("pg list topics: %w", err)
	}
	defer rows.Close()

	var topics []model.Topic
	for rows.Next() {
		var t model.Topic
		if err := rows.Scan(&t.ID, &t.Slug, &t.Name, &t.Description, &t.Icon, &t.SortOrder); err != nil {
			return nil, fmt.Errorf("pg list topics scan: %w", err)
		}
		topics = append(topics, t)
	}
	return topics, rows.Err()
}

func (r *PGTopicRepo) GetTopicBySlug(ctx context.Context, slug string) (*model.Topic, error) {
	const q = `
		SELECT id, slug, name, description, icon, sort_order
		FROM topics
		WHERE slug = $1`

	var t model.Topic
	err := r.db.QueryRow(ctx, q, slug).Scan(&t.ID, &t.Slug, &t.Name, &t.Description, &t.Icon, &t.SortOrder)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("pg get topic by slug: %w", err)
	}
	return &t, nil
}
