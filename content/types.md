# Types

## Introduction

Most components use a struct as the input argument. This provides a convenient way to pass default values to the components; struct fields initialize to the respective type's default value. We can use this feature to make some of the fields of a struct optional to make initialization less cumbersome. For example, we can pass the `components.AnchorProps` argument to an anchor element, `<a>`:

```go
package components

type AnchorProps struct {
	Href      string
	Label     string
	LeftIcon  templ.Component
	RightIcon templ.Component
	Attrs     templ.Attributes
	Class     string
}

templ Anchor(props AnchorProps) {
	<a
		if props.Href != "" {
			href={ templ.SafeURL(props.Href) }
		}
		class={ "group flex items-center cursor-pointer", props.Class }
		{ props.Attrs... }
	>
		if props.LeftIcon != nil {
			<div class="inline-block mr-1">
				@props.LeftIcon
			</div>
		}
		{ props.Label }
		if props.RightIcon != nil {
			<div class="inline-block ml-1">
				@props.RightIcon
			</div>
		}
	</a>
}
```

Here we can define the anchor element to optionally have an `href` attribute, or if we choose, a `hx-get` instead:

```go
@components.Anchor(components.AnchorProps{Href: "/"})
@components.Anchor(components.AnchorProps{Attrs: templ.Attributes{"hx-get": "/"}})
```

## Deprecated types

```go
package model

import (
	"time"

	"github.com/a-h/templ"
)

type Accordion struct {
	Label string
	Type  string
	Name  string
}

type ActiveSearchInput struct {
	ID     string
	URL    string
	Target string
	Input  Input
}

type Anchor struct {
	Href      string
	Label     string
	LeftIcon  templ.Component
	RightIcon templ.Component
	Attrs     templ.Attributes
	Class     string
}

type Avatar struct {
	AvatarClass      string
	ContainerClass   string
	Source           string
	Placeholder      string
	PlaceholderClass string
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
	Class   string
}

type Chat struct {
	Messages []ChatMessage
}

type ChatMessage struct {
	AvatarURL string
	Sender    string
	Time      string
	Message   string
	Footer    string
	Location  string
	Class     string
}

type Checkbox struct {
	ID      string
	Before  string
	After   string
	Name    string
	Checked bool
	Class   string
	Attrs   templ.Attributes
}

type Collapse struct {
	Class        string
	Title        string
	TitleClass   string
	ContentClass string
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

type DatePicker struct {
	Year        int
	Month       int
	Selected    time.Time
	StartOfWeek time.Weekday
}

func (dp DatePicker) Days() []time.Time {
	days := make([]time.Time, 0, 31)
	now := time.Now().UTC()
	start := time.Date(dp.Year, time.Month(dp.Month), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, -1)
	for end.Weekday() != dp.StartOfWeek {
		end = end.AddDate(0, 0, 1)
	}
	end = end.AddDate(0, 0, -1)

	for start.Weekday() != dp.StartOfWeek {
		start = start.AddDate(0, 0, -1)
	}
	for !start.After(end) {
		days = append(days, start)
		start = start.AddDate(0, 0, 1)
	}
	return days
}

func (dp DatePicker) Months() []time.Time {
	months := make([]time.Time, 12)
	for i := 1; i <= 12; i++ {
		dt := time.Date(dp.Year, time.Month(i), 1, 0, 0, 0, 0, time.Now().Location())
		months[i-1] = dt
	}
	return months
}

type Dropdown struct {
	Label     string
	Class     string
	ListClass string
	Items     []DropdownItem
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
	ID              string
	Type            string // defaults to "text"
	Label           string
	Name            string
	Value           string
	Placeholder     string
	Err             string
	Attrs           templ.Attributes
	Class           string
	Icon            templ.Component
	Disabled        bool
	DisabledMessage string
	Required        bool
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

type Radio struct {
	Name   string
	Values map[string]string
	Class  string
}

type Range struct {
	ID    string
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

type Select struct {
	ID      string
	Label   string
	Name    string
	Options []SelectOption
	Attrs   templ.Attributes
	Class   string
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

type Swap struct {
	On    templ.Component
	Off   templ.Component
	Class string
}

type Tabs struct {
	Name         string
	Class        string
	Tabs         []Tab
	ContentClass string
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
	ID          string
	Label       string
	Name        string
	Placeholder string
	Value       string
	Rows        int
	Err         string
	Class       string
	Attrs       templ.Attributes
}

type TimelineItem struct {
	Start  string
	Middle templ.Component
	End    string
}

type Toast struct {
	Name       string
	ToastClass string
	AlertClass string
}

type Toggle struct {
	ID        string
	Before    string
	After     string
	Name      string
	Checked   bool
	Class     string
	Highlight bool
	Attrs     templ.Attributes
}

type Tooltip struct {
	Tip   string
	Class string
}
```
