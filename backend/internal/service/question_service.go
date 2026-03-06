package service

import (
	"context"
	"fmt"

	"github.com/devprep/backend/internal/model"
	"github.com/devprep/backend/internal/repository"
)

type QuestionService struct {
	repo repository.QuestionRepository
}

func NewQuestionService(repo repository.QuestionRepository) *QuestionService {
	return &QuestionService{repo: repo}
}

func (s *QuestionService) ListQuestions(ctx context.Context, f model.ListQuestionsFilter) (*model.PaginatedQuestions, error) {
	result, err := s.repo.ListQuestions(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("question service list: %w", err)
	}
	return result, nil
}

func (s *QuestionService) GetQuestionBySlug(ctx context.Context, slug string) (*model.QuestionDetail, error) {
	q, err := s.repo.GetQuestionBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("question service get: %w", err)
	}
	return q, nil
}

func (s *QuestionService) ListTags(ctx context.Context) ([]model.Tag, error) {
	tags, err := s.repo.ListTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("question service list tags: %w", err)
	}
	return tags, nil
}
