/*
@Time : 31/1/2021 公元 11:05
@Author : philiphu
@File : request_client
@Software: GoLand
*/
package common

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type RequestClient struct {
	url        string //ip:port
	client *http.Client
}


func (c *RequestClient) Send(request *RequestClient) (response string, err error) {
	myClient := http.Client{Timeout: time.Second * 2}
	resp, err := myClient.Post(c.url, "application/x-www-form-urlencoded", strings.NewReader("your postdata"))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *RequestClient) Get() (response string, err error) {
	req, err := http.NewRequest(http.MethodGet, "http://"+c.url, nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	response = string(body)
	if "" == response {
		err = errors.New("query http://" + c.url + "return null.")
	}
	return
}

func (c *RequestClient) Init(url string) (*RequestClient, error) {
	c.client = &http.Client{
		Timeout: 3 * time.Second,
	}
	c.url = url
	return c, nil
}

func NewClient(url string) (client *RequestClient, err error) {
	c := &RequestClient{}
	client, err = c.Init(url)
	return
}
func (c *RequestClient) Close() {
	c.client.CloseIdleConnections()

}
