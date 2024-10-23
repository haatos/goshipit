package internal

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/haatos/goshipit/internal/model"
)

var ComponentCodeMap model.ComponentCodeMap
var ComponentExampleCodeMap model.ComponentExampleCodeMap

func init() {
	getComponentCodeMap()
	getComponentExampleCodeMap()
}

func codeSliceToMarkdown(s []string) string {
	n := make([]string, len(s)+2)
	n[0] = "```go"
	n[len(s)+2-1] = "```"
	copy(n[1:len(s)+1], s)
	return strings.Join(n, "\n")
}

func getComponentCodeMap() {
	f, err := os.Open("generated/component_code_map.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(b, &ComponentCodeMap); err != nil {
		log.Fatal(err)
	}

	for k := range ComponentCodeMap {
		for i := range ComponentCodeMap[k] {
			ComponentCodeMap[k][i].Label = SnakeCaseToNormal(ComponentCodeMap[k][i].Name)
			ComponentCodeMap[k][i].Markdown = codeSliceToMarkdown(ComponentCodeMap[k][i].Code)
		}
	}
}

func getComponentExampleCodeMap() {
	f, err := os.Open("generated/component_example_code_map.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(b, &ComponentExampleCodeMap); err != nil {
		log.Fatal(err)
	}

	for k := range ComponentExampleCodeMap {
		for i := range ComponentExampleCodeMap[k] {
			ComponentExampleCodeMap[k][i].Label = SnakeCaseToNormal(ComponentExampleCodeMap[k][i].Name)
			ComponentExampleCodeMap[k][i].Markdown = codeSliceToMarkdown(ComponentExampleCodeMap[k][i].Code)
		}
	}
}

func SnakeCaseToNormal(s string) string {
	b := []byte(s)
	for i := range b {
		if i == 0 {
			b[i] = b[i] - ('a' - 'A')
		}
		if b[i] == '_' {
			b[i] = ' '
		}
	}
	return string(b)
}
