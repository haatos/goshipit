```go
package model

import (
	"github.com/a-h/templ"
)

type Anchor struct {
	Href  string
	Label string
	Icon  templ.Component
	Attrs templ.Attributes
}

type Banner struct {
	Title                 templ.Component
	Description           string
	CallToAction          Button
	SecondaryCallToAction Button
}

type Button struct {
	Label string
	Attrs templ.Attributes
}

type Card struct {
	Title   string
	Content string
	Source  string
	Alt     string
}

type Checkbox struct {
	Label   string
	Name    string
	Checked bool
	Class   string
	Attrs   templ.Attributes
}

type Combobox struct {
	Label    string
	Name     string
	URL      string
	Options  []string
	Selected []string
}

type CompanyInfo struct {
	Icon        templ.Component
	Name        string
	Description string
	Copyright   string
}

type DropdownItem struct {
	Label string
	Attrs templ.Attributes
}

type Feature struct {
	Icon        templ.Component
	Title       string
	Description string
	URL         string
}

type Image struct {
	Source string
	Alt    string
}

type Input struct {
	Label   string
	Name    string
	Value   string
	Err     string
	Small   bool
	Attrs   templ.Attributes
	Classes string
	Icon    templ.Component
}

type PaginationItem struct {
	URL      string
	Page     int
	Low      int
	High     int
	MaxPages int
}

type Price struct {
	Title            string
	Description      string
	Price            string
	Per              string
	IncludedFeatures []string
	ExcludedFeatures []string
	CallToAction     Button
	Footer           templ.Component
}

type Range struct {
	Label string
	Name  string
	Value int
	Min   int
	Max   int
	Step  int
	Class string
}

type Rating struct {
	Name  string
	Min   int
	Max   int
	Class string
	Value int
}

type Script struct {
	Source string
	Defer  bool
}

type Section struct {
	Title        string
	Description  string
	Source       string
	Alt          string
	CallToAction Button
	Reverse      bool
}

type Select struct {
	Label   string
	Name    string
	Options []SelectOption
	Attrs   templ.Attributes
}

type SelectOption struct {
	Label    string
	Value    string
	Selected bool
	Disabled bool
}

type Stat struct {
	Title       string
	Value       string
	Description string
}

type Status struct {
	Code         int
	Title        string
	Description  string
	ReturnButton Button
}

type Tab struct {
	Label   string
	Content templ.Component
}

type Testimonial struct {
	Avatar  templ.Component
	Name    string
	Rating  int
	Content string
}

type Textarea struct {
	Label string
	Name  string
	Value string
	Rows  int
	Err   string
	Class string
	Attrs templ.Attributes
}

type TimelineItem struct {
	Start  string
	Middle templ.Component
	End    string
}

type Toast struct {
	Name         string
	ToastClasses string
	AlertClasses string
}

type Toggle struct {
	Before    string
	After     string
	Name      string
	Checked   bool
	Class     string
	Highlight bool
	Attrs     templ.Attributes
}
```
