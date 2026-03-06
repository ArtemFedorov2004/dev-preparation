package model

type Topic struct {
	ID          int     `json:"id"`
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"`
	SortOrder   int     `json:"sort_order"`
}

type TopicWithQuestions struct {
	Topic
	Questions []QuestionListItem `json:"questions"`
}
