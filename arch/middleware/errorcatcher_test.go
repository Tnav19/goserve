package middleware

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/unusualcodeorg/goserve/arch/network"
)

func TestErrorCatcherMiddleware(t *testing.T) {
	_, rr := network.MockTestRootMiddleware(t, "GET", "/test", "/test", "",
		NewErrorCatcher(),
		func(ctx *gin.Context) {
			panic(errors.New("panic test"))
		})

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"panic test"`)
}

func TestErrorCatcherMiddleware_NonError(t *testing.T) {
	_, rr := network.MockTestRootMiddleware(t, "GET", "/test", "/test", "",
		NewErrorCatcher(),
		func(ctx *gin.Context) {
			panic(1)
		})

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"something went wrong"`)
}
