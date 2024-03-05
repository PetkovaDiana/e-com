package pay_keeper

import (
	"bytes"
	"clean_arch/internal/dto"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	reverseUrl       = "https://ufaelectro.server.paykeeper.ru/change/payment/reverse/"
	securityTokenUrl = "https://ufaelectro.server.paykeeper.ru/info/settings/token/"
	checkStatus      = "https://ufaelectro.server.paykeeper.ru/info/payments/byid/"
	receiptSender    = "https://ufaelectro.server.paykeeper.ru/change/invoice/send/"
	receiptCreator   = "https://ufaelectro.server.paykeeper.ru/change/invoice/preview/"
)

type payKeeperClient struct {
	authData  BasicAuth
	client    *http.Client
	log       *logrus.Logger
	secretKey string
}

type BasicAuth struct {
	Username string
	Password string
}

func (p *payKeeperClient) GerSecurityToken(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", securityTokenUrl, nil)
	if err != nil {
		return "", err
	}

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(p.authData.Username+":"+p.authData.Password))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var response struct {
		Token string `json:"token"`
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return response.Token, nil
}

func (p *payKeeperClient) SecretKeyMD5Hashing(orderInfo *dto.OnlineOrderChecker) (string, error) {
	hash := md5.Sum([]byte(fmt.Sprintf("%s%.2f%s%s%s", orderInfo.ID, orderInfo.Sum, orderInfo.ClientID, orderInfo.OrderID, p.secretKey)))
	if fmt.Sprintf("%x", hash) != orderInfo.Key {
		return "", fmt.Errorf("Hash mismatch")
	}

	orderInfo.ClientID = strings.TrimSuffix(strings.SplitAfter(orderInfo.ClientID, "(")[1], ")")

	responseStr := fmt.Sprintf("OK %x", md5.Sum([]byte(orderInfo.ID+p.secretKey)))
	return responseStr, nil
}

func (p *payKeeperClient) PaymentReverse(ctx context.Context, orderData *dto.PayKeeperOrderCancel) error {
	// Создание URL объекта
	u, err := url.Parse(reverseUrl)
	if err != nil {
		panic(err)
	}

	// Задание параметров в Query параметрах
	q := u.Query()
	q.Set("id", strconv.Itoa(orderData.PaymentID))
	q.Set("partial", "false")
	q.Set("token", orderData.SecurityToken)

	// Создание запроса
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBufferString(q.Encode()))

	if err != nil {
		return err
	}

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(p.authData.Username+":"+p.authData.Password))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", authHeader)

	// Отправка запроса и получение ответа
	resp, err := p.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Обработка ответа
	var response struct {
		Success string `json:"result"`
		Msg     string `json:"msg"`
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return err
	}

	if response.Success != "success" {
		return fmt.Errorf("%s", response.Msg)
	}
	return nil
}

func (p *payKeeperClient) CheckOrderStatus(orderId int) error {
	u, err := url.Parse(checkStatus)

	if err != nil {
		return err
	}

	q := u.Query()

	q.Set("id", strconv.Itoa(orderId))

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(p.authData.Username+":"+p.authData.Password))
	req.Header.Set("Authorization", authHeader)

	if err != nil {
		return err
	}

	resp, err := p.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var response []struct {
		Status string `json:"status"`
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return err
	}

	if response[0].Status != "refunded" {
		return fmt.Errorf("not refunded")
	}
	return nil
}

func (p *payKeeperClient) SendReceiptOnEmail(ctx context.Context, data *dto.PayKeeperReceiptSender) error {
	u, err := url.Parse(receiptSender)

	if err != nil {
		return err
	}

	q := u.Query()

	q.Set("id", data.InvoiceID)
	q.Set("token", data.SecurityToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBufferString(q.Encode()))

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(p.authData.Username+":"+p.authData.Password))
	req.Header.Set("Authorization", authHeader)

	resp, err := p.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var response struct {
		Success string `json:"result"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	if response.Success != "success" {
		return fmt.Errorf("not success")
	}
	return nil
}

func (p *payKeeperClient) CreatReceipt(ctx context.Context, data *dto.OrderData) (string, error) {
	u, err := url.Parse(receiptCreator)

	if err != nil {
		return "", err
	}

	q := u.Query()

	q.Set("pay_amount", fmt.Sprint(data.ReceiptData.PayAmount))
	q.Set("clientid", strconv.Itoa(data.ReceiptData.ClientID))
	q.Set("orderid", strconv.Itoa(data.NewOrderID))
	q.Set("service_name", "Оплата БЭСМ")
	q.Set("client_email", data.ReceiptData.ClientEmail)
	q.Set("client_phone", data.ReceiptData.ClientPhone)
	q.Set("expiry", "2023-03-16")
	q.Set("token", data.ReceiptData.Token)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBufferString(q.Encode()))

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(p.authData.Username+":"+p.authData.Password))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", authHeader)

	resp, err := p.client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var response struct {
		InvoiceId string `json:"invoice_id"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if response.InvoiceId == "" {
		return "", fmt.Errorf("not success")
	}

	return response.InvoiceId, nil
}
