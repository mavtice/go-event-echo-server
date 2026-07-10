package users_service

import (
	"context"
	"fmt"
)

// DeleteUser удаляет пользователя по ID, делегируя запрос репозиторию.
func (s *UsersService) DeleteUser(
	ctx context.Context,
	id int,
) error {
	if err := s.usersRepository.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
