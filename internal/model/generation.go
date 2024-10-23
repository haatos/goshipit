package model

type ComponentCode struct {
	Name string   `json:"name"`
	Code []string `json:"code"`

	Label    string `json:"-"`
	Markdown string `json:"-"`
}

type ComponentExample struct {
	HTML string
	Code string
}

type ComponentCodeMap map[string][]ComponentCode

type ComponentExampleCodeMap map[string][]ComponentCode
