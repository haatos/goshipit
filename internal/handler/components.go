package handler

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/haatos/goshipit/generated"
	"github.com/haatos/goshipit/internal"
	"github.com/haatos/goshipit/internal/markdown"
	"github.com/haatos/goshipit/internal/model"
	"github.com/haatos/goshipit/internal/views/components"
	"github.com/haatos/goshipit/internal/views/examples"
	"github.com/haatos/goshipit/internal/views/pages"
	"github.com/labstack/echo/v4"
)

func GetComponentAnchors(c echo.Context) error {
	keys := make([]string, 0, len(internal.ComponentCodeMap))
	for k := range internal.ComponentCodeMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return render(c, http.StatusOK, pages.ComponentAnchors(keys, internal.ComponentCodeMap))
}

func GetComponentPage(c echo.Context) error {
	category := c.Param("category")
	name := c.Param("name")

	var componentCode model.ComponentCode
	for i := range internal.ComponentCodeMap[category] {
		if internal.ComponentCodeMap[category][i].Name == name {
			componentCode = internal.ComponentCodeMap[category][i]
		}
	}

	exampleCodes := internal.ComponentExampleCodeMap[name]
	coms := make([]templ.Component, 0, len(exampleCodes))
	for i, exampleCode := range exampleCodes {
		tabs := []model.Tab{
			{
				Label:   "Preview",
				Content: templ.Raw(getHTMLFromComponent(generated.ExampleComponents[exampleCode.Name])),
			},
			{
				Label:   "Templ",
				Content: templ.Raw(markdown.GetHTMLFromMarkdown([]byte(exampleCode.CodeMarkdown))),
			},
		}
		if exampleCode.HandlerMarkdown != "" {
			tabs = append(
				tabs,
				model.Tab{
					Label:   "Handler",
					Content: templ.Raw(markdown.GetHTMLFromMarkdown([]byte(exampleCode.HandlerMarkdown))),
				})
		}
		coms = append(coms, components.Tabs(fmt.Sprintf("%s-%d-tab", exampleCode.Name, i), tabs))
	}

	if isHXRequest(c) {
		return render(c, http.StatusOK, pages.ComponentMain(componentCode.Label, componentCode.CodeMarkdown, coms))
	}
	return render(c, http.StatusOK, pages.ComponentPage(componentCode.Label, componentCode.CodeMarkdown, coms))
}

func GetComponentSearch(c echo.Context) error {
	search := c.FormValue("search")
	search = strings.ToLower(search)
	search = strings.TrimSpace(search)
	if search == "" {
		return nil
	}

	results := []pages.ComponentSearchItem{}
	for category := range internal.ComponentCodeMap {
		for i := range internal.ComponentCodeMap[category] {
			if strings.Contains(
				strings.ToLower(internal.SnakeCaseToCapitalized(internal.ComponentCodeMap[category][i].Name)), search,
			) {
				results = append(results, pages.ComponentSearchItem{
					Category: category,
					Name:     internal.ComponentCodeMap[category][i].Name,
				})
			}
		}
	}

	return render(c, http.StatusOK, pages.ComponentSearchListItems(results))
}

// ActiveSearchExampleTable
func GetActiveSearchExample(c echo.Context) error {
	time.Sleep(500 * time.Millisecond)

	search := strings.ToLower(strings.TrimSpace(c.FormValue("search")))

	out := make([]templ.Component, 0, len(examples.ActiveSearchTableData))

	for _, rowData := range examples.ActiveSearchTableData {
		if search == "" || (strings.Contains(strings.ToLower(rowData.FirstName), search) ||
			strings.Contains(strings.ToLower(rowData.LastName), search) ||
			strings.Contains(strings.ToLower(rowData.Email), search)) {
			out = append(out, examples.ActiveSearchTableRow(
				rowData.FirstName, rowData.LastName, rowData.Email))
		}
	}

	return render(c, http.StatusOK, examples.ActiveSearchTableRows(out))
}

func GetActiveSearchExampleTable(c echo.Context) error {
	out := make([]templ.Component, 0, len(examples.ActiveSearchTableData))
	for _, rowData := range examples.ActiveSearchTableData {
		out = append(out, examples.ActiveSearchTableRow(
			rowData.FirstName, rowData.LastName, rowData.Email))
	}

	com := examples.ActiveSearchExample(
		"active-search-table",
		[]templ.Component{
			components.PlainText("First Name"),
			components.PlainText("Last Name"),
			components.PlainText("Email"),
		},
		out)

	return render(c, http.StatusOK, com)
}

// ActiveSearchExampleTable

