package network

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BaseController interface {
	Path() string
	AuthenticationMiddleware() GroupMiddleware
	AuthorizationMiddleware() GroupMiddleware
}

type BaseService interface {
	Context() (context.Context, context.CancelFunc)
}

type Controller interface {
	BaseController
	MountRoutes(group *gin.RouterGroup)
}

type Dto[T any] interface {
	Payload() *T
	ValidateErrors(errs validator.ValidationErrors) ([]string, error)
}

type RootMiddleware interface {
	Attach(engine *gin.Engine)
	Handler(ctx *gin.Context)
}

type GroupMiddleware interface {
	Attach(group *gin.RouterGroup)
	Handler(ctx *gin.Context)
}

type GroupMiddlewareFunc func() GroupMiddleware

type Router interface {
	GetEngine() *gin.Engine
	LoadControllers(controllers ...Controller)
	LoadRootMiddlewares(middlewares ...RootMiddleware)
	Start(ip string, port uint16)
}
