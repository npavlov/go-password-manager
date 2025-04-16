//nolint:wrapcheck,lll,err113,dogsled
package file_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"io"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/npavlov/go-password-manager/gen/proto/file"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/file"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	generalutils "github.com/npavlov/go-password-manager/internal/utils"
)

func setupFileService(t *testing.T) (*file.Service, *testutils.MockDBStorage, *MockS3Storage, context.Context, string) {
	t.Helper()

	logger := zerolog.New(nil)
	masterKey, _ := utils.GenerateRandomKey()

	cfg := &config.Config{
		SecuredMasterKey: generalutils.NewString(masterKey),
		Bucket:           "test-bucket",
	}

	storage := testutils.SetupMockUserStorage(masterKey)
	mockS3 := &MockS3Storage{}
	svc := file.NewFileService(&logger, storage, cfg, mockS3)
	encryptionKey, _ := utils.GenerateRandomKey()

	encryptionKeyEncrypted, _ := utils.Encrypt(encryptionKey, masterKey)

	// Create test user
	testUser := db.User{
		ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Username:      "tester",
		Email:         "test@example.com",
		Password:      "hashed-password",
		EncryptionKey: encryptionKeyEncrypted,
	}
	storage.AddTestUser(testUser)

	// Inject user ID and encryption key into context
	ctx := testutils.InjectUserToContext(t.Context(), testUser.ID.String())

	return svc, storage, mockS3, ctx, masterKey
}

type MockS3Storage struct {
	PutObjectFunc    func(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error)
	GetObjectFunc    func(ctx context.Context, bucketName string, objectName string, opts minio.GetObjectOptions) (io.ReadCloser, error)
	RemoveObjectFunc func(ctx context.Context, bucketName string, objectName string, opts minio.RemoveObjectOptions) error
}

func (m *MockS3Storage) PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	if m.PutObjectFunc != nil {
		return m.PutObjectFunc(ctx, bucketName, objectName, reader, objectSize, opts)
	}

	return minio.UploadInfo{}, nil
}

func (m *MockS3Storage) GetObject(ctx context.Context, bucketName string, objectName string, opts minio.GetObjectOptions) (io.ReadCloser, error) {
	if m.GetObjectFunc != nil {
		return m.GetObjectFunc(ctx, bucketName, objectName, opts)
	}

	return &minio.Object{}, nil
}

func (m *MockS3Storage) RemoveObject(ctx context.Context, bucketName string, objectName string, opts minio.RemoveObjectOptions) error {
	if m.RemoveObjectFunc != nil {
		return m.RemoveObjectFunc(ctx, bucketName, objectName, opts)
	}

	return nil
}

type MockUploadStream struct {
	RecvFunc         func() (*pb.UploadFileRequest, error)
	SendAndCloseFunc func(*pb.UploadFileResponse) error
	ContextFunc      func() context.Context
}

func (m *MockUploadStream) SetHeader(md metadata.MD) error {
	panic("implement me")
}

func (m *MockUploadStream) SendHeader(md metadata.MD) error {
	panic("implement me")
}

func (m *MockUploadStream) SetTrailer(md metadata.MD) {
	panic("implement me")
}

func (m *MockUploadStream) SendMsg(m2 any) error {
	panic("implement me")
}

func (m *MockUploadStream) RecvMsg(m1 any) error {
	panic("implement me")
}

func (m *MockUploadStream) Recv() (*pb.UploadFileRequest, error) {
	if m.RecvFunc != nil {
		return m.RecvFunc()
	}

	return nil, io.EOF
}

func (m *MockUploadStream) SendAndClose(resp *pb.UploadFileResponse) error {
	if m.SendAndCloseFunc != nil {
		return m.SendAndCloseFunc(resp)
	}

	return nil
}

func (m *MockUploadStream) Context() context.Context {
	if m.ContextFunc != nil {
		return m.ContextFunc()
	}

	return context.Background()
}

