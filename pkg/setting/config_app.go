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

func IsDev() bool {
	switch viper.GetString("mod") {
	case "dev":
		return true
	case "prod":
		return false
	default:
		SetMod("prod")
	}
	return false
}
func SetMod(m string) {
	viper.Set("mod", m)
}
