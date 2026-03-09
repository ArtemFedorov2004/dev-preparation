package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/devprep/backend/internal/apperror"
	"github.com/devprep/backend/internal/model"
	"github.com/devprep/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type ctxKeyUserID struct{}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyUserID{}).(string)
	return v, ok && v != ""
}

func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ctxKeyUserID{}, userID)
}

type UserActivityHandler struct {
	svc *service.UserActivityService
}

func NewUserActivityHandler(svc *service.UserActivityService) *UserActivityHandler {
	return &UserActivityHandler{svc: svc}
}

type updateProgressRequest struct {
	Status string `json:"status"`
}

// UpdateProgress godoc
//
//	@Summary		Обновить статус изучения вопроса
//	@Description	Устанавливает статус изучения для вопроса
//	@Tags			progress
//	@Accept			json
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			slug	path		string						true	"Slug вопроса"
//	@Param			body	body		updateProgressRequest		true	"Статус"
//	@Success		204		"No Content"
//	@Failure		400		{object}	errorResponse
//	@Failure		401		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Router			/questions/{slug}/progress [post]
func (h *UserActivityHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	userID, _ := UserIDFromContext(r.Context())
	slug := chi.URLParam(r, "slug")
	var req updateProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, r, apperror.BadRequest("invalid request body"))
		return
	}
	if err := h.svc.UpdateProgress(r.Context(), userID, slug, model.ProgressStatus(req.Status)); err != nil {
		slog.Error("update progress", "slug", slug, "error", err)
		WriteError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetQuestionProgress godoc
//
//	@Summary		Прогресс по конкретному вопросу
//	@Description	Возвращает статус изучения вопроса текущим пользователем
//	@Tags			progress
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			slug	path		string				true	"Slug вопроса"
//	@Success		200		{object}	model.UserProgress
//	@Failure		401		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Router			/questions/{slug}/progress [get]
func (h *UserActivityHandler) GetQuestionProgress(w http.ResponseWriter, r *http.Request) {
	userID, _ := UserIDFromContext(r.Context())
	slug := chi.URLParam(r, "slug")
	p, err := h.svc.GetQuestionProgress(r.Context(), userID, slug)
	if err != nil {
		WriteError(w, r, err)
		return
	}
	if p == nil {
		WriteError(w, r, apperror.NotFound("progress"))
		return
	}
	writeJSON(w, http.StatusOK, p)

}

// GetMyProgress godoc
//
//	@Summary		Весь прогресс пользователя
//	@Description	Возвращает список всех вопросов со статусами изучения
//	@Tags			progress
//	@Produce		json
//	@Security		KeycloakAuth
//	@Success		200	{array}		model.UserProgress
//	@Failure		401	{object}	errorResponse
//	@Router			/me/progress [get]
func (h *UserActivityHandler) GetMyProgress(w http.ResponseWriter, r *http.Request) {
	userID, _ := UserIDFromContext(r.Context())
	list, err := h.svc.ListProgress(r.Context(), userID)
	if err != nil {
		slog.Error("get my progress", "error", err)
		WriteError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// GetMyProgressByTopic godoc
//
//	@Summary		Агрегированный прогресс по темам
//	@Description	Возвращает количество изученных, требующих повторения и незнакомых вопросов по каждой теме
//	@Tags			progress
//	@Produce		json
//	@Security		KeycloakAuth
//	@Success		200	{array}		model.TopicProgress
//	@Failure		401	{object}	errorResponse
//	@Router			/me/progress/by-topic [get]
func (h *UserActivityHandler) GetMyProgressByTopic(w http.ResponseWriter, r *http.Request) {
	userID, _ := UserIDFromContext(r.Context())
	list, err := h.svc.ListProgressByTopic(r.Context(), userID)
	if err != nil {
		slog.Error("get my progress by topic", "error", err)
		WriteError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

type bookmarkStatusResponse struct {
	Bookmarked bool `json:"bookmarked"`
}

// ToggleBookmark godoc
//
//	@Summary		Добавить / убрать закладку
//	@Description	Если закладки нет — создаёт, если есть — удаляет. Возвращает итоговый статус.
//	@Tags			bookmarks
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			slug	path		string					true	"Slug вопроса"
//	@Success		200		{object}	bookmarkStatusResponse
//	@Failure		401		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Router			/questions/{slug}/bookmark [post]
func (h *UserActivityHandler) ToggleBookmark(w http.ResponseWriter, r *http.Request) {
	userID, _ := UserIDFromContext(r.Context())
	slug := chi.URLParam(r, "slug")
	bookmarked, err := h.svc.ToggleBookmark(r.Context(), userID, slug)
	if err != nil {
		slog.Error("toggle bookmark", "slug", slug, "error", err)
		WriteError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, bookmarkStatusResponse{Bookmarked: bookmarked})
}

// GetBookmarkStatus godoc
//
//	@Summary		Статус закладки
//	@Description	Проверяет, добавлен ли вопрос в закладки текущим пользователем
//	@Tags			bookmarks
//	@Produce		json
//	@Security		KeycloakAuth
//	@Param			slug	path		string					true	"Slug вопроса"
//	@Success		200		{object}	bookmarkStatusResponse
//	@Failure		401		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Router			/questions/{slug}/bookmark [get]
func (h *UserActivityHandler) GetBookmarkStatus(w http.ResponseWriter, r *http.Request) {
	userID, _ := UserIDFromContext(r.Context())
	slug := chi.URLParam(r, "slug")
	bookmarked, err := h.svc.IsBookmarked(r.Context(), userID, slug)
	if err != nil {
		WriteError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"bookmarked": bookmarked})
}

// GetMyBookmarks godoc
//
//	@Summary		Список закладок
//	@Description	Возвращает все вопросы, добавленные пользователем в закладки
//	@Tags			bookmarks
//	@Produce		json
//	@Security		KeycloakAuth
//	@Success		200	{array}		model.Bookmark
//	@Failure		401	{object}	errorResponse
//	@Router			/me/bookmarks [get]
func (h *UserActivityHandler) GetMyBookmarks(w http.ResponseWriter, r *http.Request) {
	userID, _ := UserIDFromContext(r.Context())
	list, err := h.svc.ListBookmarks(r.Context(), userID)
	if err != nil {
		slog.Error("get my bookmarks", "error", err)
		WriteError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// RecordView godoc
//
//	@Summary		Записать просмотр вопроса
//	@Description	Фиксирует факт просмотра вопроса.
//	@Tags			history
//	@Security		KeycloakAuth
//	@Param			slug	path	string	true	"Slug вопроса"
//	@Success		204		"No Content"
//	@Failure		401		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Router			/questions/{slug}/view [post]
func (h *UserActivityHandler) RecordView(w http.ResponseWriter, r *http.Request) {
	userID, _ := UserIDFromContext(r.Context())
	slug := chi.URLParam(r, "slug")
	if err := h.svc.RecordView(r.Context(), userID, slug); err != nil {
		slog.Warn("record view", "slug", slug, "error", err)
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetMyHistory godoc
//
//	@Summary		История просмотров
//	@Description	Возвращает просмотренные вопросы, отсортированные по времени
//	@Tags			history
//	@Produce		json
//	@Security		KeycloakAuth
//	@Success		200	{array}		model.ViewHistory
//	@Failure		401	{object}	errorResponse
//	@Router			/me/history [get]
func (h *UserActivityHandler) GetMyHistory(w http.ResponseWriter, r *http.Request) {
	userID, _ := UserIDFromContext(r.Context())
	list, err := h.svc.ListHistory(r.Context(), userID)
	if err != nil {
		slog.Error("get my history", "error", err)
		WriteError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}
