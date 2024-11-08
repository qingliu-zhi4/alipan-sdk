package alipanSdk

import (
	"context"
	"testing"
)

func TestAlipanClient_RequestAuthQrCode(t *testing.T) {
	options := &AlipanOptions{
		AppID:     "c7dda905da004f76a4acd2deaa61fd46",
		AppSecret: "4fa7695c942e459781d6b043715367e2",
	}
	c := NewClient(options)
	scope := AllScope
	data, err := c.RequestAuthQrCode(context.TODO(), scope)
	if err != nil {
		t.Errorf("cont req qr, error: %v", err)
	}
	t.Log(data.QrCodeUrl)
}

func TestAlipanClient_SimpleAuthorizeUrl(t *testing.T) {
	options := &AlipanOptions{
		AppID:     "c7dda905da004f76a4acd2deaa61fd46",
		AppSecret: "4fa7695c942e459781d6b043715367e2",
	}
	c := NewClient(options)
	authUrl, err := c.SimpleAuthorizeUrl(context.TODO(), "https://zhi4ai.com")
	if err != nil {
		t.Errorf("cont req auth, error: %v", err)
	}
	t.Log(authUrl)
}
