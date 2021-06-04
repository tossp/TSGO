package setting

import (
	"github.com/spf13/viper"
)

// IsDev 是否为开发模式
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

// SetMod 设置模式
func SetMod(m string) {
	viper.Set("mod", m)
}
func GetSecret() string {
	return viper.GetString("secret")
}
func GetStaticPageURL() string {
	return viper.GetString("static")
}
