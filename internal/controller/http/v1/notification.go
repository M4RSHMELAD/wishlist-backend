package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"github.com/qwersedzxc/wishlist-backend/internal/controller/http/v1/response"
	"github.com/qwersedzxc/wishlist-backend/internal/helpers"
	"github.com/qwersedzxc/wishlist-backend/internal/usecase"
)

type NotificationHandler struct {
	uc  usecase.NotificationUseCase
	log *slog.Logger
}

func newNotificationHandler(uc usecase.NotificationUseCase, log *slog.Logger) *NotificationHandler {
	return &NotificationHandler{
		uc:  uc,
		log: log,
	}
}

// ListNotifications возвращает список уведомлений текущего пользователя
func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserIDFromCtx(r.Context())
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.NewErrorResponse(errors.New("unauthorized")))
		return
	}

	notifications, unreadCount, err := h.uc.GetUserNotifications(r.Context(), userID, 50)
	if err != nil {
		h.log.Error("failed to get notifications", "error", err, "userID", userID)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.NewErrorResponse(err))
		return
	}

	render.JSON(w, r, response.NewNotificationListResponse(notifications, unreadCount))
}

// GetUnreadCount возвращает количество непрочитанных уведомлений
func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserIDFromCtx(r.Context())
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.NewErrorResponse(errors.New("unauthorized")))
		return
	}

	count, err := h.uc.GetUnreadCount(r.Context(), userID)
	if err != nil {
		h.log.Error("failed to get unread count", "error", err, "userID", userID)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.NewErrorResponse(err))
		return
	}

	render.JSON(w, r, map[string]int{"unreadCount": count})
}

// MarkAsRead помечает уведомление как прочитанное
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	notificationIDStr := chi.URLParam(r, "id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.NewErrorResponse(errors.New("invalid notification ID")))
		return
	}

	userID, err := helpers.GetUserIDFromCtx(r.Context())
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.NewErrorResponse(errors.New("unauthorized")))
		return
	}

	if err := h.uc.MarkAsRead(r.Context(), notificationID, userID); err != nil {
		h.log.Error("failed to mark notification as read", "error", err, "notificationID", notificationID)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.NewErrorResponse(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "notification marked as read"})
}

// MarkAllAsRead помечает все уведомления как прочитанные
func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserIDFromCtx(r.Context())
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.NewErrorResponse(errors.New("unauthorized")))
		return
	}

	if err := h.uc.MarkAllAsRead(r.Context(), userID); err != nil {
		h.log.Error("failed to mark all notifications as read", "error", err, "userID", userID)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.NewErrorResponse(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "all notifications marked as read"})
}
