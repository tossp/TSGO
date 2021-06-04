package jwt

import (
	"fmt"
	"net"
	"time"

	"github.com/tossp/tsgo/pkg/log"
	"github.com/tossp/tsgo/pkg/setting"
	"github.com/tossp/tsgo/pkg/tstype"
	"github.com/tossp/tsgo/pkg/utils/crypto"

	"github.com/dgrijalva/jwt-go"
)

const (
	Bearer    = "Bearer "
	XTseToken = "X-Ts-Token"
	CookieKey = "ts-token"
)

var (
	tokenKey        = crypto.NewKeyWithKey([]byte(setting.GetSecret() + "TossP.com"))
	expiresDuration = time.Hour * 24
)

func init() {

}

type IUser interface {
	New() IUser
	GetByID(tstype.UUID) error
	ID() tstype.UUID
	HasAdmin() bool
	OnlineExtend(*TsClaims) error
	OnlineCheck(*TsClaims, net.IP) error
}

type user struct{}

func (m *user) New() IUser {
	panic("请使用 jwt.SetUserMode 初始化默认用户接口")
}

func (m *user) GetByID(id tstype.UUID) error {
	panic("请使用 jwt.SetUserMode 初始化默认用户接口")
}
func (m *user) ID() tstype.UUID {
	panic("请使用 jwt.SetUserMode 初始化默认用户接口")
}
func (m *user) HasAdmin() bool {
	panic("请使用 jwt.SetUserMode 初始化默认用户接口")
}

func (m *user) OnlineExtend(c *TsClaims) error {
	panic("请使用 jwt.SetUserMode 初始化默认用户接口")
}
func (m *user) OnlineCheck(c *TsClaims, ip net.IP) error {
	panic("请使用 jwt.SetUserMode 初始化默认用户接口")
}

type TsClaims struct {
	jwt.StandardClaims
	UserID tstype.UUID `json:"usi,omitempty"`
}

//Extend 延期Token
func (c *TsClaims) Extend(ct time.Time) *TsClaims {
	c.ExpiresAt = ct.Add(expiresDuration).Unix()
	c.NotBefore = ct.Unix()
	return c
}

//Extend 延期Token
func (c *TsClaims) SignedString() (t string) {
	t, err := jwt.NewWithClaims(jwt.SigningMethodES512, c).SignedString(tokenKey)
	if err != nil {
		log.Warn("生成token错误", err)
	}
	return
}

func init() {
	setting.SetDefault("auth.Timeout", 30)
	_ = setting.Subscribe(ReadTimeout)

}

func ReadTimeout() (timeout int64) {
	timeout = setting.GetInt64("auth.Timeout")
	if timeout < 1 {
		timeout = 60
		setting.Set("auth.Timeout", timeout)
	}
	expiresDuration = time.Minute * time.Duration(timeout)
	return
}
func SetTimeout(timeout int64) int64 {
	if timeout < 1 || timeout > 30 {
		timeout = 30
	}
	setting.Set("auth.Timeout", timeout)
	return ReadTimeout()
}

//GenerateToken 生成Token
func GenerateToken(id tstype.UUID, ct time.Time) (claims *TsClaims, t string) {
	claims = new(TsClaims)
	claims.ExpiresAt = ct.Add(expiresDuration).Unix()
	claims.NotBefore = ct.Unix()
	claims.IssuedAt = time.Now().Unix()

	claims.UserID = id
	claims.Id = utils.NewUuidV4().String()

	token := jwt.NewWithClaims(jwt.SigningMethodES512, claims)
	t, err := token.SignedString(tokenKey)
	//t, err := token.SignedString(crypto.HashKey([]byte(setting.Data().SecretKey), 32))
	if err != nil {
		log.Warn("生成token错误", err)
	}
	return
}

//ParseToken 预处理Token
func parseToken(token string) (t *jwt.Token, err error) {
	t, err = jwt.ParseWithClaims(token, &TsClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 检查签名模型
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("错误的签名模型: %v", token.Header["alg"])
		}
		// Return the key for validation
		return &tokenKey.PublicKey, nil
	})
	return
}

//ValidToken 验证Token
func ValidToken(u IUser, auth string) (user IUser, claims *TsClaims, err error) {
	l := len(Bearer)
	if len(auth) > l+1 && auth[:l] == Bearer {
		var t *jwt.Token
		if t, err = parseToken(auth[l:]); err != nil {
			return
		}
		if data, ok := t.Claims.(*TsClaims); ok && t.Valid {
			user = u.New()
			err = user.GetByID(data.UserID)
			claims = data
			return
		}
	}
	err = fmt.Errorf("没有找到可用的签名")
	return
}
