package order_service_client

import (
	"clean_arch/internal/dto"
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
)

type OrderServiceClient interface {
	SendOrderData(ctx context.Context, orderData *dto.InvoiceData) (string, error)
}

func NewOrderServiceClient(client *http.Client, log *logrus.Logger) OrderServiceClient {
	return &orderServiceClient{
		client: client,
		log:    log,
	}
}
