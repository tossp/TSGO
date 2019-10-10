package setting

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func DbUser() string {
	return viper.GetString("db.User")
}
func DbPassword() string {
	return viper.GetString("db.Password")
}

func DbPrefix() string {
	return fmt.Sprintf("tsl_%s", strings.ToLower(viper.GetString("db.Prefix")))
}

func DbHost() string {
	return viper.GetString("db.Host")
}

func DbPort() int64 {
	return viper.GetInt64("db.Port")
}

func DbName() string {
	return viper.GetString("db.Name")
}

func DbMode() string {
	return viper.GetString("db.Ssl_mode")
}

func DbMaxIdleConns() int {
	return viper.GetInt("db.Max_Idle_Conns")
}

func DbMaxOpenConns() int {
	return viper.GetInt("db.Max_Open_Conns")
}
