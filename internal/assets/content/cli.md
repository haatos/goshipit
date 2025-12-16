# CLI

## Installation

```
go install github.com/haatos/cmd/gsi
```

Then run `gsi` in your terminal to see a list of commands:

```
usage: gsi <command> [<args>...]

gsi - GoShip.it CLI

commands:
  add <name>      Add a component to internal/views/components
  remove <name>   Remove a component from internal/views/components
  list            List all available components
```

Use the CLI at the root of your application. New components will be added to `internal/views/components/`.
