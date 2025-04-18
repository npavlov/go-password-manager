package cards

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	pb "github.com/npavlov/go-password-manager/gen/proto/card"
	"github.com/npavlov/go-password-manager/internal/client/auth"
)

// Client GRPCClient handles communication with the gRPC server.
type Client struct {
	conn         *grpc.ClientConn
	Client       pb.CardServiceClient
	TokenManager auth.ITokenManager
	Log          *zerolog.Logger
}

// NewCardClient  creates a new gRPC connection.
func NewCardClient(conn *grpc.ClientConn, tokenManager auth.ITokenManager, log *zerolog.Logger) *Client {
	return &Client{
		conn:         conn,
		Client:       pb.NewCardServiceClient(conn),
		TokenManager: tokenManager,
		Log:          log,
	}
}

// GetCard sends a register request to the server.
func (as *Client) GetCard(ctx context.Context, id string) (*pb.CardData, time.Time, error) {
	resp, err := as.Client.GetCardV1(ctx, &pb.GetCardV1Request{
		CardId: id,
	})
	if err != nil {
		as.Log.Error().Err(err).Msg("error getting password")

		return nil, time.Time{}, errors.Wrap(err, "error getting password")
	}

	return resp.GetCard(), resp.GetLastUpdate().AsTime(), nil
}

func (as *Client) UpdateCard(ctx context.Context, id, cardNum, expDate, cvv, cardHolder string) error {
	_, err := as.Client.UpdateCardV1(ctx, &pb.UpdateCardV1Request{
		CardId: id,
		Data: &pb.CardData{
			CardNumber:     cardNum,
			ExpiryDate:     expDate,
			Cvv:            cvv,
			CardholderName: cardHolder,
		},
	})
	if err != nil {
		as.Log.Error().Err(err).Msg("error updating card")

		return errors.Wrap(err, "error updating card")
	}

	return nil
}

func (as *Client) StoreCard(ctx context.Context, cardNum, expDate, cvv, cardHolder string) (string, error) {
	resp, err := as.Client.StoreCardV1(ctx, &pb.StoreCardV1Request{
		Card: &pb.CardData{
			CardNumber:     cardNum,
			ExpiryDate:     expDate,
			Cvv:            cvv,
			CardholderName: cardHolder,
		},
	})
	if err != nil {
		as.Log.Error().Err(err).Msg("error storing card")

		return "", errors.Wrap(err, "error storing card")
	}

	return resp.GetCardId(), nil
}

// DeleteCard delete card.
func (as *Client) DeleteCard(ctx context.Context, id string) (bool, error) {
	resp, err := as.Client.DeleteCardV1(ctx, &pb.DeleteCardV1Request{
		CardId: id,
	})
	if err != nil {
		as.Log.Error().Err(err).Msg("error deleting card")

		return false, errors.Wrap(err, "error deleting card")
	}

	return resp.GetOk(), nil
}
