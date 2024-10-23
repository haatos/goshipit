package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/haatos/goshipit/internal"
	"github.com/haatos/goshipit/internal/markdown"
	"github.com/haatos/goshipit/internal/views/pages"
	"github.com/labstack/echo/v4"
)

var gettingStartedHTML string
var typesHTML string

func getGettingStartedHTML() {
	if gettingStartedHTML != "" {
		return
	}

	data, err := os.ReadFile("content/getting_started.md")
	if err != nil {
		log.Fatal(err)
	}

	gettingStartedHTML = markdown.GetHTMLFromMarkdown(data)
}

func getTypesHTML() {
	if typesHTML != "" {
		return
	}

	data, err := os.ReadFile("content/types.md")
	if err != nil {
		log.Fatal(err)
	}

	typesHTML = markdown.GetHTMLFromMarkdown(data)
}

func GetIndexPage(c echo.Context) error {
	return render(c, http.StatusOK, pages.IndexPage())
}

func GetAboutPage(c echo.Context) error {
	if c.Request().Header.Get("hx-request") == "" {
		return render(c, http.StatusOK, pages.AboutPage())
	} else {
		return render(c, http.StatusOK, pages.AboutPageMain())
	}
}

func GetGettingStartedPage(c echo.Context) error {
	getGettingStartedHTML()

	if c.Request().Header.Get("hx-request") == "" {
		return render(c, http.StatusOK, pages.GettingStartedPage(gettingStartedHTML))
	} else {
		return render(c, http.StatusOK, pages.GettingStartedPageMain(gettingStartedHTML))
	}
}

func GetTypesPage(c echo.Context) error {
	getTypesHTML()

	if c.Request().Header.Get("hx-request") == "" {
		return render(c, http.StatusOK, pages.TypesPage(typesHTML))
	} else {
		return render(c, http.StatusOK, pages.TypesPageMain(typesHTML))
	}
}

func GetPrivacyPolicyPage(c echo.Context) error {
	if c.Request().Header.Get("hx-request") == "" {
		return render(
			c, http.StatusOK,
			pages.PrivacyPage(internal.Settings.Domain, internal.Settings.ContactEmail))
	} else {
		return render(
			c, http.StatusOK,
			pages.PrivacyMain(internal.Settings.Domain, internal.Settings.ContactEmail))
	}
}

func GetTermsOfServicePage(c echo.Context) error {
	if c.Request().Header.Get("hx-request") == "" {
		return render(
			c, http.StatusOK,
			pages.TermsOfService(internal.Settings.Domain, internal.Settings.ContactEmail))
	} else {
		return render(
			c, http.StatusOK,
			pages.TermsOfServiceMain(internal.Settings.Domain, internal.Settings.ContactEmail))
	}
}
