package contact

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/common/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/common/parser"
)

type controller struct {
	contactService ContactService
}

func NewContactController(s ContactService) network.Controller {
	cnt := controller{
		contactService: s,
	}
	return &cnt
}

func (c *controller) MountRoutes(router *gin.Engine) {
	router.POST("/message", c.createMessageHandler)
}

func (c *controller) createMessageHandler(ctx *gin.Context) {
	var createMsg dto.CreateMessage

	if err := parser.GetBody(ctx, &createMsg); err != nil {
		network.BadRequestResponse(err).Send(ctx)
		return
	}

	_, err := c.contactService.SaveMessage(createMsg.Type, createMsg.Msg)

	if err != nil {
		network.InternalServerErrorResponse("Something went wrong")
		return
	}

	network.SuccessMsgResponse("Message received successfully!").Send(ctx)
}
