package httperror

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Error представляет структуру ошибки
// @swagger:model
type Error struct {
	// Код ошибки
	// example: 400
	code int `json:"code"`

	// Сообщение ошибки
	// example: Invalid request parameters
	message string `json:"message"`
}

func New(code int, message string) error {
	return &Error{
		code:    code,
		message: message,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s", e.message)
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Is(target error) bool {
	var t *Error
	ok := errors.As(target, &t)

	if !ok {
		return false
	}

	return e.code == t.code
}

func HandleError(c *fiber.Ctx, err error) error {
	var httpErr *Error
	if errors.As(err, &httpErr) {
		return c.Status(httpErr.Code()).JSON(fiber.Map{
			"message":    httpErr.Error(),
			"statusCode": httpErr.Code(),
		})
	}

	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"message":    "Internal server error",
		"statusCode": http.StatusInternalServerError,
	})
}
