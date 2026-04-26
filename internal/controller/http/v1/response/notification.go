package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/qwersedzxc/wishlist-backend/internal/entity"
)

// NotificationResponse ответ с данными уведомления
type NotificationResponse struct {
	ID            uuid.UUID  `json:"id"`
	Type          string     `json:"type"`
	Title         string     `json:"title"`
	Message       string     `json:"message"`
	RelatedItemID *uuid.UUID `json:"relatedItemId,omitempty"`
	RelatedUserID *uuid.UUID `json:"relatedUserId,omitempty"`
	IsRead        bool       `json:"isRead"`
	CreatedAt     time.Time  `json:"createdAt"`
}

func NewNotificationResponse(n entity.Notification) NotificationResponse {
	return NotificationResponse{
		ID:            n.ID,
		Type:          string(n.Type),
		Title:         n.Title,
		Message:       n.Message,
		RelatedItemID: n.RelatedItemID,
		RelatedUserID: n.RelatedUserID,
		IsRead:        n.IsRead,
		CreatedAt:     n.CreatedAt,
	}
}

// NotificationListResponse ответ со списком уведомлений
type NotificationListResponse struct {
	Items       []NotificationResponse `json:"items"`
	UnreadCount int                    `json:"unreadCount"`
}

func NewNotificationListResponse(notifications []entity.Notification, unreadCount int) NotificationListResponse {
	items := make([]NotificationResponse, len(notifications))
	for i, n := range notifications {
		items[i] = NewNotificationResponse(n)
	}

	return NotificationListResponse{
		Items:       items,
		UnreadCount: unreadCount,
	}
}
