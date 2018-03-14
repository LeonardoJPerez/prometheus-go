package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo"
)

// RecordRequestMiddleware serves as a metering middleware for Echo Web framework.
func RecordRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		begin := time.Now()

		err := next(c)

		code := strconv.Itoa(c.Response().Status)
		method := c.Request().Method
		handler := c.Request().URL.Path

		AddRequestMetric(code, method, handler)
		ObserveRequestDuration(method, handler, begin)

		return err
	}
}
