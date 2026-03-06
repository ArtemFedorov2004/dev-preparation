package handler

import (
	"log/slog"
	"net/http"

	"github.com/devprep/backend/internal/apperror"
	"github.com/devprep/backend/internal/model"
	"github.com/devprep/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type TopicHandler struct {
	svc *service.TopicService
}

func NewTopicHandler(svc *service.TopicService) *TopicHandler {
	return &TopicHandler{svc: svc}
}

// ListTopics godoc
//
//	@Summary		Список тем
//	@Description	Возвращает все темы, отсортированные по sort_order
//	@Tags			topics
//	@Produce		json
//	@Success		200	{array}		model.Topic
//	@Router			/topics [get]
func (h *TopicHandler) ListTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := h.svc.ListTopics(r.Context())
	if err != nil {
		slog.Error("list topics", "error", err)
		WriteError(w, r, err)
		return
	}
	if topics == nil {
		writeJSON(w, http.StatusOK, []struct{}{})
		return
	}
	writeJSON(w, http.StatusOK, topics)
}

// GetTopicBySlug godoc
//
//	@Summary		Тема со списком вопросов
//	@Description	Возвращает тему и связанные вопросы, опционально фильтруя по уровню
//	@Tags			topics
//	@Produce		json
//	@Param			slug	path		string		true	"Slug темы"				example(golang)
//	@Param			level	query		model.Level	false	"Уровень сложности"
//	@Success		200		{object}	model.TopicWithQuestions
//	@Failure		404		{object}	errorResponse
//	@Router			/topics/{slug} [get]
func (h *TopicHandler) GetTopicBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	rawLevel := r.URL.Query().Get("level")
	var level model.Level
	if model.Level(rawLevel).IsValid() {
		level = model.Level(rawLevel)
	}

	result, err := h.svc.GetTopicWithQuestions(r.Context(), slug, level)
	if err != nil {
		slog.Error("get topic by slug", "slug", slug, "error", err)
		WriteError(w, r, err)
		return
	}
	if result == nil {
		WriteError(w, r, apperror.NotFound("topic"))
		return
	}
	writeJSON(w, http.StatusOK, result)
}
