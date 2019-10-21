package setting

import (
	"crypto/ecdsa"
	"os"
	"path/filepath"

	"github.com/tossp/tsgo/pkg/utils/crypto"

	"github.com/spf13/viper"
)

const (
	configFileName = "ts_config"
	DataDirKey     = "DataDir"
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
	globalAseKey = crypto.GenerateSharedSecret(globalKey, &crypto.NewKeyWithKey([]byte("砼砼")).PublicKey)
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
