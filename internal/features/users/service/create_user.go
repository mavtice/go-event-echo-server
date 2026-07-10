package users_service

import (
	"context"
	"fmt"

	"github.com/mavtice/golang-event-echo-server/internal/core/domain"
)

// CreateUser создаёт нового пользователя: формирует доменный объект,
// валидирует его инварианты и сохраняет через репозиторий.
func (s *UsersService) CreateUser(
	ctx context.Context,
	fullName string,
	lastName string,
	dateOfBirthday string,
	email string,
	password string,
	role string,
	notificationToken *string,
) (domain.User, error) {
	user := domain.CreateUser(
		fullName,
		lastName,
		dateOfBirthday,
		email,
		password,
		role,
		notificationToken,
	)

	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	user, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("save user in repository: %w", err)
	}

	return user, nil
}
