package notification

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/qwersedzxc/wishlist-backend/internal/entity"
)

type UseCase struct {
	repo Repository
	log  *slog.Logger
}

func New(repo Repository, log *slog.Logger) *UseCase {
	return &UseCase{
		repo: repo,
		log:  log,
	}
}

// GetUserNotifications возвращает уведомления пользователя
func (uc *UseCase) GetUserNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]entity.Notification, int, error) {
	notifications, err := uc.repo.GetByUserID(ctx, userID, limit)
	if err != nil {
		uc.log.ErrorContext(ctx, "failed to get notifications", "userID", userID, "error", err)
		return nil, 0, err
	}

	unreadCount, err := uc.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		uc.log.ErrorContext(ctx, "failed to get unread count", "userID", userID, "error", err)
		return notifications, 0, nil // Возвращаем уведомления даже если не удалось получить счетчик
	}

	return notifications, unreadCount, nil
}

// GetUnreadCount возвращает количество непрочитанных уведомлений
func (uc *UseCase) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := uc.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		uc.log.ErrorContext(ctx, "failed to get unread count", "userID", userID, "error", err)
		return 0, err
	}

	return count, nil
}

// MarkAsRead помечает уведомление как прочитанное
func (uc *UseCase) MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error {
	// TODO: Проверить что уведомление принадлежит пользователю
	if err := uc.repo.MarkAsRead(ctx, notificationID); err != nil {
		uc.log.ErrorContext(ctx, "failed to mark notification as read", "notificationID", notificationID, "error", err)
		return err
	}

	return nil
}

// MarkAllAsRead помечает все уведомления пользователя как прочитанные
func (uc *UseCase) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	if err := uc.repo.MarkAllAsRead(ctx, userID); err != nil {
		uc.log.ErrorContext(ctx, "failed to mark all notifications as read", "userID", userID, "error", err)
		return err
	}

	return nil
}
