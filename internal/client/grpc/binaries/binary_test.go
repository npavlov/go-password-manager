//nolint:wrapcheck,err113,lll,exhaustruct,forcetypeassert,ireturn,revive
package binary_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	pb "github.com/npavlov/go-password-manager/gen/proto/file"
	binary "github.com/npavlov/go-password-manager/internal/client/grpc/binaries"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

// Mocks

type MockFileServiceClient struct {
	mock.Mock
}

type MockUploadStream struct {
	mock.Mock
	grpc.ClientStream
}

type MockDownloadStream struct {
	mock.Mock
	grpc.ClientStream
}

type MockTokenManager struct {
	mock.Mock
}

func (m *MockFileServiceClient) GetFilesV1(ctx context.Context, _ *pb.GetFilesV1Request, _ ...grpc.CallOption) (*pb.GetFilesV1Response, error) {
	args := m.Called(ctx)

	return args.Get(0).(*pb.GetFilesV1Response), args.Error(1)
}

func (m *MockFileServiceClient) UploadFileV1(ctx context.Context, _ ...grpc.CallOption) (pb.FileService_UploadFileV1Client, error) {
	args := m.Called(ctx)

	// Safely handle nil to avoid type assertion panic
	stream, ok := args.Get(0).(pb.FileService_UploadFileV1Client)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return stream, args.Error(1)
}

func (m *MockFileServiceClient) DownloadFileV1(ctx context.Context, req *pb.DownloadFileV1Request, _ ...grpc.CallOption) (pb.FileService_DownloadFileV1Client, error) {
	args := m.Called(ctx, req)

	return args.Get(0).(pb.FileService_DownloadFileV1Client), args.Error(1)
}

func (m *MockFileServiceClient) DeleteFileV1(ctx context.Context, req *pb.DeleteFileV1Request, _ ...grpc.CallOption) (*pb.DeleteFileV1Response, error) {
	args := m.Called(ctx, req)

	return args.Get(0).(*pb.DeleteFileV1Response), args.Error(1)
}

func (m *MockFileServiceClient) GetFileV1(ctx context.Context, req *pb.GetFileV1Request, opts ...grpc.CallOption) (*pb.GetFileV1Response, error) {
	args := m.Called(ctx, req)

	return args.Get(0).(*pb.GetFileV1Response), args.Error(1)
}

func (s *MockUploadStream) Send(req *pb.UploadFileV1Request) error {
	args := s.Called(req)

	return args.Error(0)
}

func (s *MockUploadStream) CloseAndRecv() (*pb.UploadFileV1Response, error) {
	args := s.Called()

	return args.Get(0).(*pb.UploadFileV1Response), args.Error(1)
}

func (s *MockDownloadStream) Recv() (*pb.DownloadFileV1Response, error) {
	args := s.Called()

	// Safely handle nil to avoid type assertion panic
	stream, ok := args.Get(0).(*pb.DownloadFileV1Response)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return stream, args.Error(1)
}

func TestUploadFile_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockFileServiceClient)
	mockStream := new(MockUploadStream)
	logger := zerolog.Nop()

	reader := bytes.NewReader([]byte("hello world"))

	mockClient.On("UploadFileV1", mock.Anything).Return(mockStream, nil)
	mockStream.On("Send", mock.MatchedBy(func(_ *pb.UploadFileV1Request) bool {
		return true
	})).Return(nil).Times(2)
	mockStream.On("CloseAndRecv").Return(&pb.UploadFileV1Response{FileId: "file123"}, nil)

	client := &binary.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	fileID, err := client.UploadFile(t.Context(), "hello.txt", reader)
	require.NoError(t, err)
	assert.Equal(t, "file123", fileID)
}

func TestUploadFile_StreamError(t *testing.T) {
	t.Parallel()

	mockClient := new(MockFileServiceClient)
	logger := zerolog.Nop()

	var nilUploadStream pb.FileService_UploadFileV1Client = nil
	mockClient.On("UploadFileV1", mock.Anything).Return(nilUploadStream, errors.New("stream failed"))

	client := &binary.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.UploadFile(t.Context(), "file.txt", bytes.NewReader([]byte("data")))
	require.Error(t, err)
}

func TestDownloadFile_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockFileServiceClient)
	mockStream := new(MockDownloadStream)
	writer := new(bytes.Buffer)
	logger := zerolog.Nop()

	mockClient.On("DownloadFileV1", mock.Anything, &pb.DownloadFileV1Request{FileId: "file123"}).Return(mockStream, nil)
	mockStream.On("Recv").Return(&pb.DownloadFileV1Response{Data: []byte("chunk1")}, nil).Once()

	var nilResp *pb.DownloadFileV1Response
	mockStream.On("Recv").Return(nilResp, io.EOF).Once()

	client := &binary.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	err := client.DownloadFile(t.Context(), "file123", writer)
	require.NoError(t, err)
	assert.Equal(t, "chunk1", writer.String())
}

func TestDownloadFile_RecvError(t *testing.T) {
	t.Parallel()

	mockClient := new(MockFileServiceClient)
	mockStream := new(MockDownloadStream)
	logger := zerolog.Nop()

	mockClient.On("DownloadFileV1", mock.Anything, &pb.DownloadFileV1Request{FileId: "file123"}).Return(mockStream, nil)
	mockStream.On("Recv").Return(nil, errors.New("recv failed"))

	client := &binary.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	err := client.DownloadFile(t.Context(), "file123", new(bytes.Buffer))
	require.Error(t, err)
}

func TestDeleteFile_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockFileServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeleteFileV1", mock.Anything, &pb.DeleteFileV1Request{FileId: "file123"}).
		Return(&pb.DeleteFileV1Response{Ok: true}, nil)

	client := &binary.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	ok, err := client.DeleteFile(t.Context(), "file123")
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestGetFile_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockFileServiceClient)
	logger := zerolog.Nop()

	meta := &pb.FileMeta{Id: "file123", FileName: "myfile.txt"}

	mockClient.On("GetFileV1", mock.Anything, &pb.GetFileV1Request{FileId: "file123"}).
		Return(&pb.GetFileV1Response{File: meta}, nil)

	client := &binary.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	file, err := client.GetFile(t.Context(), "file123")
	require.NoError(t, err)
	assert.Equal(t, meta, file)
}

func TestNewBinaryClient(t *testing.T) {
	t.Parallel()

	tm := new(testutils.MockTokenManager)
	logger := zerolog.Nop()

	client := binary.NewBinaryClient(nil, tm, &logger)

	assert.NotNil(t, client)
}
