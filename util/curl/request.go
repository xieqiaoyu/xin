package curl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//Request Request
type Request struct {
	Raw *http.Request
}

//NewRequest 新建一个curl 请求
func NewRequest(method, url string) (*Request, error) {
	rawRequest, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	r := &Request{
		Raw: rawRequest,
	}
	return r, nil
}

//WithJSONBody 添加一个json 请求体
func (r *Request) WithJSONBody(content interface{}) error {
	bodyStream, err := json.Marshal(content)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(bodyStream)
	r.Raw.ContentLength = int64(bodyReader.Len())
	snapshot := *bodyReader
	r.Raw.Body = ioutil.NopCloser(bodyReader)
	r.Raw.GetBody = func() (io.ReadCloser, error) {
		r := snapshot
		return ioutil.NopCloser(&r), nil
	}
	r.SetHeader("Content-Type", "application/json")
	return nil
}

//SetHeader 设置请求头
func (r *Request) SetHeader(key, value string) {
	r.Raw.Header.Set(key, value)
}

//WithBasicAuth 添加basic auth 认证
func (r *Request) WithBasicAuth(username, passwd string) error {
	r.Raw.SetBasicAuth(username, passwd)
	return nil
}

// WithBearerAuth 添加 bearer auth 认证
func (r *Request) WithBearerAuth(token string) error {
	r.SetHeader("Authorization", "Bearer "+token)
	return nil
}

//WithQuery 在request 中添加query
func (r *Request) WithQuery(v *url.Values) error {
	r.Raw.URL.RawQuery = v.Encode()
	return nil
}

//String request stringfy
func (r *Request) String() string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Raw.Method, r.Raw.URL, r.Raw.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Raw.Host))
	// Loop through headers
	for name, headers := range r.Raw.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}
	if r.Raw.ContentLength > 0 {
		b, _ := ioutil.ReadAll(r.Raw.Body)
		request = append(request, "\n")
		request = append(request, string(b))
	}

	// Return the request as a string
	return strings.Join(request, "\n")
}
