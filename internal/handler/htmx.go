package handler

import "github.com/labstack/echo/v4"

func hxRedirect(c echo.Context, path string) {
	c.Response().Writer.Header().Set("HX-Redirect", path)
}

func hxReswap(c echo.Context, swap string) {
	c.Response().Writer.Header().Set("HX-Reswap", swap)
}
