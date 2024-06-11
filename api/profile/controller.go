package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/profile/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core"
	coredto "github.com/unusualcodeorg/go-lang-backend-architecture/core/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type controller struct {
	network.BaseController
	userService    core.UserService
	profileService ProfileService
}

func NewProfileController(
	authProvider network.AuthenticationProvider,
	authorizeProvider network.AuthorizationProvider,
	userService core.UserService,
	profileService ProfileService,
) network.Controller {
	c := controller{
		BaseController: network.NewBaseController("/profile", authProvider, authorizeProvider),
		userService:    userService,
		profileService: profileService,
	}
	return &c
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.getUserHandler)
}

func (c *controller) getUserHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, &coredto.MongoId{})
	if err != nil {
		panic(network.BadRequestError(err.Error(), err))
	}

	msg, err := c.userService.FindUserById(mongoId.ID)
	if err != nil {
		panic(network.NotFoundError("message not found", err))
	}

	data, err := network.MapToDto[dto.InfoPrivateUser](msg)
	if err != nil {
		panic(network.InternalServerError("something went wrong", err))
	}

	network.SuccessDataResponse(ctx, "success", data)
}
