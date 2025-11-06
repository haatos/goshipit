package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/haatos/goshipit/internal/views"
)

var componentsDirParts = []string{"internal", "views", "components"}

func main() {
	code := run(os.Stdout, os.Stderr, os.Args)
	if code != 0 {
		os.Exit(code)
	}
}

const usageText string = `usage: gsi <command> [<args>...]

gsi - GoShip.it CLI

commands:
  add <name>      Add a component to internal/views/components
  remove <name>   Remove a component from internal/views/components
  list            List all available components

`

func run(stdout, stderr io.Writer, args []string) int {
	if len(args) < 2 {
		_, _ = fmt.Fprint(stdout, usageText)
		return 64
	}

	switch args[1] {
	case "add":
		if len(args) < 3 {
			_, _ = fmt.Fprint(stdout, usageText)
			return 64
		}
		name := strings.Join(args[2:], " ")
		return runAddCommand(stderr, name)
	case "remove":
		if len(args) < 3 {
			_, _ = fmt.Fprint(stdout, usageText)
			return 64
		}
		name := strings.Join(args[2:], " ")
		return runRemoveCommand(stdout, name)
	case "list":
		return runListCommand(stdout)
	}

	return 1
}

func runAddCommand(stderr io.Writer, name string) int {
	name = strings.ReplaceAll(name, " ", "_")
	b, err := views.ComponentFS.ReadFile("components/" + name + ".templ")
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "component '%s' does not exist\n", name)
		return 1
	}
	packageIndex := bytes.Index(b, []byte("package"))
	b = b[packageIndex:]

	dirPath := filepath.Join(componentsDirParts...)
	fileName := fmt.Sprintf("%s.templ", name)
	fullPath := filepath.Join(dirPath, fileName)

	componentExists := true
	if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
		componentExists = false
	}

	if componentExists {
		var ans string
		fmt.Printf("Component '%s' already exists. Replace [Y/n]: ", name)
		if _, err := fmt.Scan(&ans); err != nil {
			log.Fatal(err)
		}
		if ans == "n" {
			return 0
		}
		fmt.Printf("Replacing component '%s'...\n", name)
	}

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Fatal("err making directories '", dirPath+"':", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		log.Fatal("err creating file '", fullPath+"':", err)
	}
	if _, err := f.Write(b); err != nil {
		log.Fatal("err writing component to '", fullPath+"':", err)
	}

	return 0
}

func runRemoveCommand(stdout io.Writer, name string) int {
	name = strings.ReplaceAll(name, " ", "_")
	dirPath := filepath.Join(componentsDirParts...)
	fileName := fmt.Sprintf("%s.templ", name)
	templFileName := fmt.Sprintf("%s_templ.go", name)
	if err := os.Remove(filepath.Join(dirPath, fileName)); err == nil {
		_, _ = fmt.Fprintf(stdout, "Removed '%s'\n", filepath.Join(dirPath, fileName))
	}
	if err := os.Remove(filepath.Join(dirPath, templFileName)); err != nil {
		_, _ = fmt.Fprintf(stdout, "Removed '%s'\n", filepath.Join(dirPath, templFileName))
	}
	return 0
}

func runListCommand(stdout io.Writer) int {
	entries, err := views.ComponentFS.ReadDir("components")
	if err != nil {
		log.Fatal(err)
	}
	_, _ = fmt.Fprint(stdout, "Run 'gsi add <name>' to add a component\n\n")
	m := make(map[string][]string)
	for _, entry := range entries {
		b, err := views.ComponentFS.ReadFile(filepath.Join("components", entry.Name()))
		if err != nil {
			log.Fatal(err)
		}
		split := bytes.Split(b, []byte("\n"))
		category := split[0]
		category = bytes.TrimPrefix(category, []byte("// "))
		category = bytes.ReplaceAll(category, []byte("_"), []byte(" "))
		categoryStr := string(category)
		m[categoryStr] = append(m[categoryStr], strings.TrimSuffix(entry.Name(), ".templ"))
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	for _, key := range keys {
		_, _ = fmt.Fprint(stdout, strings.ToUpper(key)+"\n")
		_, _ = fmt.Fprint(stdout, "=======================\n")
		for _, comp := range m[key] {
			p := filepath.Join("internal", "views", "components", comp+".templ")
			_, err := os.Stat(p)
			exists := err == nil
			line := strings.ReplaceAll(comp, "_", " ")
			if exists {
				line += " ✓"
			}
			line += "\n"

			_, _ = fmt.Fprint(stdout, line)
		}
		_, _ = fmt.Fprint(stdout, "\n")
	}

	return 0
}
