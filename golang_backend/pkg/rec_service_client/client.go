package rec_service_client

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
)

const orderDataUrl = "http://recommendation:81/recommendation"

type recommendationServiceClient struct {
	client *http.Client
	log    *logrus.Logger
}

func (s *recommendationServiceClient) SendOrderData(ctx context.Context, uuids []string) error {

	u, err := url.Parse(orderDataUrl)

	if err != nil {
		return err
	}

	q := u.Query()

	q.Set("products_uuids", strings.Join(uuids, ","))

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), nil)

	if err != nil {
		s.log.Errorf("error occured recommendation initializng request: %v", err)
		return err
	}
	res, err := s.client.Do(req)

	if err != nil {
		s.log.Errorf("error occured recommendation client do request: %v", err)
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 422 {
		s.log.Errorf("error with status code: %d", res.StatusCode)
		return fmt.Errorf("status code is: %d", res.StatusCode)
	}
	return nil
}
