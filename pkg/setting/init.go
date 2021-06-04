package setting

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/tossp/tsgo/pkg/client"
	"github.com/tossp/tsgo/pkg/log"
	"github.com/tossp/tsgo/pkg/utils"
	"gopkg.in/resty.v1"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tossp/tsgo/pkg/utils/crypto"

	"github.com/spf13/viper"
)

const (
	configFileName = "ts"
	appName        = "TS_GO_"
)

var (
	GitVersion   = "invalid"
	BuildTime    = "invalid"
	BuildVersion = "invalid"
	BuildUser    = "invalid"
	ProjectName  = "invalid"
)

var (
	appPath    string
	configPath string
	dataPath   string
	logPath    string

	globalKey *ecdsa.PrivateKey
	gmKey     *sm2.PrivateKey
)

func init() {
	client.SetUserAgent(fmt.Sprintf("tossp/%s (compatible; go/%s; %s %s; +https://github.com/tossp)", resty.Version, runtime.Version(), runtime.GOOS, runtime.GOARCH))

	appPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	appPath = filepath.Clean(appPath)
	if err := os.Chdir(appPath); err != nil {
		panic(err)
	}
	configPath = UseAppPath("configs")
	dataPath = UseAppPath("data")
	logPath = UseDataPath("logs")
	_ = os.MkdirAll(configPath, 0755)
	_ = os.MkdirAll(dataPath, 0755)
	_ = os.MkdirAll(logPath, 0755)
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(configPath)

	viper.SetDefault("secret", utils.GetRandomString(32))
	viper.SetDefault("mod", "prod")
	viper.SetDefault("lv", "info")
	isNew := false
read:
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			touchConfigFilename()
			isNew = true
			goto read
		}
		panic(err)
	}
	log.SetConfig(IsDev(), logPath, GetString("lv"))
	if isNew {
		Write()
	}
	globalKey = crypto.NewKeyWithKey([]byte(GetSecret()))
	gmKey = crypto.NewSm2KeyWithKey([]byte(GetSecret()))
	log.Infof("配置文件：%s", viper.ConfigFileUsed())

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.SetConfig(IsDev(), logPath, GetString("lv"))
		log.Infof("配置文件[%s]：%s", e.Op.String(), viper.ConfigFileUsed())
		bus.Publish("setting")
	})
	tsSentry()
}

// AppName 获取应用名称
func AppName() string {
	return appName + ProjectName
}

func GetGlobalPubKey() string {
	return crypto.HexEncode(crypto.FromECDSAPub(&globalKey.PublicKey))
}

func GetJsGlobalPubKey() string {
	return fmt.Sprintf("%s|%s", crypto.Base64Encode(crypto.FromECDSAPub(&globalKey.PublicKey)), crypto.Base64Encode(crypto.FromsSm2Pub(&gmKey.PublicKey)))
}

func GetGlobalKey() *ecdsa.PrivateKey {
	return globalKey
}

func GmKey() *sm2.PrivateKey {
	return gmKey
}

func touchConfigFilename() {
	f, fuck := os.OpenFile(joinPath(configPath, configFileName+".toml"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if fuck != nil {
		panic(fuck)
	}
	_ = f.Close()
}
