package setting

import (
	"github.com/spf13/viper"
)

func WebAddress() string {
	return viper.GetString("web.address")
}
