// Package curl 封装go 标准http 库使其支持类似curl 的调用
package curl

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

//定义一些基本的curl 请求方法常量映射
const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	PATCH   = "PATCH"
	DELETE  = "DELETE"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
)

//Client curl client 对象
type Client struct {
	engine  *http.Client
	baseURL *url.URL
}

//NewClient 新建一个 curl 客户端
func NewClient() *Client {
	c := &Client{
		engine: &http.Client{},
	}
	return c
}

//WithBaseURL 设置client 的 base url
func (c *Client) WithBaseURL(baseurl string) error {
	urlobj, err := url.Parse(baseurl)
	if err != nil {
		return err
	}
	if urlobj.Path != "" && urlobj.Path[len(urlobj.Path)-1] != '/' {
		urlobj.Path = urlobj.Path + "/"
	}
	//TODO: modify RawPath
	c.baseURL = urlobj
	return nil
}

//Fetch 执行request 请求
func (c *Client) Fetch(request *Request) (*Response, error) {
	httpRequest := *request.Raw
	if c.baseURL != nil {
		finalURL := c.baseURL.ResolveReference(httpRequest.URL)
		httpRequest.URL = finalURL
	}
	httpResponse, err := c.engine.Do(&httpRequest)
	// TODO: 完善err 处理
	if err != nil {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	response := &Response{
		Raw:  httpResponse,
		Body: responseBody,
	}
	return response, nil
}
