package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type NutanixClient struct {
	hostname string
	username string
	password string
	apiPath  string
	client   *http.Client
}

func NewNutanixClient(hostname, username, password, apiPath string) (*NutanixClient, error) {
	return &NutanixClient{
		hostname,
		username,
		password,
		apiPath,
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}, nil
}

func (c *NutanixClient) DoRequest(method, path string, query map[string][]string, data interface{}) ([]byte, error) {
	reqStream := bytes.NewBuffer(nil)
	if data != nil {
		reqBodyBytes, err := json.Marshal(data)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		reqStream.Write(reqBodyBytes)
	}

	u := url.URL{
		Scheme: "https",
		Host:   c.hostname,
		Path:   fmt.Sprintf("%s/%s", c.apiPath, path),
	}

	if query != nil {
		values := u.Query()
		for key, val := range query {
			for _, v := range val {
				values.Add(key, v)
			}
		}
	}

	req, err := http.NewRequest(method, u.String(), reqStream)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Close = true

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Error with request, got status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	respStream := bytes.NewBuffer(nil)

	_, err = respStream.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return respStream.Bytes(), nil
}
