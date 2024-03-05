package pay_keeper

import (
	"clean_arch/internal/dto"
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
)

type PayKeeperClient interface {
	GerSecurityToken(ctx context.Context) (string, error)
	SecretKeyMD5Hashing(orderInfo *dto.OnlineOrderChecker) (string, error)
	PaymentReverse(ctx context.Context, orderData *dto.PayKeeperOrderCancel) error
	CheckOrderStatus(orderId int) error
	SendReceiptOnEmail(ctx context.Context, data *dto.PayKeeperReceiptSender) error
	CreatReceipt(ctx context.Context, data *dto.OrderData) (string, error)
}

func NewPayKeeperClient(authData BasicAuth, client *http.Client, log *logrus.Logger, secretKey string) PayKeeperClient {
	return &payKeeperClient{
		authData:  authData,
		client:    client,
		log:       log,
		secretKey: secretKey,
	}
}
