package setting

import (
	"github.com/spf13/viper"
)

func GetSecret() string {
	return viper.GetString("secret")
}
func GetStaticPageURL() string {
	return viper.GetString("static")
}
