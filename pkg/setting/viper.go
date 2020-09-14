package setting

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/tossp/tsgo/pkg/log"
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
)

func init() {

}

func config() {

read:
	err := viper.ReadInConfig()
	if err != nil {
		log.Warnf("配置文件错误: %s\n", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			configFN := joinPath(viper.GetString(ConfigDirKey), configFileName+".toml")
			f, fuck := os.OpenFile(configFN, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if fuck != nil {
				log.Warnf("尝试修复配置文件错误: %s\n", fuck)
				os.Exit(1)
			}
			_ = f.Close()
			log.Warnf("修复配置: %s\n", configFN)
			goto read
		}
		panic(err)
	}
	overwrite()
	log.SetConfig(IsDev(), viper.GetString(LogDirKey))
	log.Infow("加载配置文件成功", "filename", viper.ConfigFileUsed())
	Write()
}
func overwrite() {
	_ = os.MkdirAll(viper.GetString(DataDirKey), 0755)
	_ = os.MkdirAll(viper.GetString(ConfigDirKey), 0755)
}

func watch() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Debug("配置发生变更：", e.Op.String())
		//oldConfig := viper.New()
		//for _, v := range viper.AllKeys() {
		//	oldConfig.Set(v, viper.Get(v))
		//}
		//if err := viper.ReadInConfig(); err != nil {
		//	log.Warnf("重新加载配置文件错误: %s", viper.ConfigFileUsed())
		//	for _, v := range oldConfig.AllKeys() {
		//		viper.Set(v, viper.Get(v))
		//	}
		//	Write()
		//	return
		//}
		overwrite()
		log.SetConfig(IsDev(), viper.GetString(LogDirKey))
		log.Infow("加载配置文件成功", "filename", viper.ConfigFileUsed())
	})
}

func Write() {
	err := viper.WriteConfig()
	if err != nil {
		log.Errorf("保存配置文件错误: %s\n", err)
	}
}
