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
	examples := make([]model.ComponentExample, 0, len(exampleCodes))
	for _, exampleCode := range exampleCodes {
		examples = append(examples, model.ComponentExample{
			HTML: getHTMLFromComponent(generated.ExampleComponents[exampleCode.Name]),
			Code: markdown.GetHTMLFromMarkdown([]byte(exampleCode.Markdown)),
		})
	}

	if c.Request().Header.Get("hx-request") == "" {
		return render(c, http.StatusOK, pages.ComponentPage(componentCode.Label, componentCode.Markdown, examples))
	} else {
		return render(c, http.StatusOK, pages.ComponentMain(componentCode.Label, componentCode.Markdown, examples))
	}
}

func GetComponentSearch(c echo.Context) error {
	return render(c, http.StatusOK, pages.ComponentSearchListItems(nil))
}

func GetActiveSearch(c echo.Context) error {
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

	com := examples.ActiveSearchTableRows(out)

	return render(c, http.StatusOK, com)
}

func GetActiveSearchTable(c echo.Context) error {
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

func GetInfiniteScroll(c echo.Context) error {
	/* generate dummy data */
	data := make([][]string, 0, 10)
	for i := 0; i < 10; i++ {
		data = append(data, []string{"John Doe", fmt.Sprintf("john.doe%d@email.com", i)})
	}

	hasMore := true

	rows := make([]templ.Component, 0, len(data))
	for i := range data {
		rows = append(rows, components.InfiniteScrollRow(data[i][0], data[i][1], 1, i+1 == 10 && hasMore))
	}
	/**/

	return render(c, http.StatusOK, components.InfiniteScrollTable(rows))
}

func GetInfiniteScrollRows(c echo.Context) error {
	// short delay for loading indicator
	time.Sleep(500 * time.Millisecond)

	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return newErrorToast(http.StatusUnprocessableEntity, fmt.Sprintf("invalid page '%s'", pageStr))
	}

	/*
		generate dummy date for the specified page
	*/
	data := make([][]string, 0, 10)
	for i := (page - 1) * 10; i < (page-1)*10+10; i++ {
		data = append(data, []string{"John Doe", fmt.Sprintf("john.doe%d@email.com", i)})
	}

	hasMore := len(data) == 10 && page < 10

	rows := make([]templ.Component, 0, len(data))
	for i := range data {
		rows = append(rows, components.InfiniteScrollRow(data[i][0], data[i][1], page, i+1 == 10 && hasMore))
	}
	/**/

	return render(c, http.StatusOK, components.InfiniteScrollRows(rows))
}

func GetLazyLoad(c echo.Context) error {
	time.Sleep(3 * time.Second)

	return render(c, http.StatusOK, examples.LazyLoadResult())
}

func GetPricing(c echo.Context) error {
	yearly := c.QueryParam("period") == "on"

	return render(c, http.StatusOK, components.PriceGrid(examples.PriceDataExample(yearly)))
}

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

func GetModels(c echo.Context) error {
	make := c.FormValue("make")

	return render(c, http.StatusOK, components.SelectOptions(modelOptions[make]))
}

const (
	ItemPerPage = 10
	PagePerSide = 3
)

func GetPaginationPage(c echo.Context) error {
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.Logger().Error("invalid page", pageStr)
	}

	// count rows from database to get maxPages
	//
	// use a limit offset query to get current page's items
	// e.g.
	// ...
	// order by created_on asc
	// limit <ItemPerPage> offset <ItemPerPage * page>
	data := make([][]string, 0, 200)
	for i := range 200 {
		data = append(data, []string{fmt.Sprintf("Key%d", i), fmt.Sprintf("Value%d", i)})
	}
	out := data[(page-1)*ItemPerPage : (page-1)*ItemPerPage+ItemPerPage]

	maxPages := int(math.Ceil(200 / ItemPerPage))
	low, high := getPaginationLowAndHigh(page, maxPages, PagePerSide)

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

	// check if we're close to the start of pagination
	// move the high end accordingly
	add := pagePerSide - page
	if add >= 0 {
		high += add + 1
	}

	// check if we're close to the end of pagination
	// move the low end accordingly
	distanceFromMaxPages := maxPages - page
	if distanceFromMaxPages < pagePerSide {
		low -= (pagePerSide - distanceFromMaxPages)
	}
	return low, high
}
