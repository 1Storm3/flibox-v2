package http

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/1Storm3/flibox-api/internal/controller"
	"github.com/1Storm3/flibox-api/pkg/sys"
)

type ExternalController struct {
	externalService controller.ExternalService
	s3Service       controller.S3Service
}

func NewExternalController(externalService controller.ExternalService, s3Service controller.S3Service) *ExternalController {
	return &ExternalController{
		externalService: externalService,
		s3Service:       s3Service,
	}
}

func (e *ExternalController) UploadFile(c *fiber.Ctx) error {
	ctx := c.Context()
	file, err := c.FormFile("file")
	if err != nil {
		return sys.NewError(sys.ErrUnknown, err.Error())
	}
	fileReader, err := file.Open()
	if err != nil {
		return sys.NewError(sys.ErrUnknown, err.Error())
	}
	defer func(fileReader multipart.File) {
		err := fileReader.Close()
		if err != nil {
			return
		}
	}(fileReader)

	fileBytes, err := io.ReadAll(fileReader)
	if err != nil {
		return sys.NewError(sys.ErrUnknown, err.Error())
	}

	ext := filepath.Ext(file.Filename)
	uniqueID, err := uuid.NewUUID()
	if err != nil {
		return sys.NewError(sys.ErrUnknown, err.Error())
	}

	uniqueFilename := fmt.Sprintf("%s%s", uniqueID.String(), ext)

	url, err := e.s3Service.UploadFile(ctx, uniqueFilename, fileBytes)
	if err != nil {
		return sys.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"url": url,
	})
}
