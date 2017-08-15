package wework

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	GetAccessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
)

type GetAccessTokenResponse struct {
	Code        int    `json:"errcode,omitempty"`
	Message     string `json:"errmsg,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
}

func (c *AgentClient) getAccessToken() (string, error) {
	c.token.mu.Lock()
	defer c.token.mu.Unlock()

	now := time.Now()
	if now.Before(c.token.expiresAt) {
		return c.token.token, nil
	}

	resp, err := c.requestAccessToken()
	if err != nil {
		return "", err
	}

	c.token.token = resp.AccessToken
	c.token.expiresAt = now.Add(time.Duration(resp.ExpiresIn) * time.Second)

	return c.token.token, err
}

func (c *AgentClient) requestAccessToken() (*GetAccessTokenResponse, error) {
	v := url.Values{}
	v.Set("corpid", c.corpID)
	v.Set("corpsecret", c.secret)

	url := GetAccessTokenURL + "?" + v.Encode()

	httpResp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get error: %v", err)
	}

	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("http response read error: %v", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("illegal http status code: %v", httpResp.StatusCode)
	}

	var resp GetAccessTokenResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("json.Unmarshal error: %v", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("token.get api error: %+v", resp)
	}

	return &resp, nil
}
