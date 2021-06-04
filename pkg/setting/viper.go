package setting

import (
	"github.com/spf13/viper"
	"github.com/tossp/tsgo/pkg/log"
	"github.com/tossp/tsgo/pkg/utils/bus"
)

var (
	SetDefault     = viper.SetDefault
	Set            = viper.Set
	Get            = viper.Get
	GetString      = viper.GetString
	GetStringSlice = viper.GetStringSlice
	GetInt64       = viper.GetInt64
	GetBool        = viper.GetBool
	GetInt         = viper.GetInt
	Subscribe      = subscribe
)

//Write 写出配置文件
func Write() {
	err := viper.WriteConfig()
	if err != nil {
		log.Errorf("保存配置文件错误: %s\n", err)
	}
}

func subscribe(fn interface{}) error {
	return bus.Subscribe("setting", fn)
}
