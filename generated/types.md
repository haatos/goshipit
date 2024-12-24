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

type Anchor struct {
	Href  string
	Label string
	Icon  templ.Component
	Attrs templ.Attributes
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
	end = end.AddDate(0, 0, 7-int(end.Weekday())-1+int(dp.StartOfWeek))
	if start.Weekday() == time.Sunday {
		start = start.AddDate(0, 0, -6)
	} else {
		start = start.AddDate(0, 0, -1*int(start.Weekday())+int(dp.StartOfWeek))
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
