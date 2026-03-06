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
	keyTopicList = "topics:list"
	keyTopicSlug = "topics:slug:%s"
)

type CachedTopicRepo struct {
	repo TopicRepository
	rdb  *redis.Client
	ttl  time.Duration
}

func NewCachedTopicRepo(repo TopicRepository, rdb *redis.Client, ttl time.Duration) *CachedTopicRepo {
	return &CachedTopicRepo{repo: repo, rdb: rdb, ttl: ttl}
}

func (r *CachedTopicRepo) ListTopics(ctx context.Context) ([]model.Topic, error) {
	cached, err := r.rdb.Get(ctx, keyTopicList).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		slog.Warn("redis get failed", "key", keyTopicList, "error", err)
	} else if err == nil {
		var topics []model.Topic
		if err := json.Unmarshal(cached, &topics); err == nil {
			return topics, nil
		}
	}

	topics, err := r.repo.ListTopics(ctx)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(topics); err == nil {
		if err := r.rdb.Set(ctx, keyTopicList, data, r.ttl).Err(); err != nil {
			slog.Warn("redis set failed", "key", keyTopicList, "error", err)
		}
	}

	return topics, nil
}

func (r *CachedTopicRepo) GetTopicBySlug(ctx context.Context, slug string) (*model.Topic, error) {
	key := fmt.Sprintf(keyTopicSlug, slug)

	cached, err := r.rdb.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		slog.Warn("redis get failed", "key", key, "error", err)
	} else if err == nil {
		if string(cached) == "null" {
			return nil, nil
		}
		var topic model.Topic
		if err := json.Unmarshal(cached, &topic); err == nil {
			return &topic, nil
		}
	}

	topic, err := r.repo.GetTopicBySlug(ctx, slug)
	fmt.Print(topic)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(topic); err == nil {
		if err := r.rdb.Set(ctx, key, data, r.ttl).Err(); err != nil {
			slog.Warn("redis set failed", "key", key, "error", err)
		}
	}

	return topic, nil
}
