package service

import (
	"context"
	"fmt"

	"github.com/devprep/backend/internal/model"
	"github.com/devprep/backend/internal/repository"
)

type TopicService struct {
	repo  repository.TopicRepository
	qrepo repository.QuestionRepository
}

func NewTopicService(repo repository.TopicRepository, qrepo repository.QuestionRepository) *TopicService {
	return &TopicService{repo: repo, qrepo: qrepo}
}

func (s *TopicService) ListTopics(ctx context.Context) ([]model.Topic, error) {
	topics, err := s.repo.ListTopics(ctx)
	if err != nil {
		return nil, fmt.Errorf("topic service list: %w", err)
	}
	return topics, nil
}

func (s *TopicService) GetTopicWithQuestions(ctx context.Context, slug string, level model.Level) (*model.TopicWithQuestions, error) {
	topic, err := s.repo.GetTopicBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("topic service get: %w", err)
	}
	if topic == nil {
		return nil, nil
	}

	questions, err := s.qrepo.ListQuestionsByTopicSlug(ctx, slug, level)
	if err != nil {
		return nil, fmt.Errorf("topic service get questions: %w", err)
	}

	return &model.TopicWithQuestions{
		Topic:     *topic,
		Questions: questions,
	}, nil
}
