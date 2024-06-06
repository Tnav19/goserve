package startup

import (
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact"
	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"github.com/unusualcodeorg/go-lang-backend-architecture/middleware"
	"github.com/unusualcodeorg/go-lang-backend-architecture/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/network"
)

func Server() {
	env := config.NewEnv(".env")

	db := mongo.NewDatabase(env)
	db.Connect()

	router := network.NewRouter()

	router.LoadControllers(
		contact.NewContactController(contact.NewService(db)),
	)

	router.LoadMiddlewares(
		middleware.NewNotFoundMiddleware(),
	)

	router.Start(env.ServerHost, env.ServerPort)

	defer db.Disconnect()
}
