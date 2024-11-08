package alipanSdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const defaultHttpTimeout = 60 * time.Second

type alipanClient struct {
	appID     string
	appSecret string
	hc        *http.Client
}

type AuthCodeCallback func(code string)

// AlipanOptions 创建AlipanOptions客户端的选项
type AlipanOptions struct {
	AppID           string
	AppSecret       string
	HttpTimeout     time.Duration
	TransportConfig http.RoundTripper
}

// NewClient 新建AlipanSDK客户端
func NewClient(options ...*AlipanOptions) *alipanClient {
	appID := os.Getenv("ALIPAN_APP_ID")
	appSecret := os.Getenv("ALIPAN_APP_SECRET")
	randTripper := http.DefaultTransport
	httpTimeout := defaultHttpTimeout

	if len(options) == 1 {
		appID = options[0].AppID
		appSecret = options[0].AppSecret
		if options[0].TransportConfig == nil {
			randTripper = options[0].TransportConfig
		}
		if options[0].HttpTimeout != 0 {
			httpTimeout = options[0].HttpTimeout
		}
	}

	if len(appID) == 0 || len(appSecret) == 0 {
		return nil
	}

	client := &alipanClient{
		appID:     appID,
		appSecret: appSecret,
		hc: &http.Client{
			Transport: randTripper,
			Timeout:   httpTimeout,
		},
	}
	return client
}

func (c *alipanClient) sendReq(ctx context.Context, endpoint string, method string, token string, data interface{}) ([]byte, error) {
	var jsonData []byte
	callUrl := fmt.Sprintf("%s%s", alipanBaseUrl, endpoint)
	u, err := url.Parse(callUrl)
	if err != nil {
		return nil, fmt.Errorf("error Parse url: %v", err)
	}

	if method == http.MethodGet {
		if data != nil {
			queryParams := data.(map[string]string)
			q := u.Query()
			for key, value := range queryParams {
				q.Add(key, value)
			}
			u.RawQuery = q.Encode()
		}
	} else {
		jsonData, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("error Marshal data: %v", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	if method != http.MethodGet {
		req.Header.Set("Content-Type", contentTypeJson)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error call api: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error read resp: %v", err)
	}

	if !isAlipanSuccessCode(resp.StatusCode) {
		return nil, fmt.Errorf("error req code: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil

}

func isAlipanSuccessCode(code int) bool {
	return code >= 200 && code <= 300
}

func (c *alipanClient) ListenAuthCode(cb AuthCodeCallback) {

}
