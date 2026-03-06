package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/devprep/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGQuestionRepo struct {
	db *pgxpool.Pool
}

func NewPGQuestionRepo(db *pgxpool.Pool) *PGQuestionRepo {
	return &PGQuestionRepo{db: db}
}

func (r *PGQuestionRepo) ListQuestions(ctx context.Context, f model.ListQuestionsFilter) (*model.PaginatedQuestions, error) {
	var (
		where  []string
		args   []any
		argIdx = 1
	)

	if f.Level != "" {
		where = append(where, fmt.Sprintf("q.level = $%d", argIdx))
		args = append(args, f.Level)
		argIdx++
	}
	if f.TopicSlug != "" {
		where = append(where, fmt.Sprintf("t.slug = $%d", argIdx))
		args = append(args, f.TopicSlug)
		argIdx++
	}
	if f.TagSlug != "" {
		where = append(where, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM question_tags qt2 JOIN tags tg2 ON tg2.id = qt2.tag_id WHERE qt2.question_id = q.id AND tg2.slug = $%d)",
			argIdx,
		))
		args = append(args, f.TagSlug)
		argIdx++
	}

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = "WHERE " + strings.Join(where, " AND ")
	}

	countSQL := fmt.Sprintf(`
		SELECT COUNT(DISTINCT q.id)
		FROM questions q
		LEFT JOIN topics t ON t.id = q.topic_id
		%s`, whereSQL)

	var total int
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("pg list questions count: %w", err)
	}

	page := f.Page
	limit := f.Limit
	offset := (page - 1) * limit

	dataSQL := fmt.Sprintf(`
		SELECT
			q.id, q.title, q.slug, q.level,
			t.id,   t.slug,   t.name, t.description, t.icon, t.sort_order
		FROM questions q
		LEFT JOIN topics t ON t.id = q.topic_id
		%s
		ORDER BY q.id DESC
		LIMIT $%d OFFSET $%d`, whereSQL, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, dataSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("pg list questions: %w", err)
	}
	defer rows.Close()

	questionMap := make(map[int]*model.QuestionListItem)
	var orderedIDs []int

	for rows.Next() {
		var item model.QuestionListItem
		var topic model.Topic
		var topicID *int

		if err := rows.Scan(
			&item.ID, &item.Title, &item.Slug, &item.Level,
			&topicID, &topic.Slug, &topic.Name, &topic.Description, &topic.Icon, &topic.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("pg list questions scan: %w", err)
		}
		if topicID != nil {
			topic.ID = *topicID
			item.Topic = &topic
		}
		questionMap[item.ID] = &item
		orderedIDs = append(orderedIDs, item.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := r.attachTags(ctx, questionMap, orderedIDs); err != nil {
		return nil, err
	}

	data := make([]model.QuestionListItem, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		data = append(data, *questionMap[id])
	}

	return &model.PaginatedQuestions{
		Data: data,
		Pagination: model.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func (r *PGQuestionRepo) GetQuestionBySlug(ctx context.Context, slug string) (*model.QuestionDetail, error) {
	const q = `
		SELECT
			q.id, q.topic_id, q.title, q.slug, q.answer, q.level,
			t.id,   t.slug,   t.name, t.description, t.icon, t.sort_order
		FROM questions q
		LEFT JOIN topics t ON t.id = q.topic_id
		WHERE q.slug = $1`

	var detail model.QuestionDetail
	var topic model.Topic
	var topicID *int

	err := r.db.QueryRow(ctx, q, slug).Scan(
		&detail.ID, &topicID, &detail.Title, &detail.Slug, &detail.Answer, &detail.Level,
		&topicID, &topic.Slug, &topic.Name, &topic.Description, &topic.Icon, &topic.SortOrder,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("pg get question by slug: %w", err)
	}
	if topicID != nil {
		topic.ID = *topicID
		detail.Topic = &topic
		detail.TopicID = topicID
	}

	detail.Tags, err = r.tagsByQuestionID(ctx, detail.ID)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r *PGQuestionRepo) ListTags(ctx context.Context) ([]model.Tag, error) {
	const q = `SELECT id, slug, name FROM tags ORDER BY name ASC`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("pg list tags: %w", err)
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.ID, &tag.Slug, &tag.Name); err != nil {
			return nil, fmt.Errorf("pg list tags scan: %w", err)
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

func (r *PGQuestionRepo) ListQuestionsByTopicSlug(ctx context.Context, topicSlug string, level model.Level) ([]model.QuestionListItem, error) {
	var (
		where  = []string{"t.slug = $1"}
		args   = []any{topicSlug}
		argIdx = 2
	)

	if level != "" {
		where = append(where, fmt.Sprintf("q.level = $%d", argIdx))
		args = append(args, level)
		argIdx++
	}

	sql := fmt.Sprintf(`
		SELECT
			q.id, q.title, q.slug, q.level,
			t.id, t.slug, t.name, t.description, t.icon, t.sort_order
		FROM questions q
		JOIN topics t ON t.id = q.topic_id
		WHERE %s
		ORDER BY q.id DESC`, strings.Join(where, " AND "))

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("pg list questions by topic: %w", err)
	}
	defer rows.Close()

	questionMap := make(map[int]*model.QuestionListItem)
	var orderedIDs []int

	for rows.Next() {
		var item model.QuestionListItem
		var topic model.Topic

		if err := rows.Scan(
			&item.ID, &item.Title, &item.Slug, &item.Level,
			&topic.ID, &topic.Slug, &topic.Name, &topic.Description, &topic.Icon, &topic.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("pg list questions by topic scan: %w", err)
		}
		item.Topic = &topic
		questionMap[item.ID] = &item
		orderedIDs = append(orderedIDs, item.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := r.attachTags(ctx, questionMap, orderedIDs); err != nil {
		return nil, err
	}

	result := make([]model.QuestionListItem, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		result = append(result, *questionMap[id])
	}
	return result, nil
}

func (r *PGQuestionRepo) tagsByQuestionID(ctx context.Context, questionID int) ([]model.Tag, error) {
	const q = `
		SELECT tg.id, tg.slug, tg.name
		FROM tags tg
		JOIN question_tags qt ON qt.tag_id = tg.id
		WHERE qt.question_id = $1
		ORDER BY tg.name`

	rows, err := r.db.Query(ctx, q, questionID)
	if err != nil {
		return nil, fmt.Errorf("pg tags by question: %w", err)
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.ID, &tag.Slug, &tag.Name); err != nil {
			return nil, fmt.Errorf("pg tags by question scan: %w", err)
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

func (r *PGQuestionRepo) attachTags(ctx context.Context, questionMap map[int]*model.QuestionListItem, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	sql := fmt.Sprintf(`
		SELECT qt.question_id, tg.id, tg.slug, tg.name
		FROM question_tags qt
		JOIN tags tg ON tg.id = qt.tag_id
		WHERE qt.question_id IN (%s)
		ORDER BY tg.name`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("pg attach tags: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var qID int
		var tag model.Tag
		if err := rows.Scan(&qID, &tag.ID, &tag.Slug, &tag.Name); err != nil {
			return fmt.Errorf("pg attach tags scan: %w", err)
		}
		if item, ok := questionMap[qID]; ok {
			item.Tags = append(item.Tags, tag)
		}
	}
	return rows.Err()
}
