package files

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"gitlab.com/flattrack/flattrack/internal/common"
)

// FileAccess to access an S3 compatible backend
type FileAccess struct {
	Client     *minio.Client
	BucketName string
	Prefix     string
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
	return FileAccess{Client: mc, BucketName: bucketName}, err
}

// Init to initialise a bucket
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
	if foundBucket {
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
func (f FileAccess) Get(name string) (objectBytes []byte, objectInfo minio.ObjectInfo, err error) {
	fileName := fmt.Sprintf("%v-%v", f.Prefix, name)
	object, err := f.Client.GetObject(context.TODO(), common.GetAppMinioBucket(), fileName, minio.GetObjectOptions{})
	if err != nil {
		slog.Error("Failed to get object", "error", err)
		return []byte{}, minio.ObjectInfo{}, err
	}
	defer func() {
		if err := object.Close(); err != nil {
			slog.Error("Failed to close object", "file", fileName, "error", err)
		}
	}()
	objectInfo, err = object.Stat()
	if err != nil {
		slog.Error("Failed to stat object", "error", err)
		return []byte{}, minio.ObjectInfo{}, err
	}
	objectBytes, err = io.ReadAll(object)
	if err != nil {
		slog.Error("Failed to read object", "error", err)
		return []byte{}, minio.ObjectInfo{}, err
	}
	return objectBytes, objectInfo, err
}

// Put ...
// uploads a file
func (f FileAccess) Put(name string, data []byte) error {
	fileName := fmt.Sprintf("%v-%v", f.Prefix, name)
	reader := bytes.NewReader(data)
	info, err := f.Client.PutObject(context.TODO(), f.BucketName, fileName, reader, int64(reader.Len()), minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	slog.Info("Successfully uploaded file into bucket", "file", info.Key)
	return nil
}

// Delete ...
// deletes a file
func (f FileAccess) Delete(name string) error {
	fileName := fmt.Sprintf("%v-%v", f.Prefix, name)
	err := f.Client.RemoveObject(context.TODO(), f.BucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	slog.Info("Successfully deleted file into bucket", "file", fileName)
	return nil
}
