package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/mavtice/golang-event-echo-server/internal/core/errors"
)

// DeleteUser удаляет пользователя по ID.
// Если RowsAffected() == 0 — пользователь не существовал → ErrNotFound.
func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM event_echo.users
	WHERE id=$1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%d': %w", id, core_errors.ErrNotFound)
	}

	return nil
}
