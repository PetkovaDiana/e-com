package bitrix_client

import (
	"clean_arch/internal/dto"
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
)

type BitrixClient interface {
	Request(ctx context.Context, url string) error
	CallBackHeader(callBackInfo *dto.RequestCall) string
	CallBackFooter(callBackInfo *dto.RequestCall) string
	CallBackLK(callBackInfo *dto.RequestCall) string
	CallBackVacancy(requestVacancyInfo *dto.RequestVacancy) string
}

func NewClient(client *http.Client, log *logrus.Logger, key ...string) BitrixClient {
	return &bitrixClient{
		client: client,
		log:    log,
		methods: map[Method]string{
			CallBackHeader:  key[0],
			CallBackFooter:  key[1],
			CallBackLK:      key[2],
			CallBackVacancy: key[3],
		},
	}
}
