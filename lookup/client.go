package lookup

import (
	"context"
	lookupv1 "github.com/halooid/backend/go-shared/lookup/gen/go/lookup/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client lookupv1.LookupServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: lookupv1.NewLookupServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetValues(ctx context.Context, lookupKey string) ([]*lookupv1.LookupValue, error) {
	resp, err := c.client.GetLookupValues(ctx, &lookupv1.GetLookupValuesRequest{
		LookupKey: lookupKey,
	})
	if err != nil {
		return nil, err
	}
	return resp.Values, nil
}
