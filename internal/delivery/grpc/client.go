package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/pkg/proto/gengrpc"
	"github.com/1Storm3/flibox-api/pkg/sys"
)

type ClientConnInterface interface {
	GetRecommendations(ctx context.Context, films []*gengrpc.Film) ([]string, error)
}

type ClientConn struct {
	conn   *grpc.ClientConn
	client gengrpc.RecommendationUseCaseClient
}

func NewClient(config *config.Config) (*ClientConn, error) {
	conn, err := grpc.NewClient(config.App.GrpcServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, sys.NewError(sys.ErrGRPCConnection, fmt.Sprintf("grpc connection error: %s", err.Error()))
	}

	client := gengrpc.NewRecommendationUseCaseClient(conn)
	return &ClientConn{conn: conn, client: client}, nil
}

func (c *ClientConn) Close() error {
	return c.conn.Close()
}

func (c *ClientConn) GetRecommendations(ctx context.Context, films []*gengrpc.Film) ([]string, error) {
	request := &gengrpc.RecommendationRequest{
		Films: films,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	response, err := c.client.GetRecommendations(ctx, request)
	if err != nil {
		return nil, sys.NewError(sys.ErrGRPCConnection, fmt.Sprintf("grpc connection error: %s", err.Error()))
	}

	return response.Films, nil
}
