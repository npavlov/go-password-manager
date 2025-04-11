package metainfo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	pb "github.com/npavlov/go-password-manager/gen/proto/metadata"
	"github.com/npavlov/go-password-manager/internal/client/auth"
)

// Client GRPCClient handles communication with the gRPC server.
type Client struct {
	conn         *grpc.ClientConn
	client       pb.MetadataServiceClient
	tokenManager *auth.TokenManager
	log          *zerolog.Logger
}

// NewMetainfoClient  creates a new gRPC connection.
func NewMetainfoClient(conn *grpc.ClientConn, tokenManager *auth.TokenManager, log *zerolog.Logger) *Client {
	return &Client{
		conn:         conn,
		client:       pb.NewMetadataServiceClient(conn),
		tokenManager: tokenManager,
		log:          log,
	}
}

// GetMetainfo sends a register request to the server.
func (as *Client) GetMetainfo(ctx context.Context, id string) (map[string]string, error) {
	resp, err := as.client.GetMetaInfo(ctx, &pb.GetMetaInfoRequest{
		ItemId: id,
	})
	if err != nil {
		as.log.Error().Err(err).Msg("error getting metainfo")

		return nil, errors.Wrap(err, "error getting metainfo")
	}

	return resp.GetMetadata(), nil
}

// SetMetainfo sets meta information for the record.
func (as *Client) SetMetainfo(ctx context.Context, id string, meta map[string]string) (bool, error) {
	data, err := as.client.AddMetaInfo(ctx, &pb.AddMetaInfoRequest{
		ItemId:   id,
		Metadata: meta,
	})
	if err != nil {
		as.log.Error().Err(err).Msg("error setting metainfo")

		return false, errors.Wrap(err, "error setting metainfo")
	}

	return data.GetSuccess(), nil
}

// DeleteMetainfo deletes meta information for the record.
func (as *Client) DeleteMetainfo(ctx context.Context, id, key string) (bool, error) {
	data, err := as.client.RemoveMetaInfo(ctx, &pb.RemoveMetaInfoRequest{
		ItemId: id,
		Key:    key,
	})
	if err != nil {
		as.log.Error().Err(err).Msg("error deleting metainfo")

		return false, errors.Wrap(err, "error deleting metainfo")
	}

	return data.GetSuccess(), nil
}
