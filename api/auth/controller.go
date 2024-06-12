package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type controller struct {
	network.BaseController
	authService AuthService
}

func NewAuthController(
	authProvider network.AuthenticationProvider,
	authorizeProvider network.AuthorizationProvider,
	service AuthService,
) network.Controller {
	c := controller{
		BaseController: network.NewBaseController("/auth", authProvider, authorizeProvider),
		authService:    service,
	}
	return &c
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.POST("/signup/basic", c.signUpBasicHandler)
	group.POST("/signin/basic", c.signInBasicHandler)
	group.DELETE("/signout", c.Authentication(), c.signOutBasic)
}

func (c *controller) signUpBasicHandler(ctx *gin.Context) {
	body, err := network.ReqBody(ctx, dto.EmptySignUpBasic())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	data, err := c.authService.SignUpBasic(body)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", data)
}

func (c *controller) signInBasicHandler(ctx *gin.Context) {
	body, err := network.ReqBody(ctx, dto.EmptySignInBasic())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	dto, err := c.authService.SignInBasic(body)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", dto)
}

func (c *controller) signOutBasic(ctx *gin.Context) {
	keystore := network.ReqMustGetKeystore[model.Keystore](ctx)

	err := c.authService.SignOut(keystore)
	if err != nil {
		c.Send(ctx).InternalServerError("something went wrong", err)
		return
	}

	c.Send(ctx).SuccessMsgResponse("signout success")
}
