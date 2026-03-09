package repository

import (
	"context"

	"github.com/devprep/backend/internal/model"
)

type UserActivityRepository interface {
	UpsertProgress(ctx context.Context, userID string, questionID int64, status model.ProgressStatus) error
	GetQuestionProgress(ctx context.Context, userID string, questionID int64) (*model.ProgressStatus, error)
	ListProgress(ctx context.Context, userID string) ([]model.UserProgress, error)
	ListProgressByTopic(ctx context.Context, userID string) ([]model.TopicProgress, error)

	ToggleBookmark(ctx context.Context, userID string, questionID int64) (bookmarked bool, err error)
	IsBookmarked(ctx context.Context, userID string, questionID int64) (bool, error)
	ListBookmarks(ctx context.Context, userID string) ([]model.Bookmark, error)

	RecordView(ctx context.Context, userID string, questionID int64) error
	ListHistory(ctx context.Context, userID string) ([]model.ViewHistory, error)

	GetQuestionIDBySlug(ctx context.Context, slug string) (int64, error)
}
