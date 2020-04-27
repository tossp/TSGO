package setting

import (
	"crypto/ecdsa"
	"fmt"
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
	appPath      string
	globalKey    *ecdsa.PrivateKey
	globalAseKey []byte
	gmKey        *sm2.PrivateKey
)

func init() {
	appPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	appPath = filepath.Clean(appPath)
	if err := os.Chdir(appPath); err != nil {
		panic(err)
	}
	config()
	watch()
	globalKey = crypto.NewKeyWithKey([]byte(GetSecret()))
	gmKey = crypto.NewSm2KeyWithKey([]byte(GetSecret()))
	globalAseKey, _ = crypto.GenerateSharedSecret(globalKey, &crypto.NewKeyWithKey([]byte("砼砼")).PublicKey)
	if viper.GetString("control.pass") == "" {
		SetKeyPass("tossp")
	}
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
