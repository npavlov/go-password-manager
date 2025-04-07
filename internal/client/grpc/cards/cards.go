package cards

import (
	"context"
	"time"

	pb "github.com/npavlov/go-password-manager/gen/proto/card"
	"github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// Client GRPCClient handles communication with the gRPC server
type Client struct {
	conn         *grpc.ClientConn
	client       pb.CardServiceClient
	tokenManager *auth.TokenManager
	log          *zerolog.Logger
}

// NewCardClient  creates a new gRPC connection
func NewCardClient(conn *grpc.ClientConn, tokenManager *auth.TokenManager, log *zerolog.Logger) *Client {
	return &Client{
		conn:         conn,
		client:       pb.NewCardServiceClient(conn),
		tokenManager: tokenManager,
		log:          log,
	}
}

// GetCard sends a register request to the server
func (as *Client) GetCard(ctx context.Context, id string) (*pb.CardData, time.Time, error) {
	resp, err := as.client.GetCard(ctx, &pb.GetCardRequest{
		CardId: id,
	})

	if err != nil {
		as.log.Error().Err(err).Msg("error getting password")

		return nil, time.Time{}, errors.Wrap(err, "error getting password")
	}

	return resp.Card, resp.LastUpdate.AsTime(), nil
}

func (as *Client) UpdateCard(ctx context.Context, id, cardNum, expDate, Cvv, cardHolder string) error {
	_, err := as.client.UpdateCard(ctx, &pb.UpdateCardRequest{
		CardId: id,
		Data: &pb.CardData{
			CardNumber:     cardNum,
			ExpiryDate:     expDate,
			Cvv:            Cvv,
			CardholderName: cardHolder,
		},
	})

	if err != nil {
		as.log.Error().Err(err).Msg("error updating card")

		return errors.Wrap(err, "error updating card")
	}

	return nil
}

func (as *Client) StoreCard(ctx context.Context, cardNum, expDate, Cvv, cardHolder string) (string, error) {
	resp, err := as.client.StoreCard(ctx, &pb.StoreCardRequest{
		Card: &pb.CardData{
			CardNumber:     cardNum,
			ExpiryDate:     expDate,
			Cvv:            Cvv,
			CardholderName: cardHolder,
		},
	})

	if err != nil {
		as.log.Error().Err(err).Msg("error storing card")

		return "", errors.Wrap(err, "error storing card")
	}

	return resp.CardId, nil
}

// DeleteCard delete card
func (as *Client) DeleteCard(ctx context.Context, id string) (bool, error) {
	resp, err := as.client.DeleteCard(ctx, &pb.DeleteCardRequest{
		CardId: id,
	})

	if err != nil {
		as.log.Error().Err(err).Msg("error deleting card")

		return false, errors.Wrap(err, "error deleting card")
	}

	return resp.Ok, nil
}
