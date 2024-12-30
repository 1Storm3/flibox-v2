package rest

import (
	"github.com/1Storm3/flibox-api/internal/delivery/middleware"
	"github.com/gofiber/fiber/v2"
)

const (
	Admin = "admin"
	User  = "user"
)

type Router struct {
	filmController           FilmController
	filmSequelController     FilmSequelController
	userController           UserController
	filmSimilarController    FilmSimilarController
	userFilmController       UserFilmController
	authController           AuthController
	externalController       ExternalController
	commentController        CommentController
	collectionController     CollectionController
	collectionFilmController CollectionFilmController
	historyFilmsController   HistoryFilmsController
}

func NewRouter(
	filmController FilmController,
	filmSequelController FilmSequelController,
	userController UserController,
	filmSimilarController FilmSimilarController,
	userFilmController UserFilmController,
	authController AuthController,
	externalController ExternalController,
	commentController CommentController,
	collectionController CollectionController,
	collectionFilmController CollectionFilmController,
	historyFilmsController HistoryFilmsController,
) *Router {
	return &Router{
		filmController:           filmController,
		filmSequelController:     filmSequelController,
		userController:           userController,
		filmSimilarController:    filmSimilarController,
		userFilmController:       userFilmController,
		authController:           authController,
		externalController:       externalController,
		commentController:        commentController,
		collectionController:     collectionController,
		collectionFilmController: collectionFilmController,
		historyFilmsController:   historyFilmsController,
	}
}

func (r *Router) LoadRoutes(app fiber.Router, authMiddleware fiber.Handler) {
	apiRoute := app.Group("api")

	authRoute := apiRoute.Group("auth")
	r.setAuthRoutes(authRoute, authMiddleware)

	userFilmRoute := apiRoute.Group("user/my", authMiddleware, middleware.RoleMiddleware(Admin, User))
	r.setUserFilmRoutes(userFilmRoute)

	userRoute := apiRoute.Group("user", authMiddleware, middleware.RoleMiddleware(Admin, User))
	r.setUserRoutes(userRoute)

	filmRoute := apiRoute.Group("film")
	r.setFilmRoutes(filmRoute)

	sequelRoute := apiRoute.Group("sequel")
	r.setSequelRoutes(sequelRoute)

	similarRoute := apiRoute.Group("similar")
	r.setSimilarRoutes(similarRoute)

	externalRoute := apiRoute.Group("upload", authMiddleware, middleware.RoleMiddleware(Admin, User))
	r.setExternalRoutes(externalRoute)

	commentRoute := apiRoute.Group("comment", authMiddleware, middleware.RoleMiddleware(Admin, User))
	r.setCommentRoutes(commentRoute)

	collectionRoute := apiRoute.Group("collection", authMiddleware, middleware.RoleMiddleware(Admin, User))
	r.setCollectionRoutes(collectionRoute)

	collectionFilmRoute := apiRoute.Group("collection", authMiddleware, middleware.RoleMiddleware(Admin, User))
	r.setCollectionFilmRoutes(collectionFilmRoute)

	historyFilmsRoute := apiRoute.Group("film/history", authMiddleware, middleware.RoleMiddleware(Admin, User))
	r.setHistoryFilmsRoutes(historyFilmsRoute)
}
