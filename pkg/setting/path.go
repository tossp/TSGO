package setting

import (
	"os"
	"path/filepath"
)

func joinPath(base string, elem ...string) string {
	return filepath.Join(base, filepath.FromSlash(filepath.Clean(filepath.Join(elem...))))
}

// UseAppPath 使用系统路径
func UseAppPath(elem ...string) string {
	return joinPath(appPath, elem...)
}

// UseDataPath 使用数据存储路径
func UseDataPath(elem ...string) string {
	return joinPath(dataPath, elem...)
}

// UseConfigPath 使用配置存储路径
func UseConfigPath(elem ...string) string {
	return joinPath(configPath, elem...)
}

// MakePath 制作路径函数
func MakePath(base string) func(...string) string {
	_ = os.MkdirAll(base, 0700)
	return func(elem ...string) string {
		return filepath.Join(base, filepath.FromSlash(filepath.Clean(filepath.Join(elem...))))
	}
}
