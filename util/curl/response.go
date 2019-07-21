package curl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//Response 返回对象
type Response struct {
	Raw  *http.Response
	Body []byte
}

func (r *Response) String() string {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "%s %s\n", r.Raw.Proto, r.Raw.Status)
	for key, value := range r.Raw.Header {
		fmt.Fprintf(b, "%s:%s\n", key, strings.Join(value, "; "))
	}
	fmt.Fprintf(b, "\n")
	b.Write(r.Body)
	return b.String()
}

//ParseJSONBody 将返回body 作为JSON 解析
func (r *Response) ParseJSONBody(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

//ParseBody 根据 response 的 content-type 自动解析返回结果
func (r *Response) ParseBody(v interface{}) error {
	contentType := filterFlags(r.Raw.Header.Get("Content-Type"))
	switch contentType {
	case "application/json":
		return r.ParseJSONBody(v)
	default:
		return fmt.Errorf("content-type %s is not support", contentType)
	}
}
