package controller

import (
	"context"
	"e-montir/api/handler"
	"e-montir/model"
	"e-montir/pkg/validator"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

type paymentCtx struct {
	orderModel model.Order
	cartModel  model.Cart
}

type Payment interface {
	Pay(ctx context.Context, orderID string) (*TransactionResponse, error)
}

func NewPayment(orderModel model.Order, cartModel model.Cart) Payment {
	return &paymentCtx{
		orderModel: orderModel,
		cartModel:  cartModel,
	}
}

type (
	PaymentDetail struct {
		TransactionTime        string          `json:"transaction_time"`
		TransactionStatus      string          `json:"transaction_status"`
		TransactionID          string          `json:"transaction_id"`
		StatusMessage          string          `json:"status_message"`
		StatusCode             string          `json:"status_code"`
		SignatureKey           string          `json:"signature_key"`
		PaymentType            string          `json:"payment_type"`
		OrderID                string          `json:"order_id"`
		MerchantID             string          `json:"merchant_id"`
		MaskedCard             string          `json:"masked_card"`
		GrossAmount            string          `json:"gross_amount"`
		FraudStatus            string          `json:"fraud_status"`
		Eci                    string          `json:"eci"`
		Currency               string          `json:"currency"`
		ChannelResponseMessage string          `json:"channel_response_message"`
		ChannelResponseCode    string          `json:"channel_response_code"`
		CardType               string          `json:"card_type"`
		Bank                   string          `json:"bank"`
		ApprovalCode           string          `json:"approval_code"`
		Issuer                 string          `json:"issuer"`
		Acquirer               string          `json:"acquirer"`
		PermataVaNumber        string          `json:"permata_va_number"`
		PaidAmount             int             `json:"paid_amount"`
		BillerCode             string          `json:"biller_code"`
		BillKey                string          `json:"bill_key"`
		Store                  string          `json:"store"`
		SettlementTime         string          `json:"settlement_time"`
		FinishURL              string          `json:"finish_url"`
		PaymentNotifURL        string          `json:"payment_notification_url"`
		VANumbers              []VANumber      `json:"va_numbers"`
		PaymentAmounts         []PaymentAmount `json:"payment_amounts"`
	}

	TransactionResponse struct {
		Code    string       `json:"code"`
		Message string       `json:"message"`
		Data    RedirectData `json:"data"`
	}

	PaymentNotif struct {
		Code    string        `json:"code"`
		Message string        `json:"message"`
		Data    PaymentDetail `json:"data"`
	}

	PaymentNotifResponse struct {
		TransactionStatus string `json:"transaction_status"`
	}

	RedirectData struct {
		Token       string `json:"token"`
		RedirectURL string `json:"redirect_url"`
		FinishURL   string `json:"finish_url"`
	}

	VANumber struct {
		VANumber string `json:"va_number"`
		Bank     string `json:"bank"`
	}
	PaymentAmount struct {
		PaidAt string `json:"paid_at"`
		Amount string `json:"amount"`
	}

	PaymentRequest struct {
		OrderID          string `json:"order_id"`
		TotalPriceString string `json:"total_price"`
		TotalPrice       int
	}

	HttpRequest struct {
		Client *http.Client
	}
)

func (req *PaymentRequest) ValidatePaymentRequest() ([]handler.Fields, error) {
	var count int
	var fields []handler.Fields
	err := validator.ValidateOrderID(req.OrderID)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "order_id",
			Message: err.Error(),
		})
	}

	if count != 0 {
		return fields, errors.New("validation-failed")
	}
	return nil, nil
}

func (c *paymentCtx) Pay(ctx context.Context, orderID string) (*TransactionResponse, error) {
	order, err := c.orderModel.CheckOrder(ctx, orderID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when checkingOrder: %w", err)).Send()
		return nil, err
	}

	totalPrice := int(order.TotalPrice)
	transactionRes := TransactionResponse{}
	payload := strings.NewReader(fmt.Sprintf(`{
	    "transaction_detail":{
	        "order_id": "%s",
	        "gross_amount": %d
	    },
	    "redirect":{
	        "finish_url": "%s",
			"payment_notification_url": "%s"
	    }
	}`, orderID, totalPrice, os.Getenv("PAYMENT_REDIRECT_BASE_URL"), os.Getenv("PAYMENT_NOTIFICATION_URL")))

	client := &http.Client{}
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/transactions", os.Getenv("PAYMENT_BASE_URL")),
		payload,
	)

	if err != nil {
		return nil, &handler.InternalServerError
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("token", os.Getenv("PAYMENT_TOKEN"))

	res, err := client.Do(req)
	if err != nil {
		return nil, &handler.InternalServerError
	}

	defer res.Body.Close()

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, &handler.InternalServerError
	}

	err = json.Unmarshal([]byte(string(resData)), &transactionRes)
	if err != nil {
		return nil, err
	}

	transactionRes.Data.FinishURL = os.Getenv("PAYMENT_REDIRECT_BASE_URL")
	return &transactionRes, nil
}
