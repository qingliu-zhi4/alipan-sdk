package alipanSdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// 权限列表
// ref: https://www.yuque.com/aliyundrive/zpfszx/dspik0
// UserBaseScope 获取用户ID、头像、昵称
// FileReadScope 读取云盘所有文件
// FileWriteScope 写入云盘所有文件
// AlbumSharedReadScope 读取共享相薄文件
// FileShareWriteScope 文件分享
const (
	UserBaseScope = 1 << iota
	FileReadScope
	FileWriteScope
	AlbumSharedReadScope
	FileShareWriteScope
)

// AllScope 全Scope
const AllScope = UserBaseScope | FileReadScope | FileWriteScope | AlbumSharedReadScope | FileShareWriteScope

var scoreMap = map[int]string{
	UserBaseScope:        "user:base",
	FileReadScope:        "file:all:read",
	FileWriteScope:       "file:all:read",
	AlbumSharedReadScope: "album:shared:read",
	FileShareWriteScope:  "file:share:write",
}

type authQRReq struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
	Width        int64    `json:"width"`
	Height       int64    `json:"height"`
}

// AuthQRRes 获取Auth二维码的返回值
type AuthQRRes struct {
	QrCodeUrl string `json:"qrCodeUrl"`
	Sid       string `json:"sid"`
}

// RequestAuthQrCode 获取扫码登录用的二维码，默认大小430*430
func (c *alipanClient) RequestAuthQrCode(ctx context.Context, scope int) (*AuthQRRes, error) {
	return c.RequestAuthQrCodeWithSize(ctx, scope, 430, 430)
}

// RequestAuthQrCodeWithSize 获取扫码登录用的二维码，自定义图片大小
func (c *alipanClient) RequestAuthQrCodeWithSize(ctx context.Context, scope int, width int, height int) (*AuthQRRes, error) {
	scopes := formScopesStr(scope)
	req := &authQRReq{
		ClientID:     c.appID,
		ClientSecret: c.appSecret,
		Scopes:       scopes,
		Width:        int64(width),
		Height:       int64(height),
	}

	respData, err := c.sendReq(ctx, reqQRUrl, http.MethodPost, "", req)
	if err != nil {
		return nil, fmt.Errorf("req qr fail %v", err)
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(respData, &responseData)
	if err != nil {
		return nil, fmt.Errorf("error Unmarshal resp: %v", err)
	}

	qrData := &AuthQRRes{
		QrCodeUrl: responseData["qrCodeUrl"].(string),
		Sid:       responseData["sid"].(string),
	}

	return qrData, nil
}

// GetAuthorizeUrl 获取授权URL, 需要全参数，如果觉得麻烦可以使用SimpleAuthorizeUrl
func (c *alipanClient) GetAuthorizeUrl(ctx context.Context, redirectUri string, scope int, reLogin bool, autoLogin bool, drive int, state string) (string, error) {
	scopes := formScopesStr(scope)
	drivers := formDriversStr(drive)
	encodedURL := url.QueryEscape(redirectUri)

	authReqMap := map[string]string{
		"client_id":     c.appID,
		"redirect_uri":  encodedURL,
		"scope":         strArray2StrComma(scopes),
		"response_type": "code",
		"relogin":       bool2Str(reLogin),
		"auto_login":    bool2Str(autoLogin),
	}

	if len(state) > 0 {
		authReqMap["state"] = state
	}

	if len(drivers) > 0 {
		authReqMap["drive"] = strArray2StrComma(drivers)
	}

	authRes, err := c.sendReq(ctx, reqAuthUrl, http.MethodGet, "", authReqMap)
	if err != nil {
		return "", fmt.Errorf("req auth error: %v", err)
	}
	return string(authRes), nil
}

// SimpleAuthorizeUrl 获取授权URL的建议方法，同义使用默认参数，只需要传入授权回调地址
func (c *alipanClient) SimpleAuthorizeUrl(ctx context.Context, redirectUri string) (string, error) {
	return c.GetAuthorizeUrl(ctx, redirectUri, AllScope, false, false, AllDrive, "")
}
