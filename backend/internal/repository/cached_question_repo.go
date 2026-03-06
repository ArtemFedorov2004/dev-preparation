package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/devprep/backend/internal/model"
	"github.com/redis/go-redis/v9"
)

const (
	keyQuestionSlug     = "questions:slug:%s"
	keyQuestionsByTopic = "questions:topic:%s:level:%s"
	keyTagList          = "tags:list"
)

type CachedQuestionRepo struct {
	repo QuestionRepository
	rdb  *redis.Client
	ttl  time.Duration
}

func NewCachedQuestionRepo(repo QuestionRepository, rdb *redis.Client, ttl time.Duration) *CachedQuestionRepo {
	return &CachedQuestionRepo{repo: repo, rdb: rdb, ttl: ttl}
}

func (r *CachedQuestionRepo) ListQuestions(ctx context.Context, filter model.ListQuestionsFilter) (*model.PaginatedQuestions, error) {
	return r.repo.ListQuestions(ctx, filter)
}

func (r *CachedQuestionRepo) GetQuestionBySlug(ctx context.Context, slug string) (*model.QuestionDetail, error) {
	key := fmt.Sprintf(keyQuestionSlug, slug)

	cached, err := r.rdb.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		slog.Warn("redis get failed", "key", key, "error", err)
	} else if err == nil {
		if string(cached) == "null" {
			return nil, nil
		}
		var q model.QuestionDetail
		if err := json.Unmarshal(cached, &q); err == nil {
			return &q, nil
		}
	}

	q, err := r.repo.GetQuestionBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(q); err == nil {
		if err := r.rdb.Set(ctx, key, data, r.ttl).Err(); err != nil {
			slog.Warn("redis set failed", "key", key, "error", err)
		}
	}

	return q, nil
}

func (r *CachedQuestionRepo) ListTags(ctx context.Context) ([]model.Tag, error) {
	cached, err := r.rdb.Get(ctx, keyTagList).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		slog.Warn("redis get failed", "key", keyTagList, "error", err)
	} else if err == nil {
		var tags []model.Tag
		if err := json.Unmarshal(cached, &tags); err == nil {
			return tags, nil
		}
	}
	tags, err := r.repo.ListTags(ctx)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(tags); err == nil {
		if err := r.rdb.Set(ctx, keyTagList, data, r.ttl).Err(); err != nil {
			slog.Warn("redis set failed", "key", keyTagList, "error", err)
		}
	}

	return tags, nil
}

func (r *CachedQuestionRepo) ListQuestionsByTopicSlug(ctx context.Context, topicSlug string, level model.Level) ([]model.QuestionListItem, error) {
	key := fmt.Sprintf(keyQuestionsByTopic, topicSlug, level)

	cached, err := r.rdb.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		slog.Warn("redis get failed", "key", key, "error", err)
	} else if err == nil {
		var questions []model.QuestionListItem
		if err := json.Unmarshal(cached, &questions); err == nil {
			return questions, nil
		}
	}

	questions, err := r.repo.ListQuestionsByTopicSlug(ctx, topicSlug, level)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(questions); err == nil {
		if err := r.rdb.Set(ctx, key, data, r.ttl).Err(); err != nil {
			slog.Warn("redis set failed", "key", key, "error", err)
		}
	}

	return questions, nil
}
