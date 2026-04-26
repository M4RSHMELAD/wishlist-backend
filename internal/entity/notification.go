package entity

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType типы уведомлений
type NotificationType string

const (
	NotificationTypeGiftReserved     NotificationType = "gift_reserved"
	NotificationTypeBirthdayReminder NotificationType = "birthday_reminder"
)

// Notification представляет уведомление пользователя
type Notification struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Type          NotificationType
	Title         string
	Message       string
	RelatedItemID *uuid.UUID
	RelatedUserID *uuid.UUID
	IsRead        bool
	CreatedAt     time.Time
}
