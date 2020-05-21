package testing

import (
	"fmt"
	"github.com/xieqiaoyu/xin/util/curl"
	"github.com/xieqiaoyu/xin/util/jsonschema"
	xyaml "github.com/xieqiaoyu/xin/util/yaml"
	"net/url"
	"strconv"
	"sync"
)

type HTTPTestCase interface {
	Run(ctx *HTTPTestContext) error
}

type HTTPTestContext struct {
	values            sync.Map
	defaultCurlClient *curl.Client
	extraCurlClient   map[string]*curl.Client
}

func NewHTTPTestContext(client *curl.Client) *HTTPTestContext {
	return &HTTPTestContext{
		values:            sync.Map{},
		defaultCurlClient: client,
		extraCurlClient:   map[string]*curl.Client{},
	}
}

func (c *HTTPTestContext) Set(key, value interface{}) {
	c.values.Store(key, value)
}
func (c *HTTPTestContext) Get(key interface{}) (value interface{}, ok bool) {
	return c.values.Load(key)
}
func (c *HTTPTestContext) GetString(key interface{}) (value string, ok bool) {
	v, ok := c.Get(key)
	if !ok {
		return "", false
	}
	vstring, ok := v.(string)
	if !ok {
		return "", false
	}
	return vstring, true
}

func (c *HTTPTestContext) GetInt(key interface{}) (value int, ok bool) {
	v, ok := c.Get(key)
	if !ok {
		return 0, false
	}
	var vint int
	var err error
	switch x := v.(type) {
	case int:
		vint = x
	case int32:
		vint = int(x)
	case int64:
		vint = int(x)
	case string:
		vint, err = strconv.Atoi(x)
		if err != nil {
			ok = false
		}
	case float64:
		vint = int(x)
	default:
		ok = false

	}
	return vint, ok
}

func (c *HTTPTestContext) FetchAndAssume(req *curl.Request, assume *ApiReturnAssume) (res *curl.Response, err error) {
	return fetchAndAssume(c.defaultCurlClient, req, assume)
}

func (c *HTTPTestContext) ExtraFetchAndAssume(key string, req *curl.Request, assume *ApiReturnAssume) (res *curl.Response, err error) {
	client, ok := c.extraCurlClient[key]
	if !ok {
		return nil, fmt.Errorf("fail to get curl client by key %s", key)
	}
	return fetchAndAssume(client, req, assume)
}

func (c *HTTPTestContext) AddExtraClient(key string, client *curl.Client) {
	c.extraCurlClient[key] = client
}

func fetchAndAssume(client *curl.Client, req *curl.Request, assume *ApiReturnAssume) (res *curl.Response, err error) {
	res, err = client.Fetch(req)
	if err != nil {
		return nil, err
	}
	if assume != nil {
		err = assume.Verify(res)
		if err != nil {
			return nil, fmt.Errorf("%s \nRequest:\n%s\n\nResponse:\n%s\n", err, req, res)
		}
	}
	return res, nil
}

type ApiReturnAssume struct {
	Httpstatus int
	BodySchema string
}

func (r *ApiReturnAssume) Status(s int) *ApiReturnAssume {
	r.Httpstatus = s
	return r
}

// very Body by jsonschema write in yaml
func (r *ApiReturnAssume) BodyVerifyBy(schema string) *ApiReturnAssume {
	jsondata, err := xyaml.Yaml2Json([]byte(schema))
	if err != nil {
		panic(err)
	}
	r.BodySchema = string(jsondata)
	return r
}

func (r *ApiReturnAssume) Verify(res *curl.Response) error {
	if r.Httpstatus != 0 {
		if res.Raw.StatusCode != r.Httpstatus {
			return fmt.Errorf("HTTP status is %d insteadof %d", res.Raw.StatusCode, r.Httpstatus)
		}
	}
	if r.BodySchema != "" {
		ok, err := jsonschema.ValidJSONString(string(res.Body), r.BodySchema)
		if !ok {
			return fmt.Errorf("return body err:%s", err)
		}
	}
	return nil
}

func ApiAssume() *ApiReturnAssume {
	return new(ApiReturnAssume)
}

type TestApiAction func(ctx *HTTPTestContext, req *curl.Request) error

type ApiTestCase struct {
	Method string
	Path   string
	Action TestApiAction
}

func (c *ApiTestCase) Run(ctx *HTTPTestContext) error {
	if c.Action != nil {
		req, err := curl.NewRequest(c.Method, c.Path)
		if err != nil {
			return err
		}
		err = c.Action(ctx, req)
		if err != nil {
			return fmt.Errorf("%s %s fail: %s", c.Method, c.Path, err)
		}
		return nil
	} else {
		return fmt.Errorf("No test action for case %s %s", c.Method, c.Path)
	}
}

func NewApiCase(method, path string, action TestApiAction) *ApiTestCase {
	if action == nil {
		action = apiDefautAction
	}
	return &ApiTestCase{
		Method: method,
		Path:   path,
		Action: action,
	}
}

func apiDefautAction(ctx *HTTPTestContext, req *curl.Request) error {
	assume := ApiAssume().Status(200)
	_, err := ctx.FetchAndAssume(req, assume)
	return err
}

type JobAction func(ctx *HTTPTestContext) error

type JobTestCase struct {
	action JobAction
}

func NewJobCase(action JobAction) *JobTestCase {
	return &JobTestCase{
		action: action,
	}
}

func (c *JobTestCase) Run(ctx *HTTPTestContext) error {
	if c.action != nil {
		return c.action(ctx)
	} else {
		return fmt.Errorf("No test action for case")
	}
}

func RenderReqPath(req *curl.Request, param interface{}) error {
	reqUrl := req.Raw.URL
	switch vt := param.(type) {
	case string:
		return RenderPath(reqUrl, vt)
	case int:
		vstring := strconv.Itoa(vt)
		return RenderPath(reqUrl, vstring)
	case map[string]interface{}:
		return MapRenderPath(reqUrl, vt)
	default:
		return fmt.Errorf("Unsupport param type %T", param)
	}
}

func RenderPath(raw *url.URL, vstring string) error {
	path := raw.Path
	iMax := len(path)
	for i := 0; i < iMax; {
		if path[i] == ':' {
			i++
			if i < iMax {
				k := i
				for i < iMax {
					if path[i] == '/' {
						break
					}
					i++
				}
				if vstring != "" {
					prefix := path[:k-1]
					suffix := path[i:]
					path = prefix + vstring + suffix
					lenInc := len(vstring) - i + k - 1
					iMax += lenInc
					// move pointer
					i += lenInc + 1
				}
			}
		} else {
			i++
		}
	}
	raw.Path = path
	return nil
}

func MapRenderPath(raw *url.URL, params map[string]interface{}) error {
	path := raw.Path
	iMax := len(path)
	for i := 0; i < iMax; {
		if path[i] == ':' {
			i++
			if i < iMax {
				k := i
				for i < iMax {
					if path[i] == '/' {
						break
					}
					i++
				}
				key := path[k:i]
				value, ok := params[key]
				if ok {
					var vstring string
					switch vt := value.(type) {
					case string:
						vstring = vt
					case int:
						vstring = strconv.Itoa(vt)
					default:
						break
					}
					if vstring != "" {
						prefix := path[:k-1]
						suffix := path[i:]
						path = prefix + vstring + suffix
						lenInc := len(vstring) - i + k - 1
						iMax += lenInc
						// move pointer
						i += lenInc + 1
					}
				}
			}
		} else {
			i++
		}
	}
	raw.Path = path
	return nil
}
