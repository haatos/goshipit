package handler

import (
	"net/http"
	"strings"

	"github.com/haatos/goshipit/internal/views/pages"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	c.Logger().Errorf("Handler error: %+v\n", err)
	switch e := err.(type) {
	case ErrorToast:
		renderErrorConfirm(c, e.Status, e.Messages)
	case *echo.HTTPError:
		switch e.Code {
		case http.StatusNotFound:
			render(c, e.Code, pages.NotFound())
		case http.StatusInternalServerError:
			render(c, e.Code, pages.InternalServerError())
		case http.StatusForbidden:
			render(c, e.Code, pages.Forbidden("Invalid permissions to view this page."))
		}
	}
}

type ErrorToast struct {
	Status   int
	Messages []string
}

func (te ErrorToast) Error() string {
	return strings.Join(te.Messages, ", ")
}

func newErrorToast(status int, messages ...string) ErrorToast {
	return ErrorToast{
		Status:   status,
		Messages: messages,
	}
}

func NotFound(c echo.Context) error {
	return render(c, http.StatusNotFound, pages.NotFound())
}
