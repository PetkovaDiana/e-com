package rec_service_client

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
)

type RecommendationServiceClient interface {
	SendOrderData(ctx context.Context, uuids []string) error
}

func NewRecommendationServiceClient(client *http.Client, log *logrus.Logger) RecommendationServiceClient {
	return &recommendationServiceClient{
		client: client,
		log:    log,
	}
}
