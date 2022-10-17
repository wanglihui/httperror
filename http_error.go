package httperror

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// HTTPError 自定义HTTP 错误
type HTTPError struct {
	error
	// Http 状态码
	StatusCode int `json:"http_code"`
	// 错误消息
	Message string `json:"msg"`
	// 业务错误码
	Code int64 `json:"code"`
	// IsHttpError 标记
	IsHTTPError bool `json:"-"`
}

func (that *HTTPError) Error() string {
	return fmt.Sprintf(`{"http_code": %d, "msg": "%s", "code": %d}`, that.StatusCode, that.Message, that.Code)
}

// IsHTTPError 判断是否是 HTTPError 类型
func IsHTTPError(e interface{}) bool {
	_, ok := e.(*HTTPError)
	return ok
}

func Parse(e error, code int64) *HTTPError {
	if err, ok := e.(*HTTPError); ok {
		return err
	}
	return InternalError(e.Error(), code)
}

// New 生成 HTTPError
func New(statusCode int, message string, code int64) *HTTPError {
	if statusCode > 600 {
		panic("statusCode不能大于600")
	}
	return &HTTPError{
		error:       fmt.Errorf("%d:%d-%s", statusCode, code, message),
		StatusCode:  statusCode,
		Message:     message,
		Code:        code,
		IsHTTPError: true,
	}
}

// TryConstructHttpError 尝试使用字符串或者httperror.HttpError 构建httperror
func TryConstructHttpError(err interface{}) (*HTTPError, bool) {
	b := true
	if e, ok := err.(*HTTPError); ok {
		return e, b
	}
	if e, ok := err.(string); ok {
		msg := gjson.Get(e, "msg").String()
		errCode := gjson.Get(e, "code").Int()
		errStatusCode := gjson.Get(e, "http_code").Int()
		if errStatusCode > 0 && msg != "" {
			return New(int(errStatusCode), msg, errCode), b
		}
	}
	b = false
	return nil, b
}

// BadRequest 400 参数校验失败
func BadRequest(message string, code int64) *HTTPError {
	return New(400, message, code)
}

// InternalError 500 系统内部错误
func InternalError(message string, code int64) *HTTPError {
	return New(500, message, code)
}

// RequestNotAccept 406 请求不可接受
func RequestNotAccept(message string, code int64) *HTTPError {
	return New(406, message, code)
}
