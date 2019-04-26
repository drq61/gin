package gin

import (
	"io/ioutil"
	"strings"

	"github.com/wule61/gin/binding"

	"github.com/micro-plat/lib4go/encoding"
)

func makeFormData(ctx *Context) InputData {
	if ctx.ContentType() == binding.MIMEPOSTForm ||
		ctx.ContentType() == binding.MIMEMultipartPOSTForm {
		ctx.request.ParseForm()
		ctx.request.ParseMultipartForm(32 << 20)
	}

	return ctx.GetPostForm
}
func makeQueyStringData(ctx *Context) InputData {
	return ctx.GetQuery
}
func makeParamsData(ctx *Context) InputData {
	return ctx.Params.Get
}

func makeMapData(m map[string]interface{}) MapData {
	return m
}

func makeSettingData(ctx *Context, m map[string]string) ParamData {
	return m
}

func makeExtData(c *Context) map[string]interface{} {
	input := make(map[string]interface{})
	input["__method_"] = strings.ToLower(c.request.Method)
	input["__header_"] = c.request.Header
	input["__func_http_request_"] = c.Request
	input["__func_http_response_"] = c.Writer
	input["__binding_"] = c.ShouldBind
	input["__binding_with_"] = func(v interface{}, ct string) error {
		return c.BindWith(v, binding.Default(c.request.Method, ct))

	}
	input["__get_request_values_"] = func() map[string]interface{} {
		c.request.ParseForm()
		data := make(map[string]interface{})
		query := c.request.URL.Query()
		for k, v := range query {
			switch len(v) {
			case 1:
				data[k] = v[0]
			default:
				data[k] = strings.Join(v, ",")
			}
		}
		forms := c.request.PostForm
		for k, v := range forms {
			switch len(v) {
			case 1:
				data[k] = v[0]
			default:
				data[k] = strings.Join(v, ",")
			}
		}

		return data
	}

	input["__func_body_get_"] = func(ch string) (string, error) {
		buff, err := ioutil.ReadAll(c.request.Body)
		if err != nil {
			return "", err
		}
		nbuff, err := encoding.DecodeBytes(buff, ch)
		if err != nil {
			return "", err
		}
		return string(nbuff), nil

	}
	return input
}

type MapData map[string]interface{}

//Get 获取指定键对应的数据
func (i MapData) Get(key string) (interface{}, bool) {
	r, ok := i[key]
	return r, ok
}

//InputData 输入参数
type InputData func(key string) (string, bool)

//Get 获取指定键对应的数据
func (i InputData) Get(key string) (interface{}, bool) {
	r, ok := i(key)
	return r, ok
}

//ParamData map参数数据
type ParamData map[string]string

//Get 获取指定键对应的数据
func (i ParamData) Get(key string) (interface{}, bool) {
	r, ok := i[key]
	return r, ok
}
