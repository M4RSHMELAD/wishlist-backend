package usecase

import (
	"context"

	"github.com/qwersedzxc/wishlist-backend/internal/dto"
	"github.com/qwersedzxc/wishlist-backend/internal/entity"
	"github.com/google/uuid"
)

// WishlistUseCase описывает бизнес-операции над вишлистами
type WishlistUseCase interface {
	CreateWishlist(ctx context.Context, input dto.CreateWishlistInput) (entity.Wishlist, error)
	GetWishlist(ctx context.Context, id uuid.UUID) (entity.Wishlist, error)
	ListWishlists(ctx context.Context, filter dto.WishlistFilter) ([]entity.Wishlist, int, error)
	UpdateWishlist(ctx context.Context, id uuid.UUID, input dto.UpdateWishlistInput) (entity.Wishlist, error)
	DeleteWishlist(ctx context.Context, id uuid.UUID) error

	CreateItem(ctx context.Context, input dto.CreateWishlistItemInput) (entity.WishlistItem, error)
	GetItem(ctx context.Context, id uuid.UUID) (entity.WishlistItem, error)
	ListItems(ctx context.Context, filter dto.WishlistItemFilter) ([]entity.WishlistItem, int, error)
	UpdateItem(ctx context.Context, id uuid.UUID, input dto.UpdateWishlistItemInput) (entity.WishlistItem, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
	ReserveItem(ctx context.Context, itemID, userID uuid.UUID, isIncognito bool, hideReserverName bool) error
	UnreserveItem(ctx context.Context, itemID, userID uuid.UUID) error
}

// NotificationUseCase описывает бизнес-операции над уведомлениями
type NotificationUseCase interface {
	GetUserNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]entity.Notification, int, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
	MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
}
