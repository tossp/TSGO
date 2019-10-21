package jwt

import (
	"fmt"
	"time"

	"github.com/tossp/tsgo/pkg/log"
	"github.com/tossp/tsgo/pkg/null"
	"github.com/tossp/tsgo/pkg/setting"
	"github.com/tossp/tsgo/pkg/utils/crypto"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	expHour = time.Hour * 6

	Bearer    = "Bearer "
	XTseToken = "X-Ts-Token"
	CookieKey = "ts-token"
	headerF   = "Tse"
)

type IUser interface {
	New() IUser
	GetByID(null.UUID) error
	ID() null.UUID
}

type user struct{}

func (m *user) New() IUser {
	panic("请使用 jwt.SetUserMode 初始化默认用户接口")
}

func (m *user) GetByID(id null.UUID) error {
	panic("请使用 jwt.SetUserMode 初始化默认用户接口")
}
func (m *user) ID() null.UUID {
	panic("请使用 jwt.SetUserMode 初始化默认用户接口")
}

func setUserMode(u IUser) {
	defUser = u
}

var (
	tokenKey       = crypto.NewKeyWithKey([]byte(setting.GetSecret()))
	defUser  IUser = new(user)
)

type TsClaims struct {
	jwt.StandardClaims
	UserID null.UUID `json:"user_id"`
}

//GenerateToken 生成Token
func GenerateToken(id null.UUID, ct time.Time) (t string) {
	claims := new(TsClaims)
	claims.ExpiresAt = ct.Add(expHour).Unix()
	claims.NotBefore = ct.Unix()

	claims.UserID = id

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

//ParseToken 预处理Token
func validJwt(auth string) (user IUser, claims *TsClaims, err error) {
	l := len(Bearer)
	if len(auth) > l+1 && auth[:l] == Bearer {
		t, fuck := parseToken(auth[l:])
		if fuck != nil {
			err = fuck
			return
		}
		if data, ok := t.Claims.(*TsClaims); ok && t.Valid {
			user = defUser.New()
			err = user.GetByID(data.UserID)
			claims = data
		} else {
			err = fmt.Errorf("没有找到可用的签名")
		}
	}
	return
}
