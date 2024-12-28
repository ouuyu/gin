package common

import (
	"net/url"
	"strconv"

	"github.com/Calcium-Ion/go-epay/epay"
)

var EPayClient *epay.Client

func SetEPayClient(client *epay.Client) {
	EPayClient = client
}

func GeneratePayURL(amount float64, payType string, param string, tradeNo string) (string, error) {
	notify, _ := url.Parse(HomePage + "/api/v1/payment/notify")
	returnUrl, _ := url.Parse(HomePage)

	url, params, err := EPayClient.Purchase(&epay.PurchaseArgs{
		Type:           payType,
		ServiceTradeNo: tradeNo,
		Name:           "pay",
		Money:          strconv.FormatFloat(amount, 'f', 2, 64),
		Device:         epay.PC,
		NotifyUrl:      notify,
		ReturnUrl:      returnUrl,
	})
	if err != nil {
		return "", err
	}
	html := "<form id='alipaysubmit' name='alipaysubmit' action='" + url + "' method='POST'>"
	for key, value := range params {
		html += "<input type='hidden' name='" + key + "' value='" + value + "'/>"
	}
	html += "<input type='submit'>POST</form>"
	return html, nil
}
