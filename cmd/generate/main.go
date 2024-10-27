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

	"github.com/haatos/goshipit/internal/model"
)

const (
	componentCodeMapJSONPath        = "generated/component_code_map.json"
	componentExampleCodeMapJSONPath = "generated/component_example_code_map.json"
	componentModelsPath             = "internal/model/components.go"
	componentsDir                   = "internal/views/components/"
	componentsHandlerPath           = "internal/handler/components.go"
	examplesDir                     = "internal/views/examples/"
	generatedDir                    = "generated"
	generatedComponentsPath         = "generated/components.go"
	generatedTypesPath              = "generated/types.md"
)

func main() {
	generateComponentCodeMap()
	generateComponentExampleCodeMap()
	generateComponentMap()
	generateTypesMarkdown()
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
	var category string
	for scanner.Scan() {
		line := scanner.Text()
		if category == "" {
			category = strings.TrimPrefix(line, "// ")
			continue
		}
		if strings.HasPrefix(line, "templ ") {
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
		model.ComponentCode{Name: componentName, Code: function})
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

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "// example") {
					if inExample {
						m[componentName] = append(
							m[componentName],
							model.ComponentCode{Name: functionName, Code: functionLines})
						functionName = ""
						functionLines = []string{}
					}
					inExample = true
					continue
				}

				if strings.HasPrefix(line, "templ ") && functionName == "" {
					functionName = strings.TrimPrefix(line, "templ ")
					functionName = functionName[:strings.Index(functionName, "(")]
				}

				if inExample {
					functionLines = append(functionLines, line)
				}
			}

			f.Close()
			m[componentName] = append(
				m[componentName],
				model.ComponentCode{Name: functionName, Code: functionLines})
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
						m[comName][i].Handler = functionLines
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

func generateTypesMarkdown() {
	b, err := os.ReadFile(componentModelsPath)
	if err != nil {
		log.Fatal(err)
	}
	b = append([]byte("```go\n"), b...)
	b = append(b, '`', '`', '`', '\n')

	f, err := os.Create(generatedTypesPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Write(b)
}
