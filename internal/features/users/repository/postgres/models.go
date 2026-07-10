package users_postgres_repository

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mavtice/golang-event-echo-server/internal/core/domain"
	core_postgres_pool "github.com/mavtice/golang-event-echo-server/internal/core/repository/postgres/pool"
)

// UserModel — структура для маппинга строки таблицы `event_echo.users` в Go-тип.
// Порядок полей совпадает с порядком столбцов в SELECT-запросах репозитория.
type UserModel struct {
	ID                int
	Version           int
	FirstName         string
	LastName          string
	DateOfBirthday    pgtype.Date
	Email             string
	Password          string
	Role              string
	NotificationToken *string
}

// Scan заполняет поля модели из результата запроса к БД.
func (m *UserModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.FirstName,
		&m.LastName,
		&m.DateOfBirthday,
		&m.Email,
		&m.Password,
		&m.Role,
		&m.NotificationToken,
	)
}

// modelToDomain конвертирует модель БД в доменный объект.
func modelToDomain(model UserModel) domain.User {
	dateOfBirthday := ""

	if model.DateOfBirthday.Valid {
		dateOfBirthday = model.DateOfBirthday.Time.Format("2006-01-02")
	}

	return domain.NewUser(
		model.ID,
		model.Version,
		model.FirstName,
		model.LastName,
		dateOfBirthday,
		model.Email,
		model.Password,
		model.Role,
		// TODO: почему здесь токен как двойной указатель???
		model.NotificationToken,
	)
}

// modelsToDomains конвертирует список моделей БД в список доменных объектов.
func modelsToDomains(models []UserModel) []domain.User {
	domains := make([]domain.User, len(models))

	for i, model := range models {
		domains[i] = modelToDomain(model)
	}

	return domains
}
