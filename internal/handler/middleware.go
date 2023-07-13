package handler

import (
	"github.com/r-mol/balanser_highload_system/internal/balancer"

	"github.com/labstack/echo/v4"
)

func MiddlewareBalancer(lb *balancer.LoadBalancer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()

			lb.ServeHTTP(res, req)

			return
		}
	}
}
