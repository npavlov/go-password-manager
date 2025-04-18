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
	Client       pb.MetadataServiceClient
	TokenManager auth.ITokenManager
	Log          *zerolog.Logger
}

// NewMetainfoClient  creates a new gRPC connection.
func NewMetainfoClient(conn *grpc.ClientConn, tokenManager auth.ITokenManager, log *zerolog.Logger) *Client {
	return &Client{
		conn:         conn,
		Client:       pb.NewMetadataServiceClient(conn),
		TokenManager: tokenManager,
		Log:          log,
	}
}

// GetMetainfo sends a register request to the server.
func (as *Client) GetMetainfo(ctx context.Context, id string) (map[string]string, error) {
	resp, err := as.Client.GetMetaInfoV1(ctx, &pb.GetMetaInfoV1Request{
		ItemId: id,
	})
	if err != nil {
		as.Log.Error().Err(err).Msg("error getting metainfo")

		return nil, errors.Wrap(err, "error getting metainfo")
	}

	return resp.GetMetadata(), nil
}

// SetMetainfo sets meta information for the record.
func (as *Client) SetMetainfo(ctx context.Context, id string, meta map[string]string) (bool, error) {
	data, err := as.Client.AddMetaInfoV1(ctx, &pb.AddMetaInfoV1Request{
		ItemId:   id,
		Metadata: meta,
	})
	if err != nil {
		as.Log.Error().Err(err).Msg("error setting metainfo")

		return false, errors.Wrap(err, "error setting metainfo")
	}

	return data.GetSuccess(), nil
}

// DeleteMetainfo deletes meta information for the record.
func (as *Client) DeleteMetainfo(ctx context.Context, id, key string) (bool, error) {
	data, err := as.Client.RemoveMetaInfoV1(ctx, &pb.RemoveMetaInfoV1Request{
		ItemId: id,
		Key:    key,
	})
	if err != nil {
		as.Log.Error().Err(err).Msg("error deleting metainfo")

		return false, errors.Wrap(err, "error deleting metainfo")
	}

	return data.GetSuccess(), nil
}
