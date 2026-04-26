package notification

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qwersedzxc/wishlist-backend/internal/entity"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create создает новое уведомление
func (r *Repository) Create(ctx context.Context, notification entity.Notification) error {
	query := `
		INSERT INTO notifications (user_id, type, title, message, related_item_id, related_user_id, is_read)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.RelatedItemID,
		notification.RelatedUserID,
		notification.IsRead,
	)
	return err
}

// GetByUserID возвращает все уведомления пользователя
func (r *Repository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]entity.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, related_item_id, related_user_id, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	
	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []entity.Notification
	for rows.Next() {
		var n entity.Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Type,
			&n.Title,
			&n.Message,
			&n.RelatedItemID,
			&n.RelatedUserID,
			&n.IsRead,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, rows.Err()
}

// GetUnreadCount возвращает количество непрочитанных уведомлений
func (r *Repository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

// MarkAsRead помечает уведомление как прочитанное
func (r *Repository) MarkAsRead(ctx context.Context, notificationID uuid.UUID) error {
	query := `UPDATE notifications SET is_read = true WHERE id = $1`
	_, err := r.db.Exec(ctx, query, notificationID)
	return err
}

// MarkAllAsRead помечает все уведомления пользователя как прочитанные
func (r *Repository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE notifications SET is_read = true WHERE user_id = $1 AND is_read = false`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

// Delete удаляет уведомление
func (r *Repository) Delete(ctx context.Context, notificationID uuid.UUID) error {
	query := `DELETE FROM notifications WHERE id = $1`
	result, err := r.db.Exec(ctx, query, notificationID)
	if err != nil {
		return err
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("notification not found")
	}
	
	return nil
}
