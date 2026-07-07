package handler

import (
	"github.com/a-h/templ"
	"github.com/haatos/goshipit/internal/views/custom"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, status int, com templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := com.Render(c.Request().Context(), buf); err != nil {
		return err
	}

	return c.HTML(status, buf.String())
}

func renderErrorConfirm(c echo.Context, status int, errs []string) error {
	hxRetarget(c, "body")
	hxReswap(c, "beforeend")
	return render(c, status, custom.ToastErrorConfirm(errs...))
}

func renderInfoFade(c echo.Context, status int, messages []string) error {
	hxRetarget(c, "body")
	hxReswap(c, "beforeend")
	return render(c, status, custom.HXToastInfoFade(messages...))
}
