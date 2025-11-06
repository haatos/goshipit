# GoShip.it

Golang + Templ + HTMX (+ TailwindCSS + DaisyUI) component library to enhance developing an application using the GOTH stack.

The library contains DaisyUI components translated into Templ components that can be easily customized using both TailwindCSS and DaisyUI.

Updated to support DaisyUI 5. If you're looking for DaisyUI 4 compatible components, take a look [here](https://old.goship.it)

## Getting started

Install node dependencies:
`npm i -D`

Build TailwindCSS:
`make tw`

Generate component code/json, generate templates and run the server:
`make dev`

## Code/data generation

`cmd/generate/main.go` is used to generate JSON, markdown and Go code from source code. JSON is used to store and load component and example component source code to be displayed in HTML. The generator also generates a markdown file that contains up-to-date types for components (from `internal/model/components.go`), and .go file containing a mapping of example names to _templ_ components.

## Contributing

Components are placed in individual .templ files in `internal/views/components/`. The name of the file is used as the name of the component (converted from snake_case to Capitalized Component Name). The .templ file starts with a category name as a comment, e.g. `// data_display`.

For example

`internal/views/components/accordion.templ`

```go
// data_display
package components

type AccordionRowProps struct {
	Label string
	Type  string
	Name  string
}

templ AccordionRow(props AccordionRowProps) {
	<div class="collapse collapse-arrow bg-base-200 join-item">
		<input type="checkbox" name="templ-accordion"/>
		<div class="collapse-title text-xl font-medium">{ label }</div>
		<div class="collapse-content bg-base-300">
			{ children... }
		</div>
	</div>
}
```

Each component also has an examples file with a corresponding name in `internal/views/examples/`. The file can contain multiple examples, each starting with a comment `// example` with any lines below this line belonging to the example up to the next `// example` line or EOF. E.g.:

`internal/views/examples/textarea.templ`

```go
package examples

import "github.com/haatos/goshipit/internal/views/components"

// example
templ BasicTextarea() {
	<div class="pt-4">
		@components.Textarea(
			components.TextareaProps{
				Label: "Description",
				Name:  "description",
			},
		)
	</div>
}

// example
templ BasicTextareaWithError() {
	<div class="pt-4">
		@components.Textarea(
			components.TextareaProps{
				Label: "Description",
				Name:  "description",
				Err:   "Description cannot be empty",
			},
		)
	</div>
}
```

Some examples have corresponding handler functions to provide dummy data for the component to display its usage. E.g. `internal/handler/components.go` contains the handler for lazy-loading example:

```go
// LazyLoadExample
func GetLazyLoadExample(c echo.Context) error {
	time.Sleep(2 * time.Second)

	return render(c, http.StatusOK, examples.LazyLoadResult())
}

// LazyLoadExample
```

The handler must be enclosed with the name of the example component as comment on both sides of the handler's function(s).

## CLI

To install the goship.it CLI, run `go install github.com/haatos/goshipit/cmd/gsi`, or clone the repository and build it yourself: `go build -o gsi cmd/gsi/main.go`.

### Usage

Add components by running `gsi add button`. Components are added to _internal/views/components_ folder.

To see a list of all (and installed components) run `gsi list`.

To remove a component, run `gsi remove button`.
