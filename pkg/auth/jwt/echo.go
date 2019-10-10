package jwt

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	authErrKey  = "claims.err"
	authUserKey = "claims.me"
)

//EchoAuth jwt注入鉴定
func EchoJwt() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 跳过 WebSocket
			if (c.Request().Header.Get(echo.HeaderUpgrade)) == "websocket" {
				return next(c)
			}
			useCookie := false
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			BearerLen := len(Bearer)
			auth = c.Request().Header.Get(echo.HeaderAuthorization)
			if auth == "" {
				auth = c.Request().Header.Get(echo.HeaderAuthorization)
			}
		AUTH:
			if len(auth) > BearerLen+1 && auth[:BearerLen] == Bearer {
				user, claims, err := validJwt(auth)
				if err != nil {
					c.Set(authErrKey, err)
					return next(c)
				}
				c.Set("claims.exp", time.Unix(claims.ExpiresAt, 0))
				c.Set("claims", claims)
				c.Set(authUserKey, user)
				expTime := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
				if expTime > 0 && expTime < expHour/2 {
					token := GenerateToken(user.ID(), time.Now())
					c.Response().Header().Add(echo.HeaderAuthorization, token)
					c.SetCookie(&http.Cookie{Name: CookieKey, Value: token, HttpOnly: true})
				}
				return next(c)
			}
			if cookie, err := c.Cookie(CookieKey); err == nil && !useCookie {
				auth = Bearer + cookie.Value
				useCookie = true
				goto AUTH
			}
			c.Set(authErrKey, "没有找到可用的token")
			return next(c)
		}
	}
}
func EchoAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 跳过 WebSocket
			if (c.Request().Header.Get(echo.HeaderUpgrade)) == "websocket" {
				return next(c)
			}
			err := c.Get(authErrKey)
			if err == nil && c.Get(authUserKey) != nil {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
	}
}
