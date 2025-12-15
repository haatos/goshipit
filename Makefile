.DEFAULT_GOAL := dev

GOOS := "linux"
GOARCH := "amd64"

deploy:
	npx @tailwindcss/cli -i input.css -o ./internal/assets/public/static/css/tw.css --minify
	go run cmd/generate/main.go
	templ generate
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w" -o bin/main cmd/server/main.go
	ssh $(user)@$(ip) "sudo service goshipit stop"
	scp 'bin/main' $(user)@$(ip):/opt/goshipit/
	ssh $(user)@$(ip) "sudo service goshipit start"

gen:
	go run cmd/generate/main.go

tw:
	@npx @tailwindcss/cli -i input.css -o ./internal/assets/public/static/css/tw.css --watch

dev: gen
	@templ generate -watch -proxyport=7332 -proxy="http://localhost:8080" -open-browser=false -cmd="go run cmd/server/main.go"
