package adapter

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/npavlov/go-password-manager/internal/server/service/file"
)

type minioAdapter struct {
	client *minio.Client
}

func NewMinioAdapter(client *minio.Client) file.S3Storage {
	return &minioAdapter{client: client}
}

func (m *minioAdapter) PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	return m.client.PutObject(ctx, bucketName, objectName, reader, objectSize, opts)
}

func (m *minioAdapter) GetObject(ctx context.Context, bucketName string, objectName string, opts minio.GetObjectOptions) (io.ReadCloser, error) {
	return m.client.GetObject(ctx, bucketName, objectName, opts)
}

func (m *minioAdapter) RemoveObject(ctx context.Context, bucketName string, objectName string, opts minio.RemoveObjectOptions) error {
	return m.client.RemoveObject(ctx, bucketName, objectName, opts)
}
