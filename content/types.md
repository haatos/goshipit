# Types

## Introduction

Most components use a struct as input arguments. This provides a convenient way to pass default values to the components; struct fields initialize to the respective type's default value. We can use this feature to make some of the fields of a struct optional. For example, we can pass the `components.AnchorProps` argument to an anchor element, `<a>`:

```go
package components

type AnchorProps struct {
	Href      string
	Label     string
	LeftIcon  templ.Component
	RightIcon templ.Component
	Attrs     templ.Attributes
	Class     string
}

templ Anchor(props AnchorProps) {
	<a
		if props.Href != "" {
			href={ templ.SafeURL(props.Href) }
		}
		class={ "group flex items-center cursor-pointer", props.Class }
		{ props.Attrs... }
	>
		if props.LeftIcon != nil {
			<div class="inline-block mr-1">
				@props.LeftIcon
			</div>
		}
		{ props.Label }
		if props.RightIcon != nil {
			<div class="inline-block ml-1">
				@props.RightIcon
			</div>
		}
	</a>
}
```

Here we can define the anchor element to optionally have an `href` attribute, or if we choose, a `hx-get` instead:

```go
@components.Anchor(components.AnchorProps{Href: "/"})
@components.Anchor(components.AnchorProps{Attrs: templ.Attributes{"hx-get": "/"}})
```
