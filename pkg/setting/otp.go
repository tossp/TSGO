package setting

import (
	"github.com/spf13/viper"
	"github.com/tossp/tsgo/pkg/log"
	"github.com/tossp/tsgo/pkg/utils/crypto"
	"github.com/tossp/tsgo/pkg/utils/otp"
)

var (
	otpKey, _ = otp.NewKeyFromURL("")
)

func init() {
	otpauth, err := crypto.Base64Decode("b3RwYXV0aDovL3RvdHAvZGV2ZWxvcDpjb2RlQFRvc3NQLmNvbT9hbGdvcml0aG09U0hBNTEyJmRpZ2l0cz02Jmlzc3Vlcj1kZXZlbG9wJnBlcmlvZD05OTkmc2VjcmV0PTNYSlRTR05HUTRCSzdGMkJKWUY4QVhKTllEVjNRSFg2")
	if err != nil {
		log.Panic(err)
	}
	otpKey, err = otp.NewKeyFromURL(string(otpauth))
	if err != nil {
		log.Panic(err)
	}
	viper.SetDefault("env.devotp", true)
}

func ValidateDevelopKey(passcode string, msg ...string) bool {
	if passcode != "" && passcode == viper.GetString("env.devotpstr") {
		return true
	}
	if len(passcode) != otpKey.Digits().Length() || !viper.GetBool("env.devotp") {
		return false
	}
	ok, err := otp.ValidateByKey(passcode, otpKey)
	if err != nil {
		log.Error("开发密钥验证错误", err)
		return false
	}
	if ok {
		log.With("msg", msg).Debug("开发密钥被验证")
	}
	return ok
}

func GenerateDevelopKey() (string, error) {
	return otp.GenerateCodeCustomByKey(otpKey)
}
