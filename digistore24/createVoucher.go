package digistore24

import (
	"net/http"
	"net/url"
)

func createVoucherValues(voucher string, products string) url.Values {
	return queryValues(url.Values{}, 1, map[string]string{
		"code":             voucher,
		"product_ids":      products,
		"first_rate":       "100",
		"is_count_limited": "true",
		"count_left":       "1",
		"upgrade_policy":   "not_valid",
	})
}

func createVoucherRequest(apiKey, voucher string, products string) (req *http.Request, err error) {
	u, err := buildUrl("createVoucher", createVoucherValues(voucher, products))
	if err != nil {
		return
	}

	return createGetRequest(u, apiKey)
}

type CreateVoucherResult struct {
	Voucher string
	Status  int
	Result  D24Result
	Err     error
}

func CreateVoucher(apiKey, voucher string, products string) CreateVoucherResult {
	req, err := createVoucherRequest(apiKey, voucher, products)
	if err != nil {
		return CreateVoucherResult{Err: err}
	}

	s, r, e := callRequest(req)
	return CreateVoucherResult{voucher, s, r, e}
}

func CreateVouchers(apiKey string, vouchers []string, products string, drain chan<- CreateVoucherResult) {
	defer close(drain)
	for _, v := range vouchers {
		drain <- CreateVoucher(apiKey, v, products)
	}
}
