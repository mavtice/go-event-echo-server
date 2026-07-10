package domain

import (
	"fmt"
	"github.com/go-playground/validator/v10"

	core_errors "github.com/mavtice/golang-event-echo-server/internal/core/errors"
)

var userValidator = validator.New()

// User — доменная сущность пользователя.
//
// NotificationToken — nil означает отсутствие токена (NULL в базе данных).
// Version — счётчик для оптимистичной блокировки: см. Task.Version.
type User struct {
	ID      int
	Version int

	FirstName         string
	LastName          string
	DateOfBirthday    string
	Email             string
	Password          string
	Role              string
	NotificationToken *string
}

// NewUser — конструктор для восстановления пользователя по имеющемуся набору данных
func NewUser(
	id int,
	version int,
	firstName string,
	lastName string,
	dateOfBirthday string,
	email string,
	password string,
	role string,
	notificationToken *string,
) User {
	return User{
		ID:                id,
		Version:           version,
		FirstName:         firstName,
		LastName:          lastName,
		DateOfBirthday:    dateOfBirthday,
		Email:             email,
		Password:          password,
		Role:              role,
		NotificationToken: notificationToken,
	}
}

// CreateUser создаёт нового пользователя без id и version (сгенерируются при вставке в БД)
func CreateUser(
	firstName string,
	lastName string,
	dateOfBirthday string,
	email string,
	password string,
	role string,
	notificationToken *string,
) User {
	var (
		id      = UninitializedID
		version = UninitializedVersion
	)

	return NewUser(
		id,
		version,
		firstName,
		lastName,
		dateOfBirthday,
		email,
		password,
		role,
		notificationToken,
	)
}

// Validate проверяет инварианты пользователя. TODO
func (u *User) Validate() error {
	firstNameLen := len([]rune(u.FirstName))
	if firstNameLen < 3 || firstNameLen > 50 {
		return fmt.Errorf(
			"invalid `FirstName` len: %d: %w",
			firstNameLen,
			core_errors.ErrInvalidArgument,
		)
	}

	lastNameLen := len([]rune(u.LastName))
	if lastNameLen < 3 || lastNameLen > 50 {
		return fmt.Errorf(
			"invalid `LastName` len: %d: %w",
			lastNameLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if err := userValidator.Var(u.Email, "required,email"); err != nil {
		return fmt.Errorf(
			"invalid `Email` format: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	passwordLen := len([]rune(u.Role))
	if passwordLen > 256 {
		return fmt.Errorf(
			"invalid `Password` len: %d: %w",
			passwordLen,
			core_errors.ErrInvalidArgument,
		)
	}

	roleLen := len([]rune(u.Role))
	if roleLen < 4 || roleLen > 5 {
		return fmt.Errorf(
			"invalid `Role` len: %d: %w",
			roleLen,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

// UserPatch содержит изменения для частичного обновления пользователя (PATCH).
// Каждое поле обёрнуто в Nullable, чтобы различать «не передано» и «передано null».
// Подробнее о Nullable: см. internal/core/domain/nullable.go.
type UserPatch struct {
	FirstName         Nullable[string]
	LastName          Nullable[string]
	Email             Nullable[string]
	Password          Nullable[string]
	Role              Nullable[string]
	NotificationToken Nullable[string]
}

// NewUserPatch — конструктор UserPatch.
func NewUserPatch(
	firstName         Nullable[string],
	lastName          Nullable[string],
	email             Nullable[string],
	password          Nullable[string],
	role              Nullable[string],
	notificationToken Nullable[string],
) UserPatch {
	return UserPatch{
		FirstName:    firstName,
		LastName: lastName,
		Email: email,
		Password: password,
		Role: role,
		NotificationToken: notificationToken,
	}
}

// Validate проверяет корректность патча до его применения.
func (p *UserPatch) Validate() error {
	if p.FirstName.Set && p.FirstName.Value == nil {
		return fmt.Errorf(
			"`FirstName` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.LastName.Set && p.LastName.Value == nil {
		return fmt.Errorf(
			"`LastName` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Email.Set && p.Email.Value == nil {
		return fmt.Errorf(
			"`Email` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Password.Set && p.Password.Value == nil {
		return fmt.Errorf(
			"`Password` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Role.Set && p.Role.Value == nil {
		return fmt.Errorf(
			"`Role` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

// ApplyPatch применяет изменения к пользователю.
// Используется техника «копия → изменение → валидация → замена».
func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.FirstName.Set {
		tmp.FirstName = *patch.FirstName.Value
	}

	if patch.LastName.Set {
		tmp.LastName = *patch.LastName.Value
	}
	
	if patch.Email.Set {
		tmp.Email = *patch.Email.Value
	}
	
	if patch.Password.Set {
		tmp.Password = *patch.Password.Value
	}

	if patch.Role.Set {
		tmp.Role = *patch.Role.Value
	}

	if patch.NotificationToken.Set {
		tmp.NotificationToken = patch.NotificationToken.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}