// InfiniteScrollTableExample
func GetInfiniteScrollExample(c echo.Context) error {
	// generate dummy row data
	data := make([][]string, 0, 10)
	for i := 0; i < 10; i++ {
		data = append(data, []string{"John Doe", fmt.Sprintf("john.doe%d@email.com", i)})
	}

	hasMore := true

	rows := make([]templ.Component, 0, len(data))
	for i := range data {
		rows = append(rows, components.InfiniteScrollRow(data[i][0], data[i][1], 1, i+1 == 10 && hasMore))
	}

	return render(c, http.StatusOK, components.InfiniteScrollTable(rows))
}

func GetInfiniteScrollExampleRows(c echo.Context) error {
	// short delay for loading indicator
	time.Sleep(500 * time.Millisecond)

	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return newErrorToast(http.StatusUnprocessableEntity, fmt.Sprintf("invalid page '%s'", pageStr))
	}

	// generate dummy date for the page
	data := make([][]string, 0, 10)
	for i := (page - 1) * 10; i < (page-1)*10+10; i++ {
		data = append(data, []string{"John Doe", fmt.Sprintf("john.doe%d@email.com", i)})
	}

	hasMore := len(data) == 10 && page < 10

	rows := make([]templ.Component, 0, len(data))
	for i := range data {
		rows = append(rows, components.InfiniteScrollRow(data[i][0], data[i][1], page, i+1 == 10 && hasMore))
	}

	return render(c, http.StatusOK, components.InfiniteScrollRows(rows))
}

// InfiniteScrollTableExample

// LazyLoadExample
func GetLazyLoadExample(c echo.Context) error {
	time.Sleep(2 * time.Second)

	return render(c, http.StatusOK, examples.LazyLoadResult())
}

// LazyLoadExample

// PricingExample
func GetPricingExample(c echo.Context) error {
	yearly := c.QueryParam("period") == "on"

	return render(c, http.StatusOK, components.PriceGrid(examples.PriceDataExample(yearly)))
}

// PricingExample

// CascadingSelect
var modelOptions = map[string][]model.SelectOption{
	"audi": {
		{Label: "A1", Value: "a1"},
		{Label: "A4", Value: "a4"},
		{Label: "A6", Value: "a6"},
	},
	"bmw": {
		{Label: "316i", Value: "316i"},
		{Label: "320d", Value: "320d"},
	},
	"toyota": {
		{Label: "Corolla", Value: "corolla"},
		{Label: "Yaris", Value: "yaris"},
		{Label: "RAV4", Value: "rav4"},
	},
}

func GetCascadingSelectExample(c echo.Context) error {
	make := c.FormValue("make")

	return render(c, http.StatusOK, components.SelectOptions(modelOptions[make]))
}

// CascadingSelect

// BasicPaginationExample
const (
	ItemsPerPage = 10
	PagesPerSide = 3
)

func GetPaginationExamplePage(c echo.Context) error {
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return newErrorToast(http.StatusUnprocessableEntity, fmt.Sprintf("invalid page '%s'", pageStr))
	}

	// generate dummy data for the page
	data := make([][]string, 0, 200)
	for i := range 200 {
		data = append(data, []string{fmt.Sprintf("Key%d", i), fmt.Sprintf("Value%d", i)})
	}
	out := data[(page-1)*ItemsPerPage : (page-1)*ItemsPerPage+ItemsPerPage]

	maxPages := int(math.Ceil(200 / ItemsPerPage))
	low, high := getPaginationLowAndHigh(page, maxPages, PagesPerSide)

	com := examples.BasicPagination(
		"pagination-example-1",
		model.PaginationItem{
			URL:      c.Request().URL.Path,
			Page:     page,
			Low:      low,
			High:     high,
			MaxPages: maxPages,
		},
		out)

	return render(c, http.StatusOK, com)
}

func getPaginationLowAndHigh(page, maxPages, pagePerSide int) (int, int) {
	low := max(0, page-pagePerSide-1)
	high := min(maxPages-1, page+pagePerSide-1)

	// if there are less than 'pagePerSide' pages to the left
	// of the current page, we add more pages to the high end of the
	// pagination element by adding to 'high'
	add := pagePerSide - page
	if add >= 0 {
		high += add + 1
	}

	// if there are less than 'pagePerSide' pages to the right
	// of the current page, we add more pages to the low end of the
	// pagination element, by subtracting from 'low'
	distanceFromMaxPages := maxPages - page
	if distanceFromMaxPages < pagePerSide {
		low -= (pagePerSide - distanceFromMaxPages)
	}
	return low, high
}

// BasicPaginationExample

// BasicCombobox
func PostCombobox(c echo.Context) error {
	name := c.Param("name")
	value := c.Param("value")

	return render(c, http.StatusOK, components.ComboBadge(name, value))
}

func PostComboboxSubmit(c echo.Context) error {
	urlValues, err := c.FormParams()
	if err != nil {
		return newErrorToast(http.StatusUnprocessableEntity, "unable to parse form data")
	}

	name := c.Param("name")
	comboboxPostData := urlValues[name]

	return renderSuccessFade(c, http.StatusOK, comboboxPostData)
}

// BasicCombobox
