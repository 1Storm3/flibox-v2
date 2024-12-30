package rest

import (
	"github.com/1Storm3/flibox-api/internal/delivery/middleware"
	"github.com/1Storm3/flibox-api/internal/dto"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) setCommentRoutes(commentRoute fiber.Router) {
	commentRoute.Get("by/:filmId", r.commentController.GetAllByFilmID)
	commentRoute.Post("", middleware.ValidateMiddleware[dto.CreateCommentDTO](), r.commentController.Create)
	commentRoute.Delete(":id", r.commentController.Delete)
	commentRoute.Patch(":id", middleware.ValidateMiddleware[dto.UpdateCommentDTO](), r.commentController.Update)
}

func (r *Router) setHistoryFilmsRoutes(historyFilmsRoute fiber.Router) {
	historyFilmsRoute.Post(":Id", r.historyFilmsController.Add)
}

func (r *Router) setExternalRoutes(externalRoute fiber.Router) {
	externalRoute.Put("", r.externalController.UploadFile)
}
func (r *Router) setSequelRoutes(sequelRoute fiber.Router) {
	sequelRoute.Get(":id", r.filmSequelController.GetAll)
}

func (r *Router) setSimilarRoutes(similarRoute fiber.Router) {
	similarRoute.Get(":id", r.filmSimilarController.GetAll)
}

func (r *Router) setAuthRoutes(authRoute fiber.Router, authMiddleware fiber.Handler) {
	authRoute.Post("login", middleware.ValidateMiddleware[dto.LoginDTO](), r.authController.Login)
	authRoute.Post("register", middleware.ValidateMiddleware[dto.CreateUserDTO](), r.authController.Register)
	authRoute.Put("me", authMiddleware, r.authController.Me)
	authRoute.Post("verify/:token", r.authController.Verify)
}

func (r *Router) setUserRoutes(userRoute fiber.Router) {
	userRoute.Get(":nickName", r.userController.GetOneByNickName)
	userRoute.Patch(":id", r.userController.Update)
}

func (r *Router) setFilmRoutes(filmRoute fiber.Router) {
	filmRoute.Get(":id", r.filmController.GetOneByID)
	filmRoute.Get("", r.filmController.Search)
}

func (r *Router) setCollectionRoutes(collectionRoute fiber.Router) {
	collectionRoute.Get("", r.collectionController.GetAll)
	collectionRoute.Get("my", r.collectionController.GetAllMy)
	collectionRoute.Get(":id", r.collectionController.GetOne)

	collectionRoute.Post("", middleware.ValidateMiddleware[dto.CreateCollectionDTO](), r.collectionController.Create)
	collectionRoute.Delete(":id", r.collectionController.Delete)
	collectionRoute.Patch(":id", middleware.ValidateMiddleware[dto.UpdateCollectionDTO](), r.collectionController.Update)
}

func (r *Router) setCollectionFilmRoutes(collectionFilmRoute fiber.Router) {
	collectionFilmRoute.Post(":id/film", middleware.ValidateMiddleware[dto.CreateCollectionFilmDTO](), r.collectionFilmController.Add)
	collectionFilmRoute.Delete(":id/film", middleware.ValidateMiddleware[dto.DeleteCollectionFilmDTO](), r.collectionFilmController.Delete)
	collectionFilmRoute.Get(":id/films", r.collectionFilmController.GetFilmsByCollectionId)
}

func (r *Router) setUserFilmRoutes(userFilmRoute fiber.Router) {
	userFilmRoute.Get("/", r.userFilmController.GetAll)
	userFilmRoute.Post("/:filmId", r.userFilmController.Add)
	userFilmRoute.Delete("/:filmId", r.userFilmController.Delete)
}
