package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/devprep/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGUserActivityRepo struct {
	db *pgxpool.Pool
}

func NewPGUserActivityRepo(db *pgxpool.Pool) *PGUserActivityRepo {
	return &PGUserActivityRepo{db: db}
}

func (r *PGUserActivityRepo) GetQuestionIDBySlug(ctx context.Context, slug string) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `SELECT id FROM questions WHERE slug = $1`, slug).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("get question id by slug: %w", err)
	}
	return id, nil
}

func (r *PGUserActivityRepo) UpsertProgress(ctx context.Context, userID string, questionID int64, status model.ProgressStatus) error {
	const q = `
		INSERT INTO user_progress (user_id, question_id, status, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, question_id)
		DO UPDATE SET status = EXCLUDED.status, updated_at = NOW()`

	_, err := r.db.Exec(ctx, q, userID, questionID, status)
	if err != nil {
		return fmt.Errorf("upsert progress: %w", err)
	}
	return nil
}

func (r *PGUserActivityRepo) GetQuestionProgress(ctx context.Context, userID string, questionID int64) (*model.ProgressStatus, error) {
	const q = `
		SELECT up.status
		FROM user_progress up
		WHERE up.user_id = $1 AND up.question_id = $2`

	var status model.ProgressStatus
	err := r.db.QueryRow(ctx, q, userID, questionID).Scan(
		&status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get progress: %w", err)
	}
	return &status, nil
}

func (r *PGUserActivityRepo) ListProgress(ctx context.Context, userID string) ([]model.UserProgress, error) {
	const q = `
		SELECT up.question_id, q.slug, q.title, q.level,
			t.id, t.slug, t.name, t.description, t.icon, t.sort_order,
			up.status, up.updated_at
		FROM user_progress up
		JOIN questions q ON q.id = up.question_id
		JOIN topics t ON t.id = q.topic_id
		WHERE up.user_id = $1
		ORDER BY up.updated_at DESC`

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list progress: %w", err)
	}
	defer rows.Close()

	var result []model.UserProgress
	for rows.Next() {
		var p model.UserProgress
		if err := rows.Scan(
			&p.QuestionID, &p.Slug, &p.Title, &p.Level,
			&p.Topic.ID, &p.Topic.Slug, &p.Topic.Name, &p.Topic.Description, &p.Topic.Icon, &p.Topic.SortOrder,
			&p.Status, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("list progress scan: %w", err)
		}
		result = append(result, p)
	}
	return result, rows.Err()
}

func (r *PGUserActivityRepo) ListProgressByTopic(ctx context.Context, userID string) ([]model.TopicProgress, error) {
	const q = `
		SELECT
			t.id, t.slug, t.name, t.description, t.icon, t.sort_order,
			COUNT(q.id) AS total,
			COUNT(up.question_id) FILTER (WHERE up.status = 'learned') AS learned,
			COUNT(up.question_id) FILTER (WHERE up.status = 'need_review') AS need_review,
			COUNT(up.question_id) FILTER (WHERE up.status = 'dont_know') AS dont_know
		FROM questions q
		JOIN topics t ON t.id = q.topic_id
		LEFT JOIN user_progress up ON up.question_id = q.id AND up.user_id = $1
		GROUP BY t.id, t.slug, t.name, t.description, t.icon, t.sort_order
		ORDER BY t.sort_order`

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list progress by topic: %w", err)
	}
	defer rows.Close()

	var result []model.TopicProgress
	for rows.Next() {
		var tp model.TopicProgress
		if err := rows.Scan(
			&tp.Topic.ID, &tp.Topic.Slug, &tp.Topic.Name, &tp.Topic.Description, &tp.Topic.Icon, &tp.Topic.SortOrder,
			&tp.Total, &tp.Learned, &tp.NeedReview, &tp.DontKnow,
		); err != nil {
			return nil, fmt.Errorf("list progress by topic scan: %w", err)
		}
		result = append(result, tp)
	}
	return result, rows.Err()
}

