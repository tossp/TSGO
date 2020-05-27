package jwt

import (
	"fmt"
	"github.com/tossp/tsgo/pkg/utils"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	authErrKey       = "claims.err"
	authUserKey      = "claims.me"
	authorityKey     = "casbin.authority"
	authorityTextKey = "casbin.authority.text"
)

//EchoAuth jwt注入鉴定
func EchoJwt(u IUser) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if (c.Request().Header.Get(echo.HeaderUpgrade)) == "websocket" || // 跳过 WebSocket
				strings.Index(c.Request().URL.Path, "/v1/") != 0 { // 跳过 非验证模块
				return next(c)
			}

			cc, ok := c.(utils.LikeContextLog)
			if !ok {
				c.Set(authErrKey, "utils.LikeContextLog 未注入")
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
				user, claims, err := validJwt(u, auth)
				if err != nil {
					cc.Log("令牌", fmt.Sprintf("校验错误 %v", err))
					c.Set(authErrKey, err)
					return next(c)
				}
				if err = user.OnlineCheck(claims, cc.Ip()); err != nil {
					cc.Log("令牌", fmt.Sprintf("会话检查失败：%v", err))
					c.Set(authErrKey, err)
					return next(c)
				}

				c.Set("claims.exp", time.Unix(claims.ExpiresAt, 0))
				c.Set("claims", claims)
				c.Set(authUserKey, user)
				expTime := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
				if expTime > 0 && expTime <= expiresDuration/2 {
					oldExpiresAt := time.Unix(claims.ExpiresAt, 0).String()
					token := claims.Extend(time.Now()).SignedString()
					if err = user.OnlineExtend(claims); err != nil {
						cc.Log("令牌", fmt.Sprintf("延期失败：%v", err))
					} else {
						//c.Response().Header().Add(echo.HeaderAuthorization, token)
						c.Response().Header().Add(XTseToken, token)
						c.SetCookie(&http.Cookie{Name: CookieKey, Value: token, HttpOnly: true})
						cc.Log("令牌", fmt.Sprintf("将在%s过期，延期到%s", oldExpiresAt, time.Unix(claims.ExpiresAt, 0).String()))
					}
				}
				return next(c)
			}
			if cookie, err := c.Cookie(CookieKey); err == nil && !useCookie {
				auth = Bearer + cookie.Value
				useCookie = true
				goto AUTH
			}
			c.Set(authErrKey, "没有找到可用的token")
			cc.Log("令牌", "没有找到可用的token")
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
			if strings.Index(c.Request().URL.Path, "/v1/system/develop/debug/pprof") == 0 {
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