func TestUploadFile_Success(t *testing.T) {
	t.Parallel()

	svc, storage, mockS3, ctx, _ := setupFileService(t)

	mockStream := &MockUploadStream{
		RecvFunc: func() func() (*pb.UploadFileRequest, error) {
			calls := 0

			return func() (*pb.UploadFileRequest, error) {
				if calls == 0 {
					calls++

					return &pb.UploadFileRequest{
						Filename: "test.txt",
						Data:     []byte("test data"),
					}, nil
				}

				return nil, io.EOF
			}
		}(),
		SendAndCloseFunc: func(*pb.UploadFileResponse) error {
			return nil
		},
		ContextFunc: func() context.Context {
			return ctx
		},
	}

	mockS3.PutObjectFunc = func(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
		return minio.UploadInfo{}, nil
	}

	err := svc.UploadFile(mockStream)
	require.NoError(t, err)

	// Verify the binary was stored
	userID := testutils.GetUserIDFromContext(ctx)
	binaries, err := storage.GetBinaries(t.Context(), userID)
	require.NoError(t, err)
	require.Len(t, binaries, 1)
	require.Equal(t, "test.txt", binaries[0].FileName)
}

func TestUploadFile_InvalidMetadata(t *testing.T) {
	t.Parallel()

	svc, _, _, ctx, _ := setupFileService(t)

	mockStream := &MockUploadStream{
		RecvFunc: func() (*pb.UploadFileRequest, error) {
			return &pb.UploadFileRequest{}, nil // Missing required filename
		},
		ContextFunc: func() context.Context {
			return ctx
		},
	}

	err := svc.UploadFile(mockStream)
	require.Error(t, err)
}

func TestUploadFile_S3UploadFailure(t *testing.T) {
	t.Parallel()

	svc, _, mockS3, ctx, _ := setupFileService(t)

	mockStream := &MockUploadStream{
		RecvFunc: func() (*pb.UploadFileRequest, error) {
			return &pb.UploadFileRequest{
				Filename: "test.txt",
				Data:     []byte("test data"),
			}, nil
		},
		ContextFunc: func() context.Context {
			return ctx
		},
	}

	mockS3.PutObjectFunc = func(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
		return minio.UploadInfo{}, errors.New("S3 upload failed")
	}

	err := svc.UploadFile(mockStream)
	require.Error(t, err)
}

func TestDownloadFile_Success(t *testing.T) {
	t.Parallel()

	svc, storage, mockS3, ctx, masterKey := setupFileService(t)

	// First store a file to download
	userID := testutils.GetUserIDFromContext(ctx)
	binary, err := storage.StoreBinary(ctx, db.StoreBinaryEntryParams{
		UserID:   pgtype.UUID{Bytes: uuid.MustParse(userID), Valid: true},
		FileName: "test.txt",
		FileSize: 123,
		FileUrl:  userID + "-test.txt",
	})
	require.NoError(t, err)

	user, err := storage.GetUserByID(ctx, pgtype.UUID{Bytes: uuid.MustParse(userID), Valid: true})
	require.NoError(t, err)

	userKey, err := utils.Decrypt(user.EncryptionKey, masterKey)
	require.NoError(t, err)

	// Generate test data - 3 blocks of 1024 random bytes each
	//nolint:prealloc
	var plaintextBlocks [][]byte
	for range 3 {
		block := make([]byte, 1024)
		_, err := rand.Read(block)
		require.NoError(t, err)
		plaintextBlocks = append(plaintextBlocks, block)
	}

	// Encrypt each block with user's key
	var encryptedData []byte

	// Create an encryptor pipe to encrypt the data exactly like the service would
	pr, pw := io.Pipe()
	encryptor, err := utils.NewEncryptor(pw, userKey)
	require.NoError(t, err)

	// Write data to encryptor in a goroutine
	go func() {
		defer pw.Close()
		for _, block := range plaintextBlocks {
			_, err := encryptor.Write(block)
			if err != nil {
				pw.CloseWithError(err)

				return
			}
		}
	}()

	// Read all encrypted data
	encryptedData, err = io.ReadAll(pr)
	require.NoError(t, err)

	// Create a mock reader with encrypted test data
	mockReader := io.NopCloser(bytes.NewReader(encryptedData))

	// Mock the S3 GetObject to return our test reader
	mockS3.GetObjectFunc = func(ctx context.Context, bucketName string, objectName string, opts minio.GetObjectOptions) (io.ReadCloser, error) {
		return mockReader, nil
	}

	// Track how many times Send is called and what data is received
	var sendCount int
	var receivedData []byte

	mockStream := &MockDownloadStream{
		SendFunc: func(resp *pb.DownloadFileResponse) error {
			sendCount++
			receivedData = append(receivedData, resp.GetData()...)

			return nil
		},
		ContextFunc: func() context.Context {
			return ctx
		},
	}

	err = svc.DownloadFile(&pb.DownloadFileRequest{FileId: binary.ID.String()}, mockStream)
	require.NoError(t, err)

	// Verify the data was sent correctly
	require.Positive(t, sendCount, "Send should be called at least once")

	// Verify we received the decrypted original data
	var decryptedReceived []byte
	for _, block := range plaintextBlocks {
		decryptedReceived = append(decryptedReceived, block...)
	}
	require.Equal(t, decryptedReceived, receivedData, "Received decrypted data should match original plaintext")
}

