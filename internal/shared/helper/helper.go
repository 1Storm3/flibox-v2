package helper

import (
	"bytes"
	_ "embed"
	"errors"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/1Storm3/flibox-api/pkg/logger"
	"github.com/1Storm3/flibox-api/pkg/sys"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"strings"
	"text/template"

	"go.uber.org/zap"
)

//go:embed template/email.html
var emailTemplate string

func ExtractS3Key(photoURL string) (string, error) {
	parsedURL, err := url.Parse(photoURL)
	if err != nil {
		return "", err
	}
	segments := strings.Split(parsedURL.Path, "/")
	return segments[len(segments)-1], nil
}

func TakeHTMLTemplate(appUrl, verificationToken string) (string, error) {
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		logger.Info("Ошибка при создании шаблона: %v", zap.Error(err))
		return "", err
	}

	data := struct {
		AppUrl string
		Token  string
	}{
		AppUrl: appUrl,
		Token:  verificationToken,
	}

	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, data); err != nil {
		logger.Info("Ошибка при выполнении шаблона: %v", zap.Error(err))
		return "", err
	}
	return emailBody.String(), nil
}

func ExtractUserFilmParams(c *fiber.Ctx) (userID, filmID string, typeUserFilm dto.TypeUserFilm, err error) {
	userID = c.Locals("userClaims").(*dto.Claims).UserID
	filmID = c.Params("filmId")

	typeUserFilmReq := c.Query("type")
	if err := ParseTypeUserFilm(typeUserFilmReq, &typeUserFilm); err != nil {
		return "", "", "", sys.NewError(sys.ErrInvalidRequestData, err.Error())
	}
	return userID, filmID, typeUserFilm, nil
}

func ParseTypeUserFilm(s string, t *dto.TypeUserFilm) error {
	switch s {
	case string(dto.TypeUserFavourite):
		*t = dto.TypeUserFavourite
	case string(dto.TypeUserRecommend):
		*t = dto.TypeUserRecommend
	default:
		return errors.New("Неверный тип фильма: " + s)
	}
	return nil
}
