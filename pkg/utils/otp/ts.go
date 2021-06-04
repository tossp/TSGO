package otp

import (
	"time"

	"github.com/tossp/tsgo/pkg/utils/otp/pquerna"
	"github.com/tossp/tsgo/pkg/utils/otp/pquerna/totp"
)

func NewKeyFromURL(orig string) (*pquerna.Key, error) {
	return pquerna.NewKeyFromURL(orig)
}

func ValidateByKey(passcode string, key *pquerna.Key) (bool, error) {
	return totp.ValidateCustom(passcode, key.Secret(), time.Now(), totp.ValidateOpts{
		Period:    uint(key.Period()),
		Skew:      1,
		Digits:    key.Digits(),
		Algorithm: key.Algorithm(),
	})
}

func GenerateCodeCustomByKey(key *pquerna.Key) (string, error) {
	return totp.GenerateCodeCustom(key.Secret(), time.Now(), totp.ValidateOpts{
		Period:    uint(key.Period()),
		Skew:      1,
		Digits:    key.Digits(),
		Algorithm: key.Algorithm(),
	})
}
