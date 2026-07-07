package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/haatos/goshipit/internal/views/pages"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	c.Logger().Errorf("Handler error: %+v\n", err)
	switch e := err.(type) {
	case ErrorToast:
		if err := renderErrorConfirm(c, e.Status, e.Messages); err != nil {
			slog.ErrorContext(
				c.Request().Context(),
				"unable to render error confirm",
				slog.Any("error", err),
			)
		}
	case *echo.HTTPError:
		switch e.Code {
		case http.StatusNotFound:
			if err := render(c, e.Code, pages.NotFound()); err != nil {
				slog.ErrorContext(
					c.Request().Context(),
					"unable to render error",
					slog.Any("error", err),
				)
			}
		case http.StatusInternalServerError:
			if err := render(c, e.Code, pages.InternalServerError()); err != nil {
				slog.ErrorContext(
					c.Request().Context(),
					"unable to render 500",
					slog.Any("error", err),
				)
			}
		case http.StatusForbidden:
			if err := render(
				c, e.Code,
				pages.Forbidden("Invalid permissions to view this page."),
			); err != nil {
				slog.ErrorContext(
					c.Request().Context(),
					"unable to render 403",
					slog.Any("error", err),
				)
			}
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
