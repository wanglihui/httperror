package httperror

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

// EchoHttpErrorHandler echo httpErrorhandler 兼容 自定义 httperror
func EchoHttpErrorHandler(server *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		// fixedbug echo 错误处理会被调用2次，导致错误结果会被输出2次
		// https://github.com/globocom/echo-prometheus/issues/5
		if c.Response().Committed {
			return
		}
		if v, ok := err.(validator.ValidationErrors); ok {
			err = BadRequest(v.Error(), 4000)
		}
		if _, ok := err.(*HTTPError); !ok {
			httpcode := gjson.Get(err.Error(), "http_code").Int()
			msg := gjson.Get(err.Error(), "msg").String()
			code := gjson.Get(err.Error(), "code").Int()
			if httpcode > 0 && msg != "" {
				err = New(int(httpcode), msg, code)
			}
		}
		if v, ok := err.(*HTTPError); ok {
			if err = c.JSON(v.StatusCode, v); err != nil {
				server.Logger.Error(err)
				if !c.Response().Committed {
					server.DefaultHTTPErrorHandler(err, c)
				}
			}
		} else {
			server.DefaultHTTPErrorHandler(err, c)
		}
	}
}
