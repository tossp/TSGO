package setting

import (
	"github.com/spf13/viper"
)

func ControlPass() string {
	return viper.GetString("control.pass")
}