func (r *PGUserActivityRepo) ToggleBookmark(ctx context.Context, userID string, questionID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id=$1 AND question_id=$2)`,
		userID, questionID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("toggle bookmark check: %w", err)
	}

	if exists {
		_, err = r.db.Exec(ctx,
			`DELETE FROM bookmarks WHERE user_id=$1 AND question_id=$2`,
			userID, questionID,
		)
		if err != nil {
			return false, fmt.Errorf("toggle bookmark delete: %w", err)
		}
		return false, nil
	}

	_, err = r.db.Exec(ctx,
		`INSERT INTO bookmarks (user_id, question_id) VALUES ($1, $2)`,
		userID, questionID,
	)
	if err != nil {
		return false, fmt.Errorf("toggle bookmark insert: %w", err)
	}
	return true, nil
}

func (r *PGUserActivityRepo) IsBookmarked(ctx context.Context, userID string, questionID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id=$1 AND question_id=$2)`,
		userID, questionID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("is bookmarked: %w", err)
	}
	return exists, nil
}

func (r *PGUserActivityRepo) ListBookmarks(ctx context.Context, userID string) ([]model.Bookmark, error) {
	const q = `
		SELECT b.question_id, q.slug, q.title, q.level,
			   t.id, t.slug, t.name, t.description, t.icon, t.sort_order,
			   b.bookmarked_at
		FROM bookmarks b
		JOIN questions q ON q.id = b.question_id
		JOIN topics t ON t.id = q.topic_id
		WHERE b.user_id = $1
		ORDER BY b.bookmarked_at DESC`

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list bookmarks: %w", err)
	}
	defer rows.Close()

	var result []model.Bookmark
	for rows.Next() {
		var b model.Bookmark
		if err := rows.Scan(
			&b.QuestionID, &b.Slug, &b.Title, &b.Level,
			&b.Topic.ID, &b.Topic.Slug, &b.Topic.Name, &b.Topic.Description, &b.Topic.Icon, &b.Topic.SortOrder,
			&b.BookmarkedAt,
		); err != nil {
			return nil, fmt.Errorf("list bookmarks scan: %w", err)
		}
		result = append(result, b)
	}
	return result, rows.Err()
}

func (r *PGUserActivityRepo) RecordView(ctx context.Context, userID string, questionID int64) error {
	const q = `
		INSERT INTO view_history (user_id, question_id, viewed_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id, question_id)
		DO UPDATE SET viewed_at = NOW()`

	_, err := r.db.Exec(ctx, q, userID, questionID)
	if err != nil {
		return fmt.Errorf("record view: %w", err)
	}
	return nil
}

func (r *PGUserActivityRepo) ListHistory(ctx context.Context, userID string) ([]model.ViewHistory, error) {
	const q = `
		SELECT q.id, q.slug, q.title, q.level,
			   t.id, t.slug, t.name, t.description, t.icon, t.sort_order,
			   vh.viewed_at
		FROM view_history vh
		JOIN questions q ON q.id = vh.question_id
		JOIN topics t ON t.id = q.topic_id
		WHERE vh.user_id = $1
		ORDER BY vh.viewed_at DESC`

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list history: %w", err)
	}
	defer rows.Close()

	var result []model.ViewHistory
	for rows.Next() {
		var vh model.ViewHistory
		if err := rows.Scan(
			&vh.QuestionID, &vh.Slug, &vh.Title, &vh.Level,
			&vh.Topic.ID, &vh.Topic.Slug, &vh.Topic.Name, &vh.Topic.Description, &vh.Topic.Icon, &vh.Topic.SortOrder,
			&vh.ViewedAt,
		); err != nil {
			return nil, fmt.Errorf("list history scan: %w", err)
		}
		result = append(result, vh)
	}
	return result, rows.Err()
}
