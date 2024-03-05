package client_1c

import (
	"bytes"
	"clean_arch/internal/dto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const urlOrder = "http://81.30.179.149/ut_besm/hs/WebsiteAPI/Order"
const urlCancelOrder = "http://81.30.179.149/ut_besm/hs/WebsiteAPI/Order/OrderCancel"

type client1C struct {
	client *http.Client
	log    *logrus.Logger
	Auth   *BasicAuth
}

type BasicAuth struct {
	Username string
	Password string
}

func (c *client1C) SendOrderData(ctx context.Context, order1CInfo *dto.Order1C) error {

	bodyData, err := json.Marshal(&order1CInfo)

	if err != nil {
		c.log.Errorf("error occured client 1c initializng request: %v", err)
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlOrder, bytes.NewReader(bodyData))

	if err != nil {
		c.log.Errorf("error occured client 1c initializng request: %v", err)
		return err
	}

	req.SetBasicAuth("Website", "152684")

	res, err := c.client.Do(req)

	defer res.Body.Close()

	if err != nil {
		c.log.Errorf("error occured client 1c client do request: %v", err)
		return err
	}

	c.log.Printf("client 1c request status code: %s", res.Status)

	return nil
}

func (c *client1C) SendCancelOrderData(ctx context.Context, orderId int) {

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlCancelOrder, nil)

	if err != nil {
		c.log.Errorf("error occured client 1c initializng request: %v", err)
		return
	}

	req.SetBasicAuth("Website", "152684")

	query := req.URL.Query()
	query.Add("id", strconv.Itoa(orderId))
	req.URL.RawQuery = query.Encode()

	fmt.Println(req.URL)

	res, err := c.client.Do(req)

	defer res.Body.Close()

	if err != nil {
		c.log.Errorf("error occured client 1c client do request: %v", err)
		return
	}

	c.log.Printf("client 1c request status code: %s", res.Status)
}
