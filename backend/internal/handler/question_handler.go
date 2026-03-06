package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/devprep/backend/internal/apperror"
	"github.com/devprep/backend/internal/model"
	"github.com/devprep/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type QuestionHandler struct {
	svc *service.QuestionService
}

func NewQuestionHandler(svc *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{svc: svc}
}

// ListQuestions godoc
//
//	@Summary		Список вопросов с пагинацией
//	@Description	Возвращает вопросы с фильтрацией по теме, тегу и уровню сложности
//	@Tags			questions
//	@Produce		json
//	@Param			topic	query		string		false	"Slug топика"
//	@Param			tag		query		string		false	"Slug тега"
//	@Param			level	query		model.Level	false	"Уровень сложности"
//	@Param			page	query		int			false	"Номер страницы"	minimum(1)	default(1)
//	@Param			limit	query		int			false	"Размер страницы"	minimum(1)	default(20)
//	@Success		200		{object}	model.PaginatedQuestions
//	@Router			/questions [get]
func (h *QuestionHandler) ListQuestions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	rawLevel := q.Get("level")
	var level model.Level
	if model.Level(rawLevel).IsValid() {
		level = model.Level(rawLevel)
	}

	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	filter := model.ListQuestionsFilter{
		TopicSlug: q.Get("topic"),
		TagSlug:   q.Get("tag"),
		Level:     level,
		Page:      page,
		Limit:     limit,
	}

	result, err := h.svc.ListQuestions(r.Context(), filter)
	if err != nil {
		slog.Error("list questions", "error", err)
		WriteError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// GetQuestionBySlug godoc
//
//	@Summary		Вопрос по slug
//	@Description	Возвращает полный вопрос с ответом, темой и тегами
//	@Tags			questions
//	@Produce		json
//	@Param			slug	path		string	true	"Slug вопроса"	example(what-is-goroutine)
//	@Success		200		{object}	model.QuestionDetail
//	@Failure		404		{object}	errorResponse
//	@Router			/questions/{slug} [get]
func (h *QuestionHandler) GetQuestionBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	question, err := h.svc.GetQuestionBySlug(r.Context(), slug)
	if err != nil {
		slog.Error("get question by slug", "slug", slug, "error", err)
		WriteError(w, r, err)
		return
	}
	if question == nil {
		WriteError(w, r, apperror.NotFound("question"))
		return
	}
	writeJSON(w, http.StatusOK, question)
}

// ListTags godoc
//
//	@Summary		Список тегов
//	@Description	Возвращает все доступные теги, отсортированные по имени
//	@Tags			tags
//	@Produce		json
//	@Success		200	{array}		model.Tag
//	@Router			/tags [get]
func (h *QuestionHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.svc.ListTags(r.Context())
	if err != nil {
		slog.Error("list tags", "error", err)
		WriteError(w, r, err)
		return
	}
	if tags == nil {
		writeJSON(w, http.StatusOK, []struct{}{})
		return
	}
	writeJSON(w, http.StatusOK, tags)
}
