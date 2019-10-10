package setting

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetAccessControlEnable() bool {
	return viper.GetBool("accessControl.Enable")
}

func GetAccessControlModel() string {
	return viper.GetString("accessControl.Model")
}

func GetAccessControlPrefix() string {
	return fmt.Sprintf("%s_", DbPrefix())
}
