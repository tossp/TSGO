package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio/pkg/madmin"

	"github.com/tossp/tsgo/pkg/errors"
	"github.com/tossp/tsgo/pkg/log"
	"github.com/tossp/tsgo/pkg/setting"
)

const expires = time.Hour

var (
	_minioClient *minio.Client
	mdmClnt      *madmin.AdminClient
	bucketName   = setting.GetString("storage.Bucket")
	confStr      = ""
	storageLock  sync.RWMutex

	errMinioClientNotReady = errors.New("存储未准备就绪")
)

func minioClient() *minio.Client {
	storageLock.RLock()
	defer storageLock.RUnlock()
	return _minioClient
}
func setMinioClient(m *minio.Client) {
	storageLock.Lock()
	defer storageLock.Unlock()
	_minioClient = m
}
func init() {
	setting.SetDefault("storage.Bucket", "sites")
	setting.SetDefault("storage.Endpoint", "127.0.0.1")
	setting.SetDefault("storage.AccessKey", "Q3AM3UQ867SPQQA43P2F")
	setting.SetDefault("storage.SecretKey", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG")
	setting.SetDefault("storage.Secure", true)
	setting.SetDefault("storage.Debug", false)
	_ = setting.Subscribe(autoInitMinio)
	log.WarnErr(autoInitMinio())
}
func makeConfStr() string {
	return setting.GetString("storage.Bucket") +
		setting.GetString("storage.Endpoint") +
		setting.GetString("storage.AccessKey") +
		setting.GetString("storage.SecretKey") +
		setting.GetString("storage.Bucket")
}
func autoInitMinio() (err error) {
	if makeConfStr() == confStr {
		return
	}
	return initMinio()
}

func makeBucket() (err error) {
	has, err := minioClient().BucketExists(context.Background(), bucketName)
	if err != nil {
		err = errors.NewMessageError(err, 7100, "查询存储桶错误")
		return
	}
	if !has {
		if err = minioClient().MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: "us-east-1"}); err != nil {
			err = errors.NewMessageError(err, 7100, "创建存储桶错误")
			return
		}
	}
	return
}

func initMinio() (err error) {
	confStr = makeConfStr()
	bucketName = setting.GetString("storage.Bucket")
	endpoint := setting.GetString("storage.Endpoint")
	accessKeyID := setting.GetString("storage.AccessKey")
	secretAccessKey := setting.GetString("storage.SecretKey")
	secure := setting.GetBool("storage.Secure")

	mc, err := minio.New(endpoint, &minio.Options{
		Region: "us-east-1",
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		err = errors.NewMessageError(err, 7100, "文件存储系统初始化失败")
		log.Warn("文件存储管理初始化失败", err)
		log.Warn("bucketName", bucketName)
		log.Warn("endpoint", endpoint)
		log.Warn("accessKeyID", accessKeyID)
		//log.Warn("secretAccessKey", secretAccessKey)
		return
	}
	if setting.GetBool("storage.Debug") {
		mc.TraceOn(nil)
	} else {
		mc.TraceOff()
	}
	mc.SetAppInfo(setting.AppName(), "0.0.1")

	setMinioClient(mc)

	// TODO(zh) 电院因WAF不合理配置导致检查失败，逻辑上可能会出现丢附件的情况，其他环境中必须开启！！！
	if !setting.GetBool("env.scdy") {
		if err = makeBucket(); err != nil {
			log.Warn("准备存储桶失败", err)
			setMinioClient(nil)
			return
		}
	}

	if mdmClnt, err = madmin.New(endpoint, accessKeyID, secretAccessKey, secure); err != nil {
		err = errors.NewMessageError(err, 7100, "文件存储系统初始化失败")
		log.Warn("文件存储管理初始化失败", err)
		log.Warn("bucketName", bucketName)
		log.Warn("endpoint", endpoint)
		log.Warn("accessKeyID", accessKeyID)
		mdmClnt = nil
		return
	}
	mdmClnt.SetAppInfo(setting.AppName(), "0.0.1")
	if setting.GetBool("storage.Debug") {
		mdmClnt.TraceOn(nil)
		si, fuck := mdmClnt.ServerInfo(context.Background())
		if fuck != nil {
			log.Warn("文件存储信息查询失败", fuck)
			return
		}
		log.Debug("附件服务器", si.Servers)
	} else {
		mdmClnt.TraceOff()
	}

	return
}

func PresignedGet(objectPath string, reqParams url.Values, expiresDuration ...time.Duration) (preURL *url.URL, err error) {
	expiresTime := expires * 2
	if len(expiresDuration) == 1 {
		expiresTime = expiresDuration[0]
	}
	if minioClient() == nil {
		err = errMinioClientNotReady
		return
	}
	reqParams.Add("response-cache-control", "private")
	reqParams.Add("response-cache-control", "max-age=3600")
	preURL, err = minioClient().PresignedGetObject(context.Background(), bucketName, objectPath, expiresTime, reqParams)
	if err != nil {
		err = errors.NewMessageError(err, 7100, "presignedGet 失败")
		return
	}
	return
}
func Get(objectPath string) (obj *minio.Object, err error) {
	if minioClient() == nil {
		err = errMinioClientNotReady
		return
	}
	obj, err = minioClient().GetObject(context.Background(), bucketName, objectPath, minio.GetObjectOptions{})
	if err != nil {
		err = errors.NewMessageError(err, 7100, "GetObject 失败")
		return
	}
	return
}
func PresignedPostPolicy(policy *minio.PostPolicy) (u *url.URL, formData map[string]string, err error) {
	if minioClient() == nil {
		err = errMinioClientNotReady
		return
	}
	if err = policy.SetBucket(bucketName); err != nil {
		return
	}
	return minioClient().PresignedPostPolicy(context.Background(), policy)
}
func NewPostPolicy() *minio.PostPolicy {
	return minio.NewPostPolicy()
}
func PutObject(objectName string, reader io.Reader, objectSize int64,
	opts minio.PutObjectOptions) (info minio.UploadInfo, err error) {
	info, err = minioClient().PutObject(context.Background(), bucketName, objectName, reader, objectSize, opts)
	return
}

func ListIncompleteUploads() {
	multiPartObjectCh := minioClient().ListIncompleteUploads(context.Background(), bucketName, "", true)
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
