package users_transport_http

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/mavtice/golang-event-echo-server/internal/core/domain"
	core_logger "github.com/mavtice/golang-event-echo-server/internal/core/logger"
	core_http_request "github.com/mavtice/golang-event-echo-server/internal/core/transport/http/request"
	core_http_response "github.com/mavtice/golang-event-echo-server/internal/core/transport/http/response"
	core_http_types "github.com/mavtice/golang-event-echo-server/internal/core/transport/http/types"
)

var requestValidator = validator.New()

type PatchUserRequest struct {
	FirstName         core_http_types.Nullable[string] `json:"first_name"         swaggertype:"string" example:"Ivan"`
	LastName          core_http_types.Nullable[string] `json:"last_name"          swaggertype:"string" example:"Ivanov"`
	Email             core_http_types.Nullable[string] `json:"email"              swaggertype:"string" example:"email@example.com"`
	Password          core_http_types.Nullable[string] `json:"password"           swaggertype:"string" example:"0a6baa6ef716d..."`
	Role              core_http_types.Nullable[string] `json:"role"               swaggertype:"string" example:"user"`
	NotificationToken core_http_types.Nullable[string] `json:"notification_token" swaggertype:"string" example:"9001bba7ef1bb..."`
}

func (r *PatchUserRequest) Validate() error {
	if r.FirstName.Set {
		if r.FirstName.Value == nil {
			return fmt.Errorf("`FirstName` can't be NULL")
		}

		firstNameLen := len([]rune(*r.FirstName.Value))
		if firstNameLen < 3 || firstNameLen > 50 {
			return fmt.Errorf("`FirstName` must be between 3 and 50 symbols")
		}
	}

	if r.LastName.Set {
		if r.LastName.Value == nil {
			return fmt.Errorf("`LastName` can't be NULL")
		}

		lastNameLen := len([]rune(*r.LastName.Value))
		if lastNameLen < 3 || lastNameLen > 50 {
			return fmt.Errorf("`LastName` must be between 3 and 50 symbols")
		}
	}

	if r.Email.Set {
		if r.Email.Value == nil {
			return fmt.Errorf("`Email` can't be NULL")
		}

		if err := requestValidator.Var(*r.Email.Value, "required,email"); err != nil {
			return fmt.Errorf("`Email` must be a valid email address")
		}
	}

	if r.Password.Set {
		if r.Password.Value == nil {
			return fmt.Errorf("`Password` can't be NULL")
		}
	}

	if r.Role.Set {
		if r.Role.Value == nil {
			return fmt.Errorf("`Role` can't be NULL")
		}

		roleLen := len([]rune(*r.Role.Value))
		if roleLen < 4 || roleLen > 5 {
			return fmt.Errorf("`Role` must be between 4 and 5 symbols")
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser     godoc
// @Summary      Изменение пользователя
// @Description  Изменение информации об уже существующем в системе пользователе
// @Description  ### Логика обновления полей (Three-state logic):
// @Description  1. **Поле не передано**: `notification_token` игнорируется, значение в БД не меняется
// @Description  2. **Явно передано значение**: `"notification_token": "9001bba7ef1bb..."` - устанавливает новый токен в БД
// @Description  3. **Передан null**: `"notification_token": null` - очищает поле в БД (set to NULL)
// @Description  Ограничения: `first_name`, `last_name`, "date_of_birthday", `email`, `password`, `role` не могут быть выставлены как null
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id      path string           true            "ID изменяемого пользователя" Format(int)
// @Param        request body PatchUserRequest true            "PatchUser тело запроса"
// @Success      200 {object} PatchUserResponse                "Успешно изменённый пользователь"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      404 {object} core_http_response.ErrorResponse "User not found"
// @Failure      409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FirstName.ToDomain(),
		request.LastName.ToDomain(),
		request.Email.ToDomain(),
		request.Password.ToDomain(),
		request.Role.ToDomain(),
		request.NotificationToken.ToDomain(),
	)
}
