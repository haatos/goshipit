package model

type ComponentCode struct {
	Name    string   `json:"name"`
	Code    []string `json:"code"`
	Handler []string `json:"handler,omitempty"`

	Label           string `json:"-"`
	CodeMarkdown    string `json:"-"`
	HandlerMarkdown string `json:"-"`
}

type ComponentCodeMap map[string][]ComponentCode

type ComponentExampleCodeMap map[string][]ComponentCode
