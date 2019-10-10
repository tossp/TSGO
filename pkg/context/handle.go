package context

import (
	"net/http"

	"github.com/tossp/tsgo/pkg/errors"
	"github.com/tossp/tsgo/pkg/log"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *Context)

func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			c,
		}
		h(ctx)
	}
}

type Context struct {
	*gin.Context
}

func (c *Context) Success(data interface{}, meta ...interface{}) {
	re := gin.H{"code": 0, "data": data}
	if len(meta) > 0 {
		re["meta"] = meta[0]
	}
	if len(c.Errors) > 0 {
		re["err"] = c.Errors
	}
	c.AsciiJSON(http.StatusOK, re)
}

func (c *Context) Fail(err error, meta ...interface{}) {
	re := gin.H{}
	statusCode := http.StatusOK
	if len(meta) > 0 {
		re["meta"] = meta[0]
	}
	if err == nil {
		statusCode = http.StatusServiceUnavailable
		re["msg"] = "预期err为nil"
	} else {
		log.Debugf("%+v", err)
		re["msg"] = err
		switch true {
		case errors.Is(err, errors.ErrForbidden):
			statusCode = http.StatusForbidden
		case errors.Is(err, errors.ErrNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, errors.ErrBadRequest):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errors.ErrUnauthorized):
			statusCode = http.StatusUnauthorized
		case errors.Is(err, errors.ErrInternalServer):
			statusCode = http.StatusInternalServerError
		case errors.Is(err, errors.ErrDatabase):
		case errors.Is(err, errors.ErrExistsfail):
			statusCode = http.StatusInsufficientStorage
		case errors.Is(err, errors.ErrCodefail):
			statusCode = http.StatusVariantAlsoNegotiates
		default:
			statusCode = http.StatusServiceUnavailable
		}
	}
	if len(c.Errors) > 0 {
		re["err"] = c.Errors
	}
	c.Abort()
	c.AsciiJSON(statusCode, re)
}
