package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/haatos/goshipit/internal/views/components"
	"github.com/labstack/echo/v4"
)

var emailRegexp = regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)

var StringValidations = map[string]func(data string) string{
	"notempty": func(data string) string {
		data = strings.TrimSpace(data)
		if data == "" {
			return "must not be empty"
		}
		return ""
	},
	"email": func(data string) string {
		data = strings.TrimSpace(data)
		if !emailRegexp.Match([]byte(data)) {
			return "must be valid"
		}
		return ""
	},
	"hasupper": func(data string) string {
		for _, r := range data {
			if unicode.IsUpper(r) {
				return ""
			}
		}
		return "must contain a uppercase letter"
	},
	"haslower": func(data string) string {
		for _, r := range data {
			if unicode.IsLower(r) {
				return ""
			}
		}
		return "must contain a lowercase letter"
	},
	"hasdigit": func(data string) string {
		for _, r := range data {
			if unicode.IsNumber(r) {
				return ""
			}
		}
		return "must contain a number"
	},
	"hasspecial": func(data string) string {
		chars := `§!"@#£¤$%&/{([=?+\'*<>,;.:-_])}`
		if strings.ContainsAny(data, chars) {
			return ""
		}
		return fmt.Sprintf("must contain one of %s", chars)
	},
}

func PostValidateString(c echo.Context) error {
	name := c.Param("name")
	value := c.FormValue(name)
	validations := c.QueryParams()["v"]

	fmt.Println("received", name, value, validations)

	errors := make([]string, 0, len(validations))
	for _, validation := range validations {
		if res := StringValidations[validation](value); res != "" {
			errors = append(errors, res)
		}
	}

	e := strings.Join(errors, ", ")

	if e != "" {
		return render(c, http.StatusUnprocessableEntity, components.PlainText(e))
	}
	return nil
}
