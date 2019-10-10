package setting

import (
	"crypto/ecdsa"
	"os"
	"path"
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
	GitTag       = "syntax error"
	GitVersion   = "syntax error"
	BuildTime    = "syntax error"
	BuildVersion = "syntax error"
	ProjectName  = "syntax error"
	appPath      string
	globalKey    *ecdsa.PrivateKey
	globalAseKey []byte
)

func init() {
	appPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	appPath = filepath.Clean(appPath)
	config()
	watch()
	globalKey = crypto.NewKeyWithKey([]byte(GetSecret()))
	globalAseKey = crypto.GenerateSharedSecret(globalKey, &crypto.NewKeyWithKey([]byte("砼砼")).PublicKey)
	if viper.GetString("control.pass") == "" {
		SetKeyPass("mcit")
	}
}

func UseAppPath(elem ...string) string {
	return joinPath(appPath, elem...)
}

func joinPath(base string, elem ...string) string {
	return path.Join(base, filepath.Clean(path.Join(elem...)))
}

func UseDataPath(elem ...string) string {
	return UseAppPath(elem...)
}
