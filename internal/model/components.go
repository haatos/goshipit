package model

import (
	"github.com/a-h/templ"
)

type CompanyInfo struct {
	Icon        templ.Component
	Name        string
	Description string
	Copyright   string
}

type Image struct {
	Source string
	Alt    string
}

type Script struct {
	Source string
	Defer  bool
}