func TestDownloadFile_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, _, ctx, _ := setupFileService(t)

	mockStream := &MockDownloadStream{
		ContextFunc: func() context.Context {
			return ctx
		},
	}

	err := svc.DownloadFile(&pb.DownloadFileRequest{FileId: uuid.NewString()}, mockStream)
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestDeleteFile_Success(t *testing.T) {
	t.Parallel()

	svc, storage, mockS3, ctx, _ := setupFileService(t)

	// First store a file to delete
	userID := testutils.GetUserIDFromContext(ctx)
	binary, err := storage.StoreBinary(ctx, db.StoreBinaryEntryParams{
		UserID:   pgtype.UUID{Bytes: uuid.MustParse(userID), Valid: true},
		FileName: "test.txt",
		FileSize: 123,
		FileUrl:  userID + "-test.txt",
	})
	require.NoError(t, err)

	mockS3.RemoveObjectFunc = func(ctx context.Context, bucketName string, objectName string, opts minio.RemoveObjectOptions) error {
		return nil
	}

	resp, err := svc.DeleteFile(ctx, &pb.DeleteFileRequest{FileId: binary.ID.String()})
	require.NoError(t, err)
	require.True(t, resp.GetOk())

	// Verify file was deleted
	_, err = storage.GetBinary(ctx, binary.ID.String(), pgtype.UUID{Bytes: uuid.MustParse(userID), Valid: true})
	require.Error(t, err)
}

func TestGetFile_Success(t *testing.T) {
	t.Parallel()

	svc, storage, _, ctx, _ := setupFileService(t)

	// First store a file to get
	userID := testutils.GetUserIDFromContext(ctx)
	binary, err := storage.StoreBinary(ctx, db.StoreBinaryEntryParams{
		UserID:   pgtype.UUID{Bytes: uuid.MustParse(userID), Valid: true},
		FileName: "test.txt",
		FileSize: 123,
		FileUrl:  userID + "-test.txt",
	})
	require.NoError(t, err)

	resp, err := svc.GetFile(ctx, &pb.GetFileRequest{FileId: binary.ID.String()})
	require.NoError(t, err)
	require.Equal(t, binary.ID.String(), resp.GetFile().GetId())
	require.Equal(t, "test.txt", resp.GetFile().GetFileName())
	require.Equal(t, int64(123), resp.GetFile().GetFileSize())
	require.Equal(t, userID+"-test.txt", resp.GetFile().GetFileUrl())
}

