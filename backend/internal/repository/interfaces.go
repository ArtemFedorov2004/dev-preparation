package repository

import (
	"context"

	"github.com/devprep/backend/internal/model"
)

type TopicRepository interface {
	ListTopics(ctx context.Context) ([]model.Topic, error)
	GetTopicBySlug(ctx context.Context, slug string) (*model.Topic, error)
}

type QuestionRepository interface {
	ListQuestions(ctx context.Context, filter model.ListQuestionsFilter) (*model.PaginatedQuestions, error)
	GetQuestionBySlug(ctx context.Context, slug string) (*model.QuestionDetail, error)
	ListTags(ctx context.Context) ([]model.Tag, error)
	ListQuestionsByTopicSlug(ctx context.Context, topicSlug string, level model.Level) ([]model.QuestionListItem, error)
}
