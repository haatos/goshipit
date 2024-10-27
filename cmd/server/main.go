package main

import (
	"github.com/haatos/goshipit/internal"
	"github.com/haatos/goshipit/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	internal.ReadDotenv()
	internal.Settings = internal.NewSettings()

	e := echo.New()
	e.HTTPErrorHandler = handler.ErrorHandler
	loggerFormat := "${method} ${uri} [${status}] (${latency_human}) | ${short_file}:${line} | ${message}\n"
	e.Logger.SetHeader(loggerFormat)

	config := internal.GetRateLimiterConfig()
	e.Use(middleware.RateLimiterWithConfig(config))

	e.Use(
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: loggerFormat,
		}),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: middleware.DefaultSkipper,
			Level:   3,
		}),
	)

	e.Static("/", "static")

	e.GET("/", handler.GetIndexPage)
	e.GET("/about", handler.GetAboutPage)
	e.GET("/get-started", handler.GetGettingStartedPage)
	e.GET("/types", handler.GetTypesPage)
	e.GET("/component-anchors", handler.GetComponentAnchors)
	e.GET("/privacy", handler.GetPrivacyPolicyPage)
	e.GET("/terms-of-service", handler.GetTermsOfServicePage)

	e.GET("/components/:category/:name", handler.GetComponentPage)
	e.GET("/components/search", handler.GetComponentSearch)

	// handlers for component examples
	e.POST("/validate/string/:name", handler.PostValidateString)
	e.GET("/infinite-scroll", handler.GetInfiniteScrollExample)
	e.GET("/infinite-scroll-rows", handler.GetInfiniteScrollExampleRows)
	e.GET("/active-search-table", handler.GetActiveSearchExampleTable)
	e.GET("/active-search", handler.GetActiveSearchExample)
	e.GET("/lazy-load", handler.GetLazyLoadExample)
	e.GET("/pricing", handler.GetPricingExample)
	e.GET("/models", handler.GetCascadingSelectExample)
	e.GET("/pagination-pages", handler.GetPaginationExamplePage)
	e.POST("/combobox/:name/:value", handler.PostCombobox)
	e.POST("/combobox-submit/:name", handler.PostComboboxSubmit)
	// handlers for component examples

	internal.GracefulShutdown(e, internal.Settings.Port)
}
