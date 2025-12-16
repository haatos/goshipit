package handler

import (
	"log"
	"net/http"

	"github.com/haatos/goshipit/internal"
	"github.com/haatos/goshipit/internal/assets"
	"github.com/haatos/goshipit/internal/markdown"
	"github.com/haatos/goshipit/internal/views/pages"
	"github.com/labstack/echo/v4"
)

var gettingStartedHTML string
var typesHTML string
var cliHTML string

const (
	contentTypesMarkdownPath   = "content/types.md"
	gettingStartedMarkdownPath = "content/getting_started.md"
	cliMarkdownPath            = "content/cli.md"
)

func getGettingStartedHTML() {
	if gettingStartedHTML != "" {
		return
	}

	pageContent, err := assets.ContentFS.ReadFile(gettingStartedMarkdownPath)
	if err != nil {
		log.Fatal(err)
	}

	gettingStartedHTML = markdown.GetHTMLFromMarkdown(pageContent)
}

func getTypesHTML() {
	if typesHTML != "" {
		return
	}

	pageContent, err := assets.ContentFS.ReadFile(contentTypesMarkdownPath)
	if err != nil {
		log.Fatal(err)
	}

	typesHTML = markdown.GetHTMLFromMarkdown(pageContent)
}

func getCLIHTML() {
	if cliHTML != "" {
		return
	}

	pageContent, err := assets.ContentFS.ReadFile(cliMarkdownPath)
	if err != nil {
		log.Fatal(err)
	}
	cliHTML = markdown.GetHTMLFromMarkdown(pageContent)
}

func GetIndexPage(c echo.Context) error {
	if isHXRequest(c) {
		return render(c, http.StatusOK, pages.IndexPageContent())
	}
	return render(c, http.StatusOK, pages.IndexPage())
}

func GetAboutPage(c echo.Context) error {
	if isHXRequest(c) {
		return render(c, http.StatusOK, pages.AboutPageMain())
	}
	return render(c, http.StatusOK, pages.AboutPage())
}

func GetGettingStartedPage(c echo.Context) error {
	getGettingStartedHTML()

	if isHXRequest(c) {
		return render(c, http.StatusOK, pages.GettingStartedPageMain(gettingStartedHTML))
	}
	return render(c, http.StatusOK, pages.GettingStartedPage(gettingStartedHTML))
}

func GetTypesPage(c echo.Context) error {
	getTypesHTML()

	if isHXRequest(c) {
		return render(c, http.StatusOK, pages.TypesPageMain(typesHTML))
	}
	return render(c, http.StatusOK, pages.TypesPage(typesHTML))
}

func GetCLIPage(c echo.Context) error {
	getCLIHTML()

	if isHXRequest(c) {
		return render(c, http.StatusOK, pages.CLIPageMain(cliHTML))
	}
	return render(c, http.StatusOK, pages.CLIPage(cliHTML))
}

func GetPrivacyPolicyPage(c echo.Context) error {
	if isHXRequest(c) {
		return render(
			c, http.StatusOK,
			pages.PrivacyMain(internal.Settings.Domain, internal.Settings.ContactEmail))
	}
	return render(
		c, http.StatusOK,
		pages.PrivacyPage(internal.Settings.Domain, internal.Settings.ContactEmail))
}

func GetTermsOfServicePage(c echo.Context) error {
	if isHXRequest(c) {
		return render(
			c, http.StatusOK,
			pages.TermsOfServiceMain(internal.Settings.Domain, internal.Settings.ContactEmail))
	}
	return render(
		c, http.StatusOK,
		pages.TermsOfService(internal.Settings.Domain, internal.Settings.ContactEmail))
}
