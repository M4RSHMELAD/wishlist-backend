package wishlist

import (
	"context"

	"github.com/google/uuid"

	"github.com/qwersedzxc/wishlist-backend/internal/dto"
	"github.com/qwersedzxc/wishlist-backend/internal/entity"
)

// Repository интерфейс репозитория вишлистов
type Repository interface {
	Create(ctx context.Context, input dto.CreateWishlistInput) (entity.Wishlist, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.Wishlist, error)
	List(ctx context.Context, filter dto.WishlistFilter) ([]entity.Wishlist, int, error)
	Update(ctx context.Context, id uuid.UUID, input dto.UpdateWishlistInput) (entity.Wishlist, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ItemRepository интерфейс репозитория элементов вишлиста
type ItemRepository interface {
	Create(ctx context.Context, input dto.CreateWishlistItemInput) (entity.WishlistItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.WishlistItem, error)
	List(ctx context.Context, filter dto.WishlistItemFilter) ([]entity.WishlistItem, int, error)
	Update(ctx context.Context, id uuid.UUID, input dto.UpdateWishlistItemInput) (entity.WishlistItem, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Reserve(ctx context.Context, itemID, userID uuid.UUID, isIncognito bool) error
	Unreserve(ctx context.Context, itemID, userID uuid.UUID) error
}

// NotificationRepository интерфейс репозитория уведомлений
type NotificationRepository interface {
	Create(ctx context.Context, notification entity.Notification) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]entity.Notification, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
	MarkAsRead(ctx context.Context, notificationID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, notificationID uuid.UUID) error
}

// UserRepository интерфейс репозитория пользователей
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}
