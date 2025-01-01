package sys

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/1Storm3/flibox-api/pkg/logger"
)

type Error struct {
	Message string // Сообщение для пользователя
	Err     string // Истинная ошибка
}

func (e *Error) Error() string {
	return fmt.Sprintf("message: %s, error: %v", e.Message, e.Err)
}

func NewError(message string, err string) *Error {
	return &Error{
		Message: message,
		Err:     err,
	}
}

func HandleError(c *fiber.Ctx, err error) error {
	logger.Error(err.Error())
	var appErr *Error
	if errors.As(err, &appErr) {
		if statusCode, exist := ErrorMap[appErr.Message]; exist {
			return c.Status(statusCode).JSON(fiber.Map{
				"message":    appErr.Message,
				"statusCode": statusCode,
			})
		}
	}
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"message":    "Внутренняя ошибка сервера",
		"statusCode": http.StatusInternalServerError,
	})
}
