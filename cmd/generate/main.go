package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/haatos/goshipit/internal/markdown"
	"github.com/haatos/goshipit/internal/model"
)

const (
	componentCodeMapJSONPath        = "generated/component_code_map.json"
	componentExampleCodeMapJSONPath = "generated/component_example_code_map.json"
	componentsDir                   = "internal/views/components/"
	componentsHandlerPath           = "internal/handler/components.go"
	examplesDir                     = "internal/views/examples/"
	generatedDir                    = "generated"
	generatedComponentsPath         = "generated/components.go"
)

func main() {
	generateComponentCodeMap()
	generateComponentExampleCodeMap()
	generateComponentMap()
}

func generateComponentCodeMap() {
	m := model.ComponentCodeMap{}
	if err := filepath.Walk(componentsDir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".templ") {
			if err := getComponentCode(path, info, m); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	if err := os.RemoveAll(generatedDir); err != nil {
		log.Fatal(err)
	}
	if err := os.Mkdir(generatedDir, 0755); err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fg, err := os.Create(componentCodeMapJSONPath)
	if err != nil {
		log.Fatal(err)
	}
	defer fg.Close()

	fg.Write(b)
}

func getComponentCode(path string, info fs.FileInfo, fmap model.ComponentCodeMap) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	componentName := strings.TrimSuffix(info.Name(), ".templ")

	scanner := bufio.NewScanner(f)
	function := []string{}
	inFunction := false
	description := []string{}
	inDescription := false
	daisyUIURL := ""
	var category string
	for scanner.Scan() {
		line := scanner.Text()
		if category == "" {
			category = strings.TrimPrefix(line, "// ")
			continue
		}
		if strings.HasPrefix(line, "// https://daisyui.com") {
			daisyUIURL = strings.TrimPrefix(line, "// ")
		}
		if strings.HasPrefix(line, "/*") {
			inDescription = true
			continue
		}
		if strings.HasPrefix(line, "*/") {
			inDescription = false
			continue
		}
		if inDescription {
			description = append(description, line)
			continue
		}
		if strings.HasPrefix(line, "type ") || strings.HasPrefix(line, "templ ") {
			inFunction = true
		}
		if inFunction {
			function = append(function, line)
		}
	}
	if _, ok := fmap[category]; !ok {
		fmap[category] = []model.ComponentCode{}
	}
	fmap[category] = append(
		fmap[category],
		model.ComponentCode{
			Name:        componentName,
			Code:        markdown.CodeSliceToMarkdown(function),
			Description: strings.Join(description, "\n"),
			DaisyUIURL:  daisyUIURL,
		})
	return nil
}

func generateComponentExampleCodeMap() {
	m := model.ComponentExampleCodeMap{}
	if err := filepath.Walk(examplesDir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".templ") {
			f, err := os.Open(path)
			if err != nil {
				return err
			}

			functionLines := []string{}
			var functionName string
			componentName := strings.TrimSuffix(info.Name(), ".templ")
			m[componentName] = []model.ComponentCode{}
			inExample := false
			description := []string{}
			title := ""
			inDescription := false

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "// example") {
					if inExample {
						m[componentName] = append(
							m[componentName],
							model.ComponentCode{
								Name:        functionName,
								Code:        markdown.CodeSliceToMarkdown(functionLines),
								Title:       title,
								Description: strings.Join(description, "\n"),
							})
						functionName = ""
						functionLines = []string{}
						description = []string{}
					}
					inExample = true
					continue
				}

				if strings.HasPrefix(line, "// ") && !strings.HasPrefix(line, "// example") {
					title = strings.TrimPrefix(line, "// ")
					continue
				}

				if strings.HasPrefix(line, "/*") {
					inDescription = true
					continue
				}
				if strings.HasPrefix(line, "*/") {
					inDescription = false
					continue
				}
				if inDescription {
					description = append(description, line)
				}

				if strings.HasPrefix(line, "templ ") && functionName == "" {
					functionName = strings.TrimPrefix(line, "templ ")
					functionName = functionName[:strings.Index(functionName, "(")]
				}

				if inExample && !inDescription {
					functionLines = append(functionLines, line)
				}
			}

			f.Close()
			m[componentName] = append(
				m[componentName],
				model.ComponentCode{
					Name:        functionName,
					Code:        markdown.CodeSliceToMarkdown(functionLines),
					Title:       title,
					Description: strings.Join(description, "\n"),
				})
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	for comName := range m {
		for i := range m[comName] {
			f, err := os.Open(componentsHandlerPath)
			if err != nil {
				log.Fatal(err)
			}

			inExampleHandler := false
			functionLines := []string{}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, fmt.Sprintf("// %s", m[comName][i].Name)) {
					if inExampleHandler {
						inExampleHandler = false
						m[comName][i].Handler = markdown.CodeSliceToMarkdown(functionLines)
						break
					} else {
						inExampleHandler = true
					}
					continue
				}

				if inExampleHandler {
					functionLines = append(functionLines, line)
				}
			}

			f.Close()
		}
	}

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fg, err := os.Create(componentExampleCodeMapJSONPath)
	if err != nil {
		log.Fatal(err)
	}
	defer fg.Close()

	fg.Write(b)
}

func generateComponentMap() {
	functionNames := []string{}
	if err := filepath.Walk(examplesDir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".templ") {
			f, err := os.Open(path)
			if err != nil {
				return err
			}

			inExample := false

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "// example") {
					inExample = true
				}

				if inExample && strings.HasPrefix(line, "templ ") {
					functionName := strings.TrimPrefix(line, "templ ")
					functionName = functionName[:strings.Index(functionName, "(")]
					functionNames = append(functionNames, functionName)
					inExample = false
				}
			}
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	writeGeneratedFunctions(functionNames)
}

func writeGeneratedFunctions(functionNames []string) {
	// write functions into a buffer
	src := bytes.NewBuffer(nil)
	src.WriteString("package generated\n\n")
	src.WriteString("import (\n")
	src.WriteString("\t\"github.com/a-h/templ\"\n")
	src.WriteString("\t\"github.com/haatos/goshipit/internal/views/examples\"\n")
	src.WriteString(")\n\n")
	src.WriteString("var ExampleComponents = map[string]templ.Component{\n")
	for _, name := range functionNames {
		src.WriteString(fmt.Sprintf("\t\"%s\": examples.%s(),\n", name, name))
	}
	src.WriteString("}\n")

	// format the buffer's bytes using gofmt
	b, err := format.Source(src.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	// write the file
	fg, err := os.Create(generatedComponentsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer fg.Close()
	fg.Write(b)
}
