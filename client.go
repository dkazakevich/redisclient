package main

import (
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"io"
	"strings"
	"strconv"
	"errors"
)

type Client struct {
	url string
	httpClient http.Client
}

func newClient(url string) *Client {

	if strings.HasSuffix(url, "/") == false {
		url += "/"
	}

	return &Client{url, http.Client{Timeout: time.Second * 2}}
}

func (c *Client) keys() ([]interface{}, error) {

	respValue, err := c.httpCall(http.MethodGet, c.url+"keys", nil)

	return respValue["keys"].([]interface{}), err
}

func (c *Client) get(key string) (interface{}, error) {

	respValue, err := c.httpCall(http.MethodGet, c.url + "values/" + key,  nil)

	return respValue["value"], err
}

func (c *Client) getListElement(key string, pos int) (interface{}, error) {

	respValue, err := c.httpCall(http.MethodGet, c.url + "values/" + key + "?pos=" + strconv.Itoa(pos),  nil)

	return respValue["value"], err
}

func (c *Client) getDictValue(key string, dictKey string) (interface{}, error) {

	respValue, err := c.httpCall(http.MethodGet, c.url + "values/" + key + "?key=" + dictKey,  nil)

	return respValue["value"], err
}

func (c *Client) put(key string, value interface{}) (interface{}, error) {

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(value)

	respValue, err := c.httpCall(http.MethodPut, c.url + "values/" + key,  b)

	return respValue, err
}

func (c *Client) httpCall(method, url string, body io.Reader) (map[string]interface{}, error) {

	var respValue map[string]interface{}
	req, err := http.NewRequest(method, url,  body)
	if err == nil {
		req.Header.Set("Content-Type", "application/json")

		resp, doErr := c.httpClient.Do(req)
		if doErr != nil {
			err = doErr
		} else {
			defer resp.Body.Close()

			err = json.NewDecoder(resp.Body).Decode(&respValue)
			if err == nil {
				if respValue["error"] != nil {
					err = errors.New(respValue["error"].(string))
				}
			}
		}
	}

	return respValue, err
}