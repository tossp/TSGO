package setting

import (
	"github.com/spf13/viper"
)

func StorageBucketName() string {
	return viper.GetString("storage.Bucket")
}
func StorageEndpoint() string {
	return viper.GetString("storage.Endpoint")
}
func StorageAccessKey() string {
	return viper.GetString("storage.AccessKey")
}
func StorageSecretKey() string {
	return viper.GetString("storage.SecretKey")
}
func StorageSecure() bool {
	return viper.GetBool("storage.Secure")
}
