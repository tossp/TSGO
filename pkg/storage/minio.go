package storage

import (
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	"time"

	minio "github.com/minio/minio-go/v6"
	"github.com/minio/minio/pkg/madmin"
	"github.com/tossp/tsgo/pkg/errors"
)

const expires = time.Hour

var (
	minioClient *minio.Client
	mdmClnt     *madmin.AdminClient
	bucketName  = viper.GetString("storage.Bucket")
	bucketOk    = false
)

func init() {
	viper.SetDefault("storage.Bucket", "sites")
	viper.SetDefault("storage.Endpoint", "127.0.0.1")
	viper.SetDefault("storage.AccessKey", "Q3AM3UQ867SPQQA43P2F")
	viper.SetDefault("storage.SecretKey", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG")
	viper.SetDefault("storage.Secure", true)
}

func makeBucket() (err error) {
	has, err := minioClient.BucketExists(bucketName)
	if err != nil {
		err = errors.NewMessageError(err, 7100, "查询存储桶错误")
		return
	}
	if !has {
		if err = minioClient.MakeBucket(bucketName, ""); err != nil {
			err = errors.NewMessageError(err, 7100, "创建存储桶错误")
			return
		}
	}
	bucketOk = true
	return
}

func initMinio() (err error) {
	bucketName = viper.GetString("storage.Bucket")
	endpoint := viper.GetString("storage.Endpoint")
	accessKeyID := viper.GetString("storage.AccessKey")
	secretAccessKey := viper.GetString("storage.SecretKey")
	secure := viper.GetBool("storage.Secure")

	if mdmClnt, err = madmin.New(endpoint, accessKeyID, secretAccessKey, secure); err != nil {
		err = errors.NewMessageError(err, 7100, "文件存储系统初始化失败")
		mdmClnt = nil
		return
	}
	mdmClnt.SetAppInfo("sites", "0.0.1")
	//si, err := mdmClnt.ServerInfo(context.Background())
	//if err != nil {
	//    fmt.Println(err)
	//    return
	//}
	//log.Debug("附件服务器", si)

	if minioClient, err = minio.New(endpoint, accessKeyID, secretAccessKey, secure); err != nil {
		err = errors.NewMessageError(err, 7100, "文件存储系统初始化失败")
		minioClient = nil
		return
	}
	minioClient.TraceOff()
	//minioClient.TraceOn(nil)
	minioClient.SetAppInfo("sites", "0.0.1")
	err = makeBucket()
	return
}

func tryMinio() (err error) {
	if minioClient != nil && bucketOk {
		return
	}
	err = initMinio()
	return
}

func PresignedGetInline(objectPath string) (presignedURL *url.URL, err error) {
	reqParams := make(url.Values)
	presignedURL, err = PresignedGet(objectPath, reqParams)
	return
}

func PresignedGetAttachment(objectPath string, filename string) (presignedURL *url.URL, err error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	presignedURL, err = PresignedGet(objectPath, reqParams)
	return
}

func PresignedGet(objectPath string, reqParams url.Values, expiresDuration ...time.Duration) (preURL *url.URL, err error) {
	expiresTime := expires * 2
	if len(expiresDuration) == 1 {
		expiresTime = expiresDuration[0]
	}
	if err = tryMinio(); err != nil {
		return
	}
	reqParams.Add("response-cache-control", "private")
	reqParams.Add("response-cache-control", "max-age=3600")
	preURL, err = minioClient.PresignedGetObject(bucketName, objectPath, expiresTime, reqParams)
	if err != nil {
		err = errors.NewMessageError(err, 7100, "presignedGet 失败")
		return
	}
	return
}

func PresignedPut(objectPath string) (presignedURL string, err error) {
	if err = tryMinio(); err != nil {
		return
	}

	preUrl, err := minioClient.PresignedPutObject(bucketName, objectPath, expires)
	err = errors.NewMessageError(err, 7100, "创建 PresignedPutObject 失败")
	if err != nil {
		return
	}
	presignedURL = preUrl.String()
	return
}
func PresignedPostPolicy(policy *minio.PostPolicy) (u *url.URL, formData map[string]string, err error) {
	if err = tryMinio(); err != nil {
		return
	}
	if err = policy.SetBucket(bucketName); err != nil {
		return
	}
	return minioClient.PresignedPostPolicy(policy)
}
func NewPostPolicy() *minio.PostPolicy {
	return minio.NewPostPolicy()
}

func ListIncompleteUploads() {
	doneCh := make(chan struct{})
	defer close(doneCh)
	//isRecursive := true // Recursively list everything at 'myprefix'
	multiPartObjectCh := minioClient.ListIncompleteUploads(bucketName, "", true, doneCh)
	fmt.Println("listIncompleteUploads", "开始清理")
	for multiPartObject := range multiPartObjectCh {
		if multiPartObject.Err != nil {
			fmt.Println("listIncompleteUploads", "错误", multiPartObject.Err)
			return
		}
		fmt.Println("listIncompleteUploads", "碎片", multiPartObject)
	}
	fmt.Println("listIncompleteUploads", "清理完成")
}
