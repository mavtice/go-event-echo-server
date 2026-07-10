package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/mavtice/golang-event-echo-server/internal/core/errors"
)


// GetIntPathValue извлекает переменную пути по ключу key и парсит её как int.
func GetIntPathValue(r *http.Request, key string) (int, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf(
			"no key='%s' in path values: %w",
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	val, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf(
			"path value='%s' by key='%s' not a valid integer: %v: %w",
			pathValue,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return val, nil
}