func TestGetFiles_Success(t *testing.T) {
	t.Parallel()

	svc, storage, _, ctx, _ := setupFileService(t)

	// First store some files
	userID := testutils.GetUserIDFromContext(ctx)
	_, err := storage.StoreBinary(ctx, db.StoreBinaryEntryParams{
		UserID:   pgtype.UUID{Bytes: uuid.MustParse(userID), Valid: true},
		FileName: "test1.txt",
		FileSize: 123,
		FileUrl:  userID + "-test1.txt",
	})
	require.NoError(t, err)

	_, err = storage.StoreBinary(ctx, db.StoreBinaryEntryParams{
		UserID:   pgtype.UUID{Bytes: uuid.MustParse(userID), Valid: true},
		FileName: "test2.txt",
		FileSize: 456,
		FileUrl:  userID + "-test2.txt",
	})
	require.NoError(t, err)

	_, err = svc.GetFiles(ctx, &pb.GetFilesRequest{})
	require.NoError(t, err)
}

func TestUploadFile_EmptyFile(t *testing.T) {
	t.Parallel()

	svc, _, mockS3, ctx, _ := setupFileService(t)

	mockStream := &MockUploadStream{
		RecvFunc: func() (*pb.UploadFileRequest, error) {
			return &pb.UploadFileRequest{
				Filename: "empty.txt",
				Data:     []byte{},
			}, io.EOF
		},
		ContextFunc: func() context.Context {
			return ctx
		},
	}

	mockS3.PutObjectFunc = func(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
		return minio.UploadInfo{}, nil
	}

	err := svc.UploadFile(mockStream)
	require.Error(t, err)
}

