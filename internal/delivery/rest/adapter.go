package rest

import "github.com/gofiber/fiber/v2"

type FilmController interface {
	Search(c *fiber.Ctx) error
	GetOneByID(c *fiber.Ctx) error
}

type FilmSequelController interface {
	GetAll(c *fiber.Ctx) error
}

type FilmSimilarController interface {
	GetAll(c *fiber.Ctx) error
}

type ExternalController interface {
	UploadFile(c *fiber.Ctx) error
}

type CollectionController interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetOne(c *fiber.Ctx) error
	GetAllMy(c *fiber.Ctx) error
}

type CollectionFilmController interface {
	Add(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetFilmsByCollectionId(c *fiber.Ctx) error
}

type HistoryFilmsController interface {
	Add(c *fiber.Ctx) error
}

type UserController interface {
	GetOneByNickName(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}

type UserFilmController interface {
	GetAll(c *fiber.Ctx) error
	Add(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type CommentController interface {
	GetAllByFilmID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type AuthController interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Verify(c *fiber.Ctx) error
	Me(c *fiber.Ctx) error
}
