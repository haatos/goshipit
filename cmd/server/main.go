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
	e.GET("/search", handler.GetComponentSearch)

	/*
		handlers for component examples
	*/
	e.POST("/validate/string/:name", handler.PostValidateString)
	e.GET("/infinite-scroll", handler.GetInfiniteScroll)
	e.GET("/infinite-scroll-rows", handler.GetInfiniteScrollRows)
	e.GET("/active-search-table", handler.GetActiveSearchTable)
	e.GET("/active-search", handler.GetActiveSearch)
	e.GET("/lazy-load", handler.GetLazyLoad)
	e.GET("/pricing", handler.GetPricing)
	e.GET("/models", handler.GetModels)
	e.GET("/pagination-pages", handler.GetPaginationPage)
	/**/

	internal.GracefulShutdown(e, internal.Settings.Port)
}
