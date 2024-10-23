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

func main() {
	generateComponentCodeMap()
	generateComponentExampleCodeMap()
	generatePreviewMap()
}

func generateComponentCodeMap() {
	path := "internal/views/components/"
	m := model.ComponentCodeMap{}
	if err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".templ") {
			if err := getComponentCode(path, info, m); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	if err := os.RemoveAll("generated"); err != nil {
		log.Fatal(err)
	}
	if err := os.Mkdir("generated", 0755); err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fg, err := os.Create("generated/component_code_map.json")
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
	fmap[category] = append(fmap[category], model.ComponentCode{Name: componentName, Code: function})
	return nil
}

func generateComponentExampleCodeMap() {
	path := "internal/views/examples/"
	m := model.ComponentExampleCodeMap{}
	if err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
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
						m[componentName] = append(m[componentName], model.ComponentCode{Name: functionName, Code: functionLines})
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

			m[componentName] = append(m[componentName], model.ComponentCode{Name: functionName, Code: functionLines})
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fg, err := os.Create("generated/component_example_code_map.json")
	if err != nil {
		log.Fatal(err)
	}
	defer fg.Close()

	fg.Write(b)
}

func generatePreviewMap() {
	path := "internal/views/examples/"
	functionNames := []string{}
	if err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
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

	b, err := format.Source(src.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	fg, err := os.Create("generated/components.go")
	if err != nil {
		log.Fatal(err)
	}
	defer fg.Close()
	fg.Write(b)
}