func TestUploadFile_LargeFile(t *testing.T) {
	t.Parallel()

	svc, _, mockS3, ctx, _ := setupFileService(t)

	// Use a channel to safely communicate the received size
	sizeChan := make(chan int64, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	// Simulate a 5MB file in chunks
	chunkCount := 0
	mockStream := &MockUploadStream{
		RecvFunc: func() (*pb.UploadFileRequest, error) {
			if chunkCount == 0 {
				chunkCount++

				return &pb.UploadFileRequest{
					Filename: "large.bin",
				}, nil
			}
			if chunkCount < 50 { // 50 chunks of 100KB each = 5MB
				chunkCount++
				data := make([]byte, 100*1024) // 100KB chunk

				return &pb.UploadFileRequest{
					Data: data,
				}, nil
			}

			return nil, io.EOF
		},
		ContextFunc: func() context.Context {
			return ctx
		},
	}

	mockS3.PutObjectFunc = func(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
		defer wg.Done()
		cursor, err := io.Copy(io.Discard, reader)
		if err != nil {
			sizeChan <- 0

			return minio.UploadInfo{}, err
		}
		sizeChan <- cursor

		return minio.UploadInfo{}, nil
	}

	err := svc.UploadFile(mockStream)
	require.NoError(t, err)

	// Wait for the S3 upload to complete
	wg.Wait()
	close(sizeChan)

	receivedSize := <-sizeChan
	assert.Positive(t, receivedSize)
}

func TestUploadFile_StorageFailureAfterS3Upload(t *testing.T) {
	t.Parallel()

	svc, storage, _, ctx, _ := setupFileService(t)

	// Make storage fail after S3 upload succeeds
	storage.CallError = errors.New("database failure")

	mockStream := &MockUploadStream{
		RecvFunc: func() (*pb.UploadFileRequest, error) {
			return &pb.UploadFileRequest{
				Filename: "test.txt",
				Data:     []byte("test data"),
			}, io.EOF
		},
		ContextFunc: func() context.Context {
			return ctx
		},
	}

	err := svc.UploadFile(mockStream)
	require.Error(t, err)
}

func TestUploadFile_EncryptionError(t *testing.T) {
	t.Parallel()

	svc, storage, _, ctx, _ := setupFileService(t)

	// Corrupt the user's encryption key
	userID := testutils.GetUserIDFromContext(ctx)
	userIDPG := pgtype.UUID{
		Bytes: uuid.MustParse(userID),
		Valid: true,
	}
	user := storage.UsersByID[userIDPG]
	user.EncryptionKey = "invalid-key"
	storage.UsersByID[userIDPG] = user

	mockStream := &MockUploadStream{
		RecvFunc: func() (*pb.UploadFileRequest, error) {
			return &pb.UploadFileRequest{
				Filename: "test.txt",
				Data:     []byte("test data"),
			}, nil
		},
		ContextFunc: func() context.Context {
			return ctx
		},
	}

	err := svc.UploadFile(mockStream)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error getting user id")
}

func TestDownloadFile_DecryptionError(t *testing.T) {
	t.Parallel()

	svc, storage, mockS3, ctx, _ := setupFileService(t)

	// Store a file with invalid encryption
	userID := testutils.GetUserIDFromContext(ctx)
	binary, err := storage.StoreBinary(ctx, db.StoreBinaryEntryParams{
		UserID:   pgtype.UUID{Bytes: uuid.MustParse(userID), Valid: true},
		FileName: "test.txt",
		FileSize: 123,
		FileUrl:  userID + "-test.txt",
	})
	require.NoError(t, err)

	// Generate test data - 3 blocks of 1024 random bytes each
	var plaintextBlocks []byte
	for range 3 {
		block := make([]byte, 1024)
		_, err := rand.Read(block)
		require.NoError(t, err)
		plaintextBlocks = append(plaintextBlocks, block...)
	}

	// Create a mock reader with encrypted test data
	mockReader := io.NopCloser(bytes.NewReader(plaintextBlocks))

	// Mock the S3 GetObject to return our test reader
	mockS3.GetObjectFunc = func(ctx context.Context, bucketName string, objectName string, opts minio.GetObjectOptions) (io.ReadCloser, error) {
		return mockReader, nil
	}

	// Track how many times Send is called and what data is received
	var sendCount int
	var receivedData []byte

	mockStream := &MockDownloadStream{
		SendFunc: func(resp *pb.DownloadFileResponse) error {
			sendCount++
			receivedData = append(receivedData, resp.GetData()...)

			return nil
		},
		ContextFunc: func() context.Context {
			return ctx
		},
	}

	err = svc.DownloadFile(&pb.DownloadFileRequest{FileId: binary.ID.String()}, mockStream)
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}

func TestDownloadFile_UnauthorizedAccess(t *testing.T) {
	t.Parallel()

	svc, storage, _, ctx, _ := setupFileService(t)

	// Store a file with a different user
	otherUserID := uuid.New()
	binary, err := storage.StoreBinary(t.Context(), db.StoreBinaryEntryParams{
		UserID:   pgtype.UUID{Bytes: otherUserID, Valid: true},
		FileName: "test.txt",
		FileSize: 123,
		FileUrl:  otherUserID.String() + "-test.txt",
	})
	require.NoError(t, err)

	mockStream := &MockDownloadStream{
		ContextFunc: func() context.Context {
			return ctx
		},
	}

	err = svc.DownloadFile(&pb.DownloadFileRequest{FileId: binary.ID.String()}, mockStream)
	require.Error(t, err)
}

func TestDeleteFile_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, _, ctx, _ := setupFileService(t)

	resp, err := svc.DeleteFile(ctx, &pb.DeleteFileRequest{FileId: uuid.NewString()})
	require.Error(t, err)
	require.False(t, resp.GetOk())
}

func TestDeleteFile_Unauthorized(t *testing.T) {
	t.Parallel()

	svc, storage, _, ctx, _ := setupFileService(t)

	// Store a file with a different user
	otherUserID := uuid.New()
	binary, err := storage.StoreBinary(t.Context(), db.StoreBinaryEntryParams{
		UserID:   pgtype.UUID{Bytes: otherUserID, Valid: true},
		FileName: "test.txt",
		FileSize: 123,
		FileUrl:  otherUserID.String() + "-test.txt",
	})
	require.NoError(t, err)

	resp, err := svc.DeleteFile(ctx, &pb.DeleteFileRequest{FileId: binary.ID.String()})
	require.Error(t, err)
	require.False(t, resp.GetOk())
}

