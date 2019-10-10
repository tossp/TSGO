package setting

import (
	"github.com/tossp/tsgo/pkg/utils/crypto"

	"github.com/spf13/viper"
)

var ()

func GetKeyPass() string {
	return deKeyPass(viper.GetString("control.pass"))
}

func SetKeyPass(pass string) {
	pass = enKeyPass(pass)
	viper.Set("control.pass", pass)
	write()
}

func enKeyPass(pass string) string {
	pass = crypto.HexEncode(crypto.AesEncrypt([]byte(pass), globalAseKey))
	return pass
}

func deKeyPass(pass string) string {
	pass = string(crypto.AesDecrypt(crypto.HexDecode(pass), globalAseKey))
	return pass
}
