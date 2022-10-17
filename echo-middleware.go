package httperror

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// HttpErrorHandleMiddleware 错误处理
func HttpErrorHandleMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if err = next(c); err != nil {
				v, ok := err.(validator.ValidationErrors)
				if ok {
					err = BadRequest(v.Error(), 4000)
				}
				return
			}
			return
		}
	}
}
