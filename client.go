package redisclient

import (
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"strings"
	"strconv"
	"errors"
)

type Client struct {
	url string
	httpClient http.Client
}

func NewClient(url string) *Client {

	if strings.HasSuffix(url, "/") == false {
		url += "/"
	}

	return &Client{url, http.Client{Timeout: time.Second * 2}}
}

func (c *Client) Keys() ([]string, error) {

	respValue, err := c.httpCall(http.MethodGet, c.url+"keys", nil)
	keys := respValue["keys"].([]interface{})

	result := make([]string, len(keys))
	for i := range keys {
		result[i] = keys[i].(string)
	}

	return result, err
}

func (c *Client) Get(key string) (interface{}, error) {

	respValue, err := c.httpCall(http.MethodGet, c.url + "values/" + key,  nil)

	return respValue["value"], err
}

func (c *Client) GetTtl(key string) (interface{}, error) {

	respValue, err := c.httpCall(http.MethodGet, c.url + "ttl/" + key,  nil)

	return respValue["value"], err
}

func (c *Client) GetListElement(key string, pos int) (interface{}, error) {

	respValue, err := c.httpCall(http.MethodGet, c.url + "values/" + key + "?pos=" + strconv.Itoa(pos),  nil)

	return respValue["value"], err
}

func (c *Client) GetDictValue(key string, dictKey string) (interface{}, error) {

	respValue, err := c.httpCall(http.MethodGet, c.url + "values/" + key + "?key=" + dictKey,  nil)

	return respValue["value"], err
}

func (c *Client) Put(key string, value interface{}) (interface{}, error) {

	return c.httpCall(http.MethodPut, c.url + "values/" + key,  value)
}

func (c *Client) PutWithExpire(key string, value interface{}, expire int) (interface{}, error) {

	return c.httpCall(http.MethodPut, c.url + "values/" + key + "?expire=" + strconv.Itoa(expire),  value)
}

func (c *Client) Expire(key string, expire int) (interface{}, error) {

	return c.httpCall(http.MethodPut, c.url + "expire/" + key,  expire)
}

func (c *Client) Delete(key string) (interface{}, error) {

	respValue, err := c.httpCall(http.MethodDelete, c.url + "values/" + key, nil)

	return respValue["value"], err
}

func (c *Client) Persist() (interface{}, error) {

	respValue, err := c.httpCall(http.MethodPost, c.url + "persist", nil)

	return respValue["value"], err
}

func (c *Client) Reload() (interface{}, error) {

	respValue, err := c.httpCall(http.MethodPost, c.url + "reload", nil)

	return respValue["value"], err
}

func (c *Client) httpCall(method, url string, value interface{}) (map[string]interface{}, error) {

	var respValue map[string]interface{}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(value)

	req, err := http.NewRequest(method, url,  b)
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