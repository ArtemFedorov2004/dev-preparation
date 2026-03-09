package model

import "time"

type ProgressStatus string

const (
	StatusLearned    ProgressStatus = "learned"
	StatusNeedReview ProgressStatus = "need_review"
	StatusDontKnow   ProgressStatus = "dont_know"
)

func (s ProgressStatus) IsValid() bool {
	switch s {
	case StatusLearned, StatusNeedReview, StatusDontKnow:
		return true
	}
	return false
}

type UserProgress struct {
	QuestionID int64          `json:"question_id"`
	Slug       string         `json:"slug"`
	Title      string         `json:"title"`
	Level      Level          `json:"level"`
	Topic      Topic          `json:"topic"`
	Status     ProgressStatus `json:"status"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type TopicProgress struct {
	Topic      Topic `json:"topic"`
	Total      int   `json:"total"`
	Learned    int   `json:"learned"`
	NeedReview int   `json:"need_review"`
	DontKnow   int   `json:"dont_know"`
}

type Bookmark struct {
	QuestionID   int64     `json:"question_id"`
	Slug         string    `json:"slug"`
	Title        string    `json:"title"`
	Level        Level     `json:"level"`
	Topic        Topic     `json:"topic"`
	BookmarkedAt time.Time `json:"bookmarked_at"`
}

type ViewHistory struct {
	QuestionID int64     `json:"question_id"`
	Slug       string    `json:"slug"`
	Title      string    `json:"title"`
	Level      Level     `json:"level"`
	Topic      Topic     `json:"topic"`
	ViewedAt   time.Time `json:"viewed_at"`
}
