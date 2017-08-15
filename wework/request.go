package wework

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	ContentTypeJson string = "application/json"
)

func (c *AgentClient) getJSON(url string, resp interface{}) error {
	log.Printf("Get %s", url)

	reqURL, err := c.urlWithToken(url)
	if err != nil {
		return err
	}

	httpResp, err := http.Get(reqURL)
	if err != nil {
		return fmt.Errorf("http.Get error: %v", err)
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("http response read error: %v", err)
	}

	log.Printf("Response %d %s", httpResp.StatusCode, body)

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("illegal http status code: %v", httpResp.StatusCode)
	}

	if err := json.Unmarshal(body, resp); err != nil {
		return fmt.Errorf("json.Unmarshal error: %v", err)
	}

	return nil
}

func (c *AgentClient) postJSON(url string, req interface{}, resp interface{}) error {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(req); err != nil {
		return fmt.Errorf("json.Encode error: %v", err)
	}

	reqURL, err := c.urlWithToken(url)
	if err != nil {
		return err
	}

	log.Printf("Post %s Body: %s", url, buf.String())

	httpResp, err := http.Post(reqURL, ContentTypeJson, buf)
	if err != nil {
		return fmt.Errorf("http.Post error: %v", err)
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("http response read error: %v", err)
	}

	log.Printf("Response %d %s", httpResp.StatusCode, body)

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("illegal http status code: %v", httpResp.StatusCode)
	}

	if err := json.Unmarshal(body, resp); err != nil {
		return fmt.Errorf("json.Unmarshal error: %v", err)
	}

	return nil
}

func (c *AgentClient) urlWithToken(url string) (string, error) {
	token, err := c.getAccessToken()
	if err != nil {
		return "", fmt.Errorf("get access token error: %v", err)
	}

	reqURL := url + "?access_token=" + token
	return reqURL, nil
}
