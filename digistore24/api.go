package digistore24

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	baseUrl = "https://www.digistore24.com/api/call/%s/"
)

type D24Result struct {
	ApiVersion  string `json:"api_version"`
	CurrentTime string `json:"current_time"`
	Result      string `json:"result"`
	Message     string `json:"message"`
	Code        int    `json:"code"`

	Data struct {
		Code     string      `json:"code"`
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

func buildUrl(function string, values url.Values) (*url.URL, error) {
	u, err := url.Parse(fmt.Sprintf(baseUrl, function))
	if err != nil {
		return nil, err
	}
	u.RawQuery = values.Encode()

	return u, nil
}

func queryValues(v url.Values, id int, values map[string]string) url.Values {
	for key, val := range values {
		v.Set(fmt.Sprintf("arg%d[%s]", id, key), val)
	}
	return v
}

func queryValue(v url.Values, id int, val string) url.Values {
	v.Set(fmt.Sprintf("arg%d", id), val)
	return v
}

func createGetRequest(u *url.URL, apiKey string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return
	}

	for k, v := range header(apiKey) {
		req.Header.Add(k, v)
	}

	return
}

func callRequest(req *http.Request) (status int, result D24Result, err error) {
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
