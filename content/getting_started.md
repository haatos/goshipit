# Getting started

To get started with goship.it in a new project using _Echo_ router:

- Create a new folder for your project
- Initialize the module by running
  - `go mod init github.com/my/package`
  - replace the package with your own repository!
- Get Go modules:
  - `go get -u github.com/labstack/echo/v4`
  - `go get -u github.com/a-h/templ`
- Install Templ CLI:
  - `go install github.com/a-h/templ/cmd/templ@latest`
- Install TailwindCSS and DaisyUI:
  - `npm i -D tailwindcss @tailwindcss/typography daisyui`
- Initialize TailwindCSS:
  - `npx tailwindcss init`
- Configure `tailwind.config.js`, which the previous command generated, to look like this:

```javascript
module.exports = {
  content: ["internal/views/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    themes: ["light"],
  },
};
```

- Create `input.css` at the base of your project with the following contents:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

- Create `Makefile` at the base of your project with the following contents:

```make
tw:
	@npx tailwindcss -i input.css -o public/static/css/tw.css --watch

dev:
	@templ generate -watch -proxy="http://localhost:8080" -open-browser=false -cmd="go run main.go"
```

- Place the following rows in `main.go` (remember to update the components package import path to match your project):

```go
package main

import (
	"embed"
	"log"
	"net/http"

	"<your project>/internal/views/components"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Static("/", "public")

	e.GET("/", func(c echo.Context) error {
		accordion := components.AccordionExample()
        return render(accordion)
	})

	e.Start(":8080")
}

func render(c echo.Context, component templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := component.Render(c.Request().Context(), buf); err != nil {
		    return err
	}
	return c.HTML(http.StatusOK, buf.String())
}

```

- Create the accordion component `internal/views/components/accordion.templ`:

```go
package components

type AccordionRowProps struct {
	Label string
	Type  string
	Name  string
}

templ AccordionRow(props AccordionRowProps) {
	<div class="collapse collapse-arrow bg-base-300 join-item">
		<input
			if props.Type == "" {
				type="checkbox"
			} else {
				type={ props.Type }
			}
			name={ props.Name }
		/>
		<div class="collapse-title text-xl font-medium">{ props.Label }</div>
		<div class="collapse-content bg-base-200">
			{ children... }
		</div>
	</div>
}

templ AccordionExample() {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<link rel="stylesheet" type="text/css" href="/static/css/tw.css"/>
			<title>Document</title>
		</head>
		<body class="w-full h-full min-h-svh">
			<main>
				<div>
					@AccordionRow(AccordionRowProps{Label: "Accordion row 1", Type: "checkbox"}) {
						<p>This is the first content</p>
					}
					@AccordionRow(AccordionRowProps{Label: "Accordion row 2", Type: "checkbox"}) {
						<p>This is the second content</p>
					}
				</div>
			</main>
		</body>
	</html>
}
```

At this point, the filetree of your project should look something like this:

```sh
.
├── Makefile
├── go.mod
├── go.sum
├── input.css
├── internal
│   └── views
├── main.go
├── node_modules
├── package-lock.json
├── package.json
├── public
│   └── static
│       └── css
│           └── tw.css
└── tailwind.config.js
```

If you are using VSCode as your IDE, you should also add a `.vscode/settings.json` with the following contents (or place these settings in some other VSCode configuration file):

```json
{
  "[templ]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "a-h.templ"
  },
  "tailwindCSS.includeLanguages": {
    "templ": "html"
  },
  "emmet.includeLanguages": {
    "templ": "html"
  }
}
```

These will enable TailwindCSS autocompletions and HTML element autocompletions (emmet), as well as automatically formatting `.templ` files when saving.

Finally, you can run the example application by running `make tw` and `make dev` in two separate terminals. The site with the accordion should now be visible at http://localhost:8080.
