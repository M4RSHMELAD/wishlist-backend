package notification

import (
	"context"

	"github.com/google/uuid"
	"github.com/qwersedzxc/wishlist-backend/internal/entity"
)

// Repository интерфейс репозитория уведомлений
type Repository interface {
	Create(ctx context.Context, notification entity.Notification) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]entity.Notification, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
	MarkAsRead(ctx context.Context, notificationID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, notificationID uuid.UUID) error
}
