package sys

import (
	"net/http"
)

var (
	// ErrInvalidCredentials Ошибка при неверном логине или пароле
	ErrInvalidCredentials = "Неверный логин или пароль"

	// ErrUserNotFound Ошибка, когда пользователь не найден
	ErrUserNotFound = "Пользователь не найден"

	// ErrUserAlreadyExists Ошибка при попытке зарегистрировать уже существующего пользователя
	ErrUserAlreadyExists = "Пользователь с таким email уже существует"

	// ErrInvalidEmailFormat Ошибка при неверном формате email
	ErrInvalidEmailFormat = "Некорректный формат email"

	// ErrDatabaseFailure Ошибка при сбое в работе базы данных
	ErrDatabaseFailure = "Ошибка базы данных"

	// ErrInvalidToken Ошибка, когда токен недействителен
	ErrInvalidToken = "Неверный или истекший токен"

	// ErrAccessDenied Ошибка при попытке получить доступ без прав
	ErrAccessDenied = "Нет доступа"

	// ErrRateLimitExceeded Ошибка при превышении лимита запросов
	ErrRateLimitExceeded = "Превышен лимит запросов"

	// ErrInvalidRequestData Ошибка при неверных данных в запросе
	ErrInvalidRequestData = "Некорректные данные запроса"

	// ErrInsufficientPermissions Ошибка при отсутствии прав для выполнения операции
	ErrInsufficientPermissions = "Недостаточно прав для выполнения операции"

	// ErrEntityNotFound Ошибка, когда сущность не найдена
	ErrEntityNotFound = "Сущность не найдена"

	// ErrTokenGeneration Ошибка при сбое в процессе генерации токена
	ErrTokenGeneration = "Ошибка генерации токена"

	// ErrUnknown Неизвестная ошибка
	ErrUnknown = "Неизвестная ошибка"

	ErrPasswordHashGeneration = "Ошибка генерации хэша пароля"

	ErrCreateToken = "Ошибка генерации токена"

	ErrCreateUser = "Ошибка создания пользователя"

	ErrUpdateUser = "Ошибка обновления пользователя"

	ErrFilmNotFound = "Фильм не найден"

	ErrGRPCConnection = "Ошибка подключения к gRPC"

	ErrUserBlocked = "Пользователь заблокирован"

	ErrUserNotVerified = "Пользователь не верифицирован"

	ErrCollectionNotFound = "Коллекция не найдена"

	ErrFilmAlreadyAdded = "Фильм уже добавлен в коллекцию"

	ErrParentCommentNotFound = "Родительский комментарий не найден"

	ErrCommentNotFound = "Комментарий не найден"

	ErrFavouriteNotFound = "Избранное не найдено"

	ErrRecommendationsNotFound = "Рекомендации не найдены"
)

var ErrorMap = map[string]int{
	ErrInvalidCredentials:      http.StatusUnauthorized,
	ErrUserNotFound:            http.StatusNotFound,
	ErrUserAlreadyExists:       http.StatusConflict,
	ErrInvalidEmailFormat:      http.StatusBadRequest,
	ErrDatabaseFailure:         http.StatusInternalServerError,
	ErrInvalidToken:            http.StatusUnauthorized,
	ErrAccessDenied:            http.StatusForbidden,
	ErrRateLimitExceeded:       http.StatusTooManyRequests,
	ErrInvalidRequestData:      http.StatusBadRequest,
	ErrInsufficientPermissions: http.StatusForbidden,
	ErrEntityNotFound:          http.StatusNotFound,
	ErrTokenGeneration:         http.StatusInternalServerError,
	ErrUnknown:                 http.StatusInternalServerError,
	ErrPasswordHashGeneration:  http.StatusInternalServerError,
	ErrCreateToken:             http.StatusInternalServerError,
	ErrCreateUser:              http.StatusInternalServerError,
	ErrUpdateUser:              http.StatusInternalServerError,
	ErrFilmNotFound:            http.StatusNotFound,
	ErrGRPCConnection:          http.StatusInternalServerError,
	ErrUserBlocked:             http.StatusForbidden,
	ErrUserNotVerified:         http.StatusForbidden,
	ErrCollectionNotFound:      http.StatusNotFound,
	ErrFilmAlreadyAdded:        http.StatusConflict,
	ErrParentCommentNotFound:   http.StatusNotFound,
	ErrCommentNotFound:         http.StatusNotFound,
	ErrFavouriteNotFound:       http.StatusNotFound,
	ErrRecommendationsNotFound: http.StatusNotFound,
}
