package agent

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/mohaijiang/computeshare-client/internal/conf"
)

func NewHttpConnection(confData *conf.Data) (*transhttp.Client, func(), error) {
	client, err := transhttp.NewClient(
		context.Background(),
		transhttp.WithMiddleware(
			recovery.Recovery(),
		),
		transhttp.WithEndpoint(confData.ComputerPowerApi),
	)

	cleanup := func() {
		_ = client.Close()
	}

	return client, cleanup, err
}
