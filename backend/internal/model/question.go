package model

type Level string

const (
	LevelJunior Level = "junior"
	LevelMiddle Level = "middle"
	LevelSenior Level = "senior"
)

func (l Level) IsValid() bool {
	switch l {
	case LevelJunior, LevelMiddle, LevelSenior:
		return true
	}
	return false
}

type Question struct {
	ID      int    `json:"id"`
	TopicID *int   `json:"topic_id,omitempty"`
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Answer  string `json:"answer"`
	Level   Level  `json:"level"`
}

type QuestionDetail struct {
	Question
	Topic *Topic `json:"topic,omitempty"`
	Tags  []Tag  `json:"tags"`
}

type QuestionListItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Level Level  `json:"level"`
	Topic *Topic `json:"topic,omitempty"`
	Tags  []Tag  `json:"tags"`
}

type ListQuestionsFilter struct {
	TopicSlug string
	TagSlug   string
	Level     Level
	Page      int
	Limit     int
}

type PaginatedQuestions struct {
	Data       []QuestionListItem `json:"data"`
	Pagination Pagination         `json:"pagination"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}
