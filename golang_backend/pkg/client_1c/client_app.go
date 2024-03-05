package client_1c

import (
	"clean_arch/internal/dto"
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Client1C interface {
	SendOrderData(ctx context.Context, order1CInfo *dto.Order1C) error
	SendCancelOrderData(ctx context.Context, orderId int)
}

func NewClient1C(client *http.Client, log *logrus.Logger, username, password string) Client1C {
	return &client1C{
		client: client,
		log:    log,
		Auth: &BasicAuth{
			Username: username,
			Password: password,
		},
	}
}
