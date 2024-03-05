package order_service_client

import (
	"bytes"
	"clean_arch/internal/dto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

const orderDataUrl = "http://order:79/create_payment_invoice"

type orderServiceClient struct {
	client *http.Client
	log    *logrus.Logger
}

func (s *orderServiceClient) SendOrderData(ctx context.Context, orderData *dto.InvoiceData) (string, error) {

	bodyData, err := json.Marshal(&orderData)

	if err != nil {
		return "", err
	}

	u, err := url.Parse(orderDataUrl)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(bodyData))

	if err != nil {
		s.log.Errorf("error occured order initializng request: %v", err)
		return "", err
	}
	//return "", nil
	resp, err := s.client.Do(req)

	if err != nil {
		s.log.Errorf("error occured order service client do request: %v", err)
		return "", err
	}

	defer resp.Body.Close()

	var response struct {
		Status int    `json:"status"`
		Url    string `json:"url"`
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		s.log.Errorf("error with status code: %d", resp.StatusCode)
		return "", fmt.Errorf("status code is: %d", resp.StatusCode)
	}
	return response.Url, nil
}
