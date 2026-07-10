package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mavtice/golang-event-echo-server/internal/core/domain"
	core_errors "github.com/mavtice/golang-event-echo-server/internal/core/errors"
	core_postgres_pool "github.com/mavtice/golang-event-echo-server/internal/core/repository/postgres/pool"
)

// UpdateUser обновляет пользователя в БД с оптимистичной блокировкой.
// WHERE id=$7 AND version=$8 — условие гарантирует, что запись не была
// параллельно изменена. Если version не совпадает — возвращаем ErrConflict.
func (r *UsersRepository) UpdateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE event_echo.users
	SET
		first_name=$1,
		last_name=$2,
		email=$3,
		password=$4,
		role=$5,
		notification_token=$6,
		version=version+1
	WHERE id=$7 AND version=$8
	RETURNING
		id, 
		version,
		first_name,
		last_name,
		date_of_birthday,
		email,
		password,
		role,
		notification_token;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Role,
		user.NotificationToken,
		user.ID,
		user.Version,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.DateOfBirthday,
		&userModel.Email,
		&userModel.Password,
		&userModel.Role,
		&userModel.NotificationToken,

	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%d' concurrently accessed: %w",
				user.ID,
				core_errors.ErrConflict,
			)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := modelToDomain(userModel)

	return userDomain, nil
}
