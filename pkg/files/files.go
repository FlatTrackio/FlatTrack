package files

import (
	"context"
	"io/ioutil"
	"log"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"gitlab.com/flattrack/flattrack/pkg/common"
)

type Client struct {
	minio.Client
}

// Open ...
// open a Minio client
func Open(endpoint string, accessKey string, secretKey string, useSSL bool) (*minio.Client, error) {
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
}

func Init(mc *minio.Client, bucketName string) error {
	buckets, err := mc.ListBuckets(context.TODO())
	if err != nil {
		return err
	}
	foundBucket := false
	for _, b := range buckets {
		if b.Name == bucketName {
			foundBucket = true
		}
	}
	if foundBucket == true {
		return nil
	}
	err = mc.MakeBucket(context.TODO(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return err
	}
	return nil
}

// Get ...
// retrieves a given object
func Get(minioClient *minio.Client, filePath string) (objectBytes []byte, objectInfo minio.ObjectInfo, err error) {
	object, err := minioClient.GetObject(context.TODO(), common.GetAppMinioBucket(), filePath, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("%#v\n", err)
		return []byte{}, minio.ObjectInfo{}, err
	}
	defer object.Close()
	objectInfo, err = object.Stat()
	if err != nil {
		log.Printf("%#v\n", err)
		return []byte{}, minio.ObjectInfo{}, err
	}
	objectBytes, err = ioutil.ReadAll(object)
	if err != nil {
		log.Printf("%#v\n", err)
		return []byte{}, minio.ObjectInfo{}, err
	}
	return objectBytes, objectInfo, err
}