func TestDeleteFile_S3DeleteFailure(t *testing.T) {
	t.Parallel()

	svc, storage, _, ctx, _ := setupFileService(t)

	// First store a file to delete
	userID := testutils.GetUserIDFromContext(ctx)
	binary, err := storage.StoreBinary(ctx, db.StoreBinaryEntryParams{
		UserID:   pgtype.UUID{Bytes: uuid.MustParse(userID), Valid: true},
		FileName: "test.txt",
		FileSize: 123,
		FileUrl:  userID + "-test.txt",
	})
	require.NoError(t, err)

	storage.CallError = errors.New("error deleting file")

	resp, err := svc.DeleteFile(ctx, &pb.DeleteFileRequest{FileId: binary.ID.String()})
	require.Error(t, err)
	require.False(t, resp.GetOk())
}

func TestGetFile_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, _, ctx, _ := setupFileService(t)

	resp, err := svc.GetFile(ctx, &pb.GetFileRequest{FileId: uuid.NewString()})
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestGetFile_Unauthorized(t *testing.T) {
	t.Parallel()

	svc, storage, _, ctx, _ := setupFileService(t)

	// Store a file with a different user
	otherUserID := uuid.New()
	binary, err := storage.StoreBinary(t.Context(), db.StoreBinaryEntryParams{
		UserID:   pgtype.UUID{Bytes: otherUserID, Valid: true},
		FileName: "test.txt",
		FileSize: 123,
		FileUrl:  otherUserID.String() + "-test.txt",
	})
	require.NoError(t, err)

	resp, err := svc.GetFile(ctx, &pb.GetFileRequest{FileId: binary.ID.String()})
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestGetFiles_EmptyResult(t *testing.T) {
	t.Parallel()

	svc, _, _, ctx, _ := setupFileService(t)

	resp, err := svc.GetFiles(ctx, &pb.GetFilesRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestGetFiles_WithPagination(t *testing.T) {
	t.Parallel()

	svc, storage, _, ctx, _ := setupFileService(t)

	// Store multiple files
	userID := testutils.GetUserIDFromContext(ctx)
	for range 10 {
		_, err := storage.StoreBinary(ctx, db.StoreBinaryEntryParams{
			UserID:   pgtype.UUID{Bytes: uuid.MustParse(userID), Valid: true},
			FileName: "test.txt",
			FileSize: 123,
			FileUrl:  userID + "-test.txt",
		})
		require.NoError(t, err)
	}

	resp, err := svc.GetFiles(ctx, &pb.GetFilesRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
}

type MockDownloadStream struct {
	SendFunc    func(*pb.DownloadFileResponse) error
	ContextFunc func() context.Context
}

func (m *MockDownloadStream) Send(resp *pb.DownloadFileResponse) error {
	if m.SendFunc != nil {
		return m.SendFunc(resp)
	}

	return nil
}

func (m *MockDownloadStream) Context() context.Context {
	if m.ContextFunc != nil {
		return m.ContextFunc()
	}

	return context.Background()
}

// The following are required by the gRPC stream interface but not used in your test.
func (m *MockDownloadStream) SetHeader(md metadata.MD) error {
	panic("implement me")
}

func (m *MockDownloadStream) SendHeader(md metadata.MD) error {
	panic("implement me")
}

func (m *MockDownloadStream) SetTrailer(md metadata.MD) {
	panic("implement me")
}

func (m *MockDownloadStream) SendMsg(msg any) error {
	panic("implement me")
}

func (m *MockDownloadStream) RecvMsg(msg any) error {
	panic("implement me")
}
