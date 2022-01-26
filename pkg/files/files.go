package files

import (
	"context"
	"io/ioutil"
	"log"
	"fmt"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"gitlab.com/flattrack/flattrack/pkg/common"
)

type FileAccess struct {
	Client *minio.Client
	BucketName string
}

// Open ...
// open a Minio client
func Open(endpoint string, accessKey string, secretKey string, bucketName string, useSSL bool) (FileAccess, error) {
	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return FileAccess{}, err
	}
	return FileAccess{mc, bucketName}, err
}

func (f FileAccess) Init() error {
	if f.BucketName == "" {
		return fmt.Errorf("Error: cannot initialise a bucket, because no bucket name was provided")
	}
	buckets, err := f.Client.ListBuckets(context.TODO())
	if err != nil {
		return err
	}
	foundBucket := false
	for _, b := range buckets {
		if b.Name == f.BucketName {
			foundBucket = true
		}
	}
	if foundBucket == true {
		return nil
	}
	err = f.Client.MakeBucket(context.TODO(), f.BucketName, minio.MakeBucketOptions{})
	if err != nil {
		return err
	}
	return nil
}

// Get ...
// retrieves a given object
func (f FileAccess) Get(filePath string) (objectBytes []byte, objectInfo minio.ObjectInfo, err error) {
	object, err := f.Client.GetObject(context.TODO(), common.GetAppMinioBucket(), filePath, minio.GetObjectOptions{})
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

// Put ...
// uploads a file
func (f FileAccess) Put() {}
