package service

import (
	"context"
	"fmt"

	"github.com/devprep/backend/internal/apperror"
	"github.com/devprep/backend/internal/model"
	"github.com/devprep/backend/internal/repository"
)

type UserActivityService struct {
	repo repository.UserActivityRepository
}

func NewUserActivityService(repo repository.UserActivityRepository) *UserActivityService {
	return &UserActivityService{repo: repo}
}

func (s *UserActivityService) resolveQuestionID(ctx context.Context, slug string) (int64, error) {
	id, err := s.repo.GetQuestionIDBySlug(ctx, slug)
	if err != nil {
		return 0, fmt.Errorf("resolve question id: %w", err)
	}
	if id == 0 {
		return 0, apperror.NotFound("question")
	}
	return id, nil
}

func (s *UserActivityService) UpdateProgress(ctx context.Context, userID, slug string, status model.ProgressStatus) error {
	if !status.IsValid() {
		return apperror.BadRequest("invalid status")
	}
	qID, err := s.resolveQuestionID(ctx, slug)
	if err != nil {
		return err
	}
	if err := s.repo.UpsertProgress(ctx, userID, qID, status); err != nil {
		return fmt.Errorf("update progress: %w", err)
	}
	return nil
}

func (s *UserActivityService) GetQuestionProgress(ctx context.Context, userID, slug string) (*model.ProgressStatus, error) {
	qID, err := s.resolveQuestionID(ctx, slug)
	if err != nil {
		return nil, err
	}
	p, err := s.repo.GetQuestionProgress(ctx, userID, qID)
	if err != nil {
		return nil, fmt.Errorf("get question progress: %w", err)
	}
	return p, nil
}

func (s *UserActivityService) ListProgress(ctx context.Context, userID string) ([]model.UserProgress, error) {
	list, err := s.repo.ListProgress(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list progress: %w", err)
	}
	if list == nil {
		return []model.UserProgress{}, nil
	}
	return list, nil
}

func (s *UserActivityService) ListProgressByTopic(ctx context.Context, userID string) ([]model.TopicProgress, error) {
	list, err := s.repo.ListProgressByTopic(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list progress by topic: %w", err)
	}
	if list == nil {
		return []model.TopicProgress{}, nil
	}
	return list, nil
}

func (s *UserActivityService) ToggleBookmark(ctx context.Context, userID, slug string) (bool, error) {
	qID, err := s.resolveQuestionID(ctx, slug)
	if err != nil {
		return false, err
	}
	bookmarked, err := s.repo.ToggleBookmark(ctx, userID, qID)
	if err != nil {
		return false, fmt.Errorf("toggle bookmark: %w", err)
	}
	return bookmarked, nil
}

func (s *UserActivityService) IsBookmarked(ctx context.Context, userID, slug string) (bool, error) {
	qID, err := s.resolveQuestionID(ctx, slug)
	if err != nil {
		return false, err
	}
	ok, err := s.repo.IsBookmarked(ctx, userID, qID)
	if err != nil {
		return false, fmt.Errorf("is bookmarked: %w", err)
	}
	return ok, nil
}

func (s *UserActivityService) ListBookmarks(ctx context.Context, userID string) ([]model.Bookmark, error) {
	list, err := s.repo.ListBookmarks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list bookmarks: %w", err)
	}
	if list == nil {
		return []model.Bookmark{}, nil
	}
	return list, nil
}

func (s *UserActivityService) RecordView(ctx context.Context, userID, slug string) error {
	qID, err := s.resolveQuestionID(ctx, slug)
	if err != nil {
		return err
	}
	if err := s.repo.RecordView(ctx, userID, qID); err != nil {
		return fmt.Errorf("record view: %w", err)
	}
	return nil
}

func (s *UserActivityService) ListHistory(ctx context.Context, userID string) ([]model.ViewHistory, error) {
	list, err := s.repo.ListHistory(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list history: %w", err)
	}
	if list == nil {
		return []model.ViewHistory{}, nil
	}
	return list, nil
}
