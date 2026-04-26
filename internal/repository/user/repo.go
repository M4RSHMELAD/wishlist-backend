package user

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/qwersedzxc/wishlist-backend/internal/entity"
)

type Repository struct {
	db *pgxpool.Pool
	qb squirrel.StatementBuilderType
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

var userColumns = []string{
	"id", "email", "username", "password_hash", "provider", "provider_id",
	"avatar_url", "full_name", "birth_date", "bio", "phone", "city", "notification_email",
	"push_notifications", "email_notifications", "show_fulfilled_wishes",
	"created_at", "updated_at",
}

func scanUser(row pgx.Row, user *entity.User) error {
	return row.Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.Provider, &user.ProviderID, &user.AvatarURL,
		&user.FullName, &user.BirthDate, &user.Bio, &user.Phone, &user.City, &user.NotificationEmail,
		&user.PushNotifications, &user.EmailNotifications, &user.ShowFulfilledWishes,
		&user.CreatedAt, &user.UpdatedAt,
	)
}

func scanUserFromRows(rows pgx.Rows, user *entity.User) error {
	return rows.Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.Provider, &user.ProviderID, &user.AvatarURL,
		&user.FullName, &user.BirthDate, &user.Bio, &user.Phone, &user.City, &user.NotificationEmail,
		&user.PushNotifications, &user.EmailNotifications, &user.ShowFulfilledWishes,
		&user.CreatedAt, &user.UpdatedAt,
	)
}

func (r *Repository) Create(ctx context.Context, user *entity.User) error {
	// Устанавливаем значения по умолчанию для настроек если они не заданы
	if !user.PushNotifications && !user.EmailNotifications && !user.ShowFulfilledWishes {
		user.PushNotifications = true
		user.EmailNotifications = false
		user.ShowFulfilledWishes = true
	}
	
	query, args, err := r.qb.
		Insert("users").
		Columns("id", "email", "username", "password_hash", "provider", "provider_id", "avatar_url", "push_notifications", "email_notifications", "show_fulfilled_wishes").
		Values(user.ID, user.Email, user.Username, user.PasswordHash, user.Provider, user.ProviderID, user.AvatarURL, user.PushNotifications, user.EmailNotifications, user.ShowFulfilledWishes).
		Suffix("RETURNING created_at, updated_at").
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	return r.db.QueryRow(ctx, query, args...).Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query, args, err := r.qb.
		Select(userColumns...).From("users").
		Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	user := &entity.User{}
	if err := scanUser(r.db.QueryRow(ctx, query, args...), user); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("execute query: %w", err)
	}
	return user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query, args, err := r.qb.
		Select(userColumns...).From("users").
		Where(squirrel.Eq{"email": email}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	user := &entity.User{}
	if err := scanUser(r.db.QueryRow(ctx, query, args...), user); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("execute query: %w", err)
	}
	return user, nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query, args, err := r.qb.
		Select(userColumns...).From("users").
		Where(squirrel.Eq{"username": username}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	user := &entity.User{}
	if err := scanUser(r.db.QueryRow(ctx, query, args...), user); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("execute query: %w", err)
	}
	return user, nil
}

func (r *Repository) GetByProviderID(ctx context.Context, provider, providerID string) (*entity.User, error) {
	query, args, err := r.qb.
		Select(userColumns...).From("users").
		Where(squirrel.Eq{"provider": provider, "provider_id": providerID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	user := &entity.User{}
	if err := scanUser(r.db.QueryRow(ctx, query, args...), user); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("execute query: %w", err)
	}
	return user, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, id uuid.UUID, fullName, bio, phone, city *string, birthDate *string) (*entity.User, error) {
	qb := r.qb.Update("users").Where(squirrel.Eq{"id": id})

	if fullName != nil {
		qb = qb.Set("full_name", *fullName)
	}
	if bio != nil {
		qb = qb.Set("bio", *bio)
	}
	if phone != nil {
		qb = qb.Set("phone", *phone)
	}
	if city != nil {
		qb = qb.Set("city", *city)
	}
	if birthDate != nil {
		qb = qb.Set("birth_date", *birthDate)
	}
	qb = qb.Set("updated_at", squirrel.Expr("now()"))

	query, args, err := qb.Suffix("RETURNING " + joinColumns(userColumns)).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	user := &entity.User{}
	if err := scanUser(r.db.QueryRow(ctx, query, args...), user); err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	return user, nil
}

func (r *Repository) SearchUsers(ctx context.Context, query string, limit int) ([]*entity.User, error) {
	sqlQuery, args, err := r.qb.
		Select("id", "email", "username", "avatar_url", "full_name", "created_at", "updated_at").
		From("users").
		Where(squirrel.Or{
			squirrel.ILike{"email": "%" + query + "%"},
			squirrel.ILike{"username": "%" + query + "%"},
		}).
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.db.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.AvatarURL, &user.FullName, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *Repository) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	query, args, err := r.qb.
		Select(userColumns...).
		From("users").
		Where("birth_date IS NOT NULL"). // Только пользователи с указанной датой рождения
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		if err := scanUserFromRows(rows, user); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func joinColumns(cols []string) string {
	result := ""
	for i, c := range cols {
		if i > 0 {
			result += ", "
		}
		result += c
	}
	return result
}

// UpdateSettings обновляет настройки пользователя
func (r *Repository) UpdateSettings(ctx context.Context, id uuid.UUID, pushNotifications, emailNotifications, showFulfilledWishes *bool) error {
	qb := r.qb.Update("users").Where(squirrel.Eq{"id": id})

	if pushNotifications != nil {
		qb = qb.Set("push_notifications", *pushNotifications)
	}
	if emailNotifications != nil {
		qb = qb.Set("email_notifications", *emailNotifications)
	}
	if showFulfilledWishes != nil {
		qb = qb.Set("show_fulfilled_wishes", *showFulfilledWishes)
	}
	qb = qb.Set("updated_at", squirrel.Expr("now()"))

	query, args, err := qb.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}
	return nil
}
