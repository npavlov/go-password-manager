//nolint:wrapcheck,lll
package adapter

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioAdapter struct {
	client *minio.Client
}

func NewMinioAdapter(client *minio.Client) *MinioAdapter {
	return &MinioAdapter{client: client}
}

func (m *MinioAdapter) PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	return m.client.PutObject(ctx, bucketName, objectName, reader, objectSize, opts)
}

func (m *MinioAdapter) GetObject(ctx context.Context, bucketName string, objectName string, opts minio.GetObjectOptions) (io.ReadCloser, error) {
	return m.client.GetObject(ctx, bucketName, objectName, opts)
}

func (m *MinioAdapter) RemoveObject(ctx context.Context, bucketName string, objectName string, opts minio.RemoveObjectOptions) error {
	return m.client.RemoveObject(ctx, bucketName, objectName, opts)
}
