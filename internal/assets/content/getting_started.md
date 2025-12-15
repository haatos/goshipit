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
  - `npm i -D tailwindcss @tailwindcss/cli @tailwindcss/typography daisyui@latest`
- Create `input.css` at the root of the project with the following contents:

```input.css
@import "tailwindcss" source(none);
@plugin "@tailwindcss/typography";
@source "./internal/views/**/*.templ";
@plugin "daisyui" {
    themes:
        light --default,
        dark --prefersdark;
}
@plugin "daisyui/theme" {
    name: "light";
    default: true;
    prefersdark: false;
    color-scheme: "light";
    --color-base-100: oklch(92% 0 0);
    --color-base-200: oklch(87% 0 0);
    --color-base-300: oklch(70% 0 0);
    --color-base-content: oklch(0% 0 0);
    --color-primary: oklch(55% 0.135 66.442);
    --color-primary-content: oklch(98% 0.018 155.826);
    --color-secondary: oklch(53% 0.157 131.589);
    --color-secondary-content: oklch(98% 0.031 120.757);
    --color-accent: oklch(64% 0.222 41.116);
    --color-accent-content: oklch(98% 0.016 73.684);
    --color-neutral: oklch(43% 0 0);
    --color-neutral-content: oklch(98% 0 0);
    --color-info: oklch(52% 0.105 223.128);
    --color-info-content: oklch(98% 0.019 200.873);
    --color-success: oklch(50% 0.118 165.612);
    --color-success-content: oklch(98% 0.018 155.826);
    --color-warning: oklch(55% 0.195 38.402);
    --color-warning-content: oklch(98% 0.016 73.684);
    --color-error: oklch(52% 0.223 3.958);
    --color-error-content: oklch(97% 0.014 343.198);
    --radius-selector: 0.5rem;
    --radius-field: 0.5rem;
    --radius-box: 0.5rem;
    --size-selector: 0.28125rem;
    --size-field: 0.28125rem;
    --border: 1px;
    --depth: 1;
    --noise: 0;
}
@plugin "daisyui/theme" {
    name: "dark";
    default: false;
    prefersdark: true;
    color-scheme: "dark";
    --color-base-100: oklch(26% 0 0);
    --color-base-200: oklch(20% 0 0);
    --color-base-300: oklch(14% 0 0);
    --color-base-content: oklch(97% 0 0);
    --color-primary: oklch(79% 0.184 86.047);
    --color-primary-content: oklch(28% 0.066 53.813);
    --color-secondary: oklch(64% 0.2 131.684);
    --color-secondary-content: oklch(98% 0.031 120.757);
    --color-accent: oklch(64% 0.222 41.116);
    --color-accent-content: oklch(98% 0.016 73.684);
    --color-neutral: oklch(14% 0 0);
    --color-neutral-content: oklch(98% 0 0);
    --color-info: oklch(71% 0.143 215.221);
    --color-info-content: oklch(98% 0.019 200.873);
    --color-success: oklch(72% 0.219 149.579);
    --color-success-content: oklch(98% 0.018 155.826);
    --color-warning: oklch(70% 0.213 47.604);
    --color-warning-content: oklch(98% 0.016 73.684);
    --color-error: oklch(65% 0.241 354.308);
    --color-error-content: oklch(97% 0.014 343.198);
    --radius-selector: 0.5rem;
    --radius-field: 0.5rem;
    --radius-box: 0.5rem;
    --size-selector: 0.28125rem;
    --size-field: 0.28125rem;
    --border: 1px;
    --depth: 1;
    --noise: 0;
}
```

- Create `Makefile` at the base of your project with the following contents:

```make
tw:
	@npx @tailwindcss/cli -i input.css -o ./public/static/css/tw.css --watch

dev:
	@templ generate -watch -proxyport=7332 -proxy="http://localhost:8080" -open-browser=false -cmd="go run main.go"
```

- Place the following rows in `main.go` (remember to update the components package import path to match your project):

```go
package main

import (
	"net/http"

	"github.com/my/package/internal/views/components"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Static("/", "public")

	e.GET("/", func(c echo.Context) error {
		accordion := components.AccordionExample()
		return render(c, accordion)
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
│       └── components
│           └── accordion.templ
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
