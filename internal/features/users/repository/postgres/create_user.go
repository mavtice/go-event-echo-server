package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mavtice/golang-event-echo-server/internal/core/domain"
	core_errors "github.com/mavtice/golang-event-echo-server/internal/core/errors"
	core_postgres_pool "github.com/mavtice/golang-event-echo-server/internal/core/repository/postgres/pool"
)

// CreateUser вставляет нового пользователя в БД и возвращает сохранённую версию.
// RETURNING позволяет получить итоговое состояние записи одним запросом.
func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO event_echo.users (first_name, last_name, date_of_birthday, email, password, role, notification_token)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, version, first_name, last_name, date_of_birthday, email, password, role, notification_token;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.DateOfBirthday,
		user.Email,
		user.Password,
		user.Role,
		user.NotificationToken,
	)

	var userModel UserModel
	if err := userModel.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesUniqueKey) {
			return domain.User{}, fmt.Errorf(
				"create user with duplicated email: %w",
				core_errors.ErrConflict,
			)
		}

		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	userDomain := modelToDomain(userModel)

	return userDomain, nil
}
