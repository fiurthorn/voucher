package digistore24

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/fiurthorn/voucher/config"
)

const (
	baseUrl = "https://www.digistore24.com/api/call/createVoucher/"
)

type D24Result struct {
	ApiVersion  string `json:"api_version"`
	CurrentTime string `json:"current_time"`
	Result      string `json:"result"`
	Message     string `json:"message"`
	Code        int    `json:"code"`

	Data struct {
		Code     json.Number `json:"code"`
		CouponId json.Number `json:"coupon_id"`
		Note     string      `json:"note"`
	} `json:"data"`
}

func header(apikey string) map[string]string {
	return map[string]string{
		"X-DS-API-KEY": apikey,
		"Accept":       "application/json",
		"User-Agent":   "node-XMLHttpRequest",
	}
}

func query(voucher string, products []string) url.Values {
	return url.Values{
		"arg1[code]":             []string{voucher},
		"arg1[product_ids]":      []string{strings.Join(products, ",")},
		"arg1[first_rate]":       []string{"100"},
		"arg1[is_count_limited]": []string{"true"},
		"arg1[count_left]":       []string{"1"},
		"arg1[upgrade_policy]":   []string{"not_valid"},
	}
}

func buildUrl(voucher string, products []string) (*url.URL, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = query(voucher, products).Encode()

	return u, nil
}

func buildRequest(apiKey, voucher string, products []string) (req *http.Request, err error) {
	u, err := buildUrl(voucher, products)
	if err != nil {
		return
	}

	req, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return
	}

	for k, v := range header(apiKey) {
		req.Header.Add(k, v)
	}

	return
}

func call(apiKey, voucher string, products []string) (status int, result D24Result, err error) {
	result = D24Result{}

	req, err := buildRequest(apiKey, voucher, products)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	status = resp.StatusCode

	if status == 200 {
		var body []byte

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			return
		}
	}

	return
}

type CallResult struct {
	Voucher string
	Status  int
	Result  D24Result
	Err     error
}

func Call(vouchers string, drain chan<- CallResult) {
	voucherIds := strings.Split(strings.ReplaceAll(vouchers, "\r\n", "\n"), "\n")

	for _, v := range voucherIds {
		s, r, e := call(config.Config.ApiKey, v, config.Config.Products)
		drain <- CallResult{v, s, r, e}
	}

	close(drain)
}
