package items

import (
	"context"

	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	"github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// Client GRPCClient handles communication with the gRPC server
type Client struct {
	conn         *grpc.ClientConn
	client       pb.ItemServiceClient
	tokenManager *auth.TokenManager
	log          *zerolog.Logger
}

// NewItemsClient  creates a new gRPC connection
func NewItemsClient(conn *grpc.ClientConn, tokenManager *auth.TokenManager, log *zerolog.Logger) *Client {
	return &Client{
		conn:         conn,
		client:       pb.NewItemServiceClient(conn),
		tokenManager: tokenManager,
		log:          log,
	}
}

// GetItems sends a register request to the server
func (as *Client) GetItems(ctx context.Context, page, pageSize int32) ([]*pb.ItemData, int32, error) {
	resp, err := as.client.GetItems(ctx, &pb.GetItemsRequest{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, 0, errors.Wrapf(err, "GetItems failed, page=%d, pageSize=%d", page, pageSize)
	}

	return resp.Items, resp.TotalCount, nil
}
