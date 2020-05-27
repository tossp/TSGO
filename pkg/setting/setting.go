package setting

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/tossp/tsgo/pkg/utils"
	"os"
	"path/filepath"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tossp/tsgo/pkg/utils/crypto"

	"github.com/spf13/viper"
)

const (
	configFileName = "ts_config"
	DataDirKey     = "DataDir"
	LogDirKey      = "LogDir"
	ConfigDirKey   = "ConfigDir"
	AppName        = "TS_SITES"
)

var (
	GitTag       = "invalid"
	GitVersion   = "invalid"
	BuildTime    = "invalid"
	BuildVersion = "invalid"
	ProjectName  = "invalid"
)
var (
	appPath   string
	globalKey *ecdsa.PrivateKey
	gmKey     *sm2.PrivateKey
)

func init() {
	appPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	appPath = filepath.Clean(appPath)
	if err := os.Chdir(appPath); err != nil {
		panic(err)
	}

	viper.SetDefault(ConfigDirKey, UseConfigPath("configs"))
	viper.SetDefault(DataDirKey, UseAppPath("data"))
	viper.SetDefault(LogDirKey, UseAppPath("data", "logs"))
	viper.SetDefault("secret", utils.GetRandomString(32))
	viper.SetDefault("mod", "prod")

	//viper.SetEnvPrefix("ts")
	//viper.AutomaticEnv()
	//_ = viper.BindEnv(DataDirKey)
	//_ = viper.BindEnv(ConfigDirKey)

	flag.String(ConfigDirKey, viper.GetString(ConfigDirKey), "基础配置文件路径")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	_ = os.MkdirAll(viper.GetString(ConfigDirKey), 0755)

	viper.SetConfigName(configFileName)
	viper.AddConfigPath(viper.GetString(ConfigDirKey))
	//viper.AddConfigPath(UseAppPath("configs"))
	viper.AddConfigPath(appPath)
	//viper.AddConfigPath("$HOME/.ts_site")
	//viper.AddConfigPath(".")

	config()
	watch()
	globalKey = crypto.NewKeyWithKey([]byte(GetSecret()))
	gmKey = crypto.NewSm2KeyWithKey([]byte(GetSecret()))
}

func UseAppPath(elem ...string) string {
	return joinPath(appPath, elem...)
}

func joinPath(base string, elem ...string) string {
	return filepath.Join(base, filepath.Clean(filepath.Join(elem...)))
}

func UseDataPath(elem ...string) string {
	return joinPath(viper.GetString(DataDirKey), elem...)
}
func UseConfigPath(elem ...string) string {
	return joinPath(viper.GetString(ConfigDirKey), elem...)
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
