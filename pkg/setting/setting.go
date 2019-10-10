package setting

import (
	"os"
	"path"
	"path/filepath"
)

const (
	configFileName = "ts_config"
	DataDirKey     = "DataDir"
	ConfigDirKey   = "ConfigDir"
	AppName        = "TS_SITES"
)

var (
	appPath string
)

func init() {
	appPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	appPath = filepath.Clean(appPath)
	config()
	watch()
}

func UseAppPath(elem ...string) string {
	return joinPath(appPath, elem...)
}

func joinPath(base string, elem ...string) string {
	return path.Join(base, filepath.Clean(path.Join(elem...)))
}

func UseDataPath(elem ...string) string {
	return filepath.Clean(path.Join(append([]string{appPath})...))
	return UseAppPath()
}
