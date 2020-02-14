package jwt

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tossp/tsgo/pkg/casbin"
	"github.com/tossp/tsgo/pkg/log"
	"github.com/tossp/tsgo/pkg/utils/crypto"
)

func EchoEnforcer() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			u := c.Get(authUserKey).(IUser)
			if u == nil {
				return echo.NewHTTPError(http.StatusForbidden, "没有找到可用的token")
			}
			obj := c.Request().URL.Path
			act := c.Request().Method
			if ok, err := casbin.E().Enforce(u.ID().String(), obj, act, "*"); err != nil {
				log.Error("RBAC", err)
				c.Set(authorityKey, "ERR")
				return echo.NewHTTPError(http.StatusForbidden, "鉴权系统错误，请稍后再试")
			} else if ok {
				c.Set(authorityKey, "RBAC")
				return next(c)
			} else if u.HasAdmin() {
				c.Set(authorityKey, "ADMIN")
				c.Response().Header().Set("x-ts-authority-admin", crypto.Base64Encode([]byte(fmt.Sprintf("%s;%s", act, obj))))
				return next(c)
			}
			c.Set(authorityKey, "FAIL")
			return echo.NewHTTPError(http.StatusForbidden, "权限不足")
		}
	}
}

func GetAuthority(c echo.Context) (data string) {
	data, ok := c.Get(authorityKey).(string)
	if !ok {
		data = ""
	}
	return
}
