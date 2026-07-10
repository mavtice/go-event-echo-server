package users_transport_http

import (
	"github.com/mavtice/golang-event-echo-server/internal/core/domain"
)

// UserDTOResponse — DTO для представления пользователя в API-ответе.
type UserDTOResponse struct {
	ID                int     `json:"id"                   example:"1"`
	Version           int     `json:"version"              example:"3"`
	FirstName         string  `json:"first_name"           example:"Ivan"`
	LastName          string  `json:"last_name"            example:"Ivanov"`
	DateOfBirthday    string  `json:"date_of_birthday"     example:"2006-01-02"`
	Email             string  `json:"email"                example:"email@example.com"`
	Role              string  `json:"role"                 example:"user"`
	NotificationToken *string `json:"notification_token"   example:"9001bba7ef1bb..."`
}

// userDTOFromDomain конвертирует доменный объект User в DTO для HTTP-ответа.
func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:                user.ID,
		Version:           user.Version,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		DateOfBirthday:    user.DateOfBirthday,
		Email:             user.Email,
		Role:              user.Role,
		NotificationToken: user.NotificationToken,
	}
}

// usersDTOFromDomains конвертирует список доменных объектов в список DTO.
func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}
