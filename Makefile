.DEFAULT_GOAL := dev

GOOS := "linux"
GOARCH := "amd64"

deploy:
	npx tailwindcss -o ./public/static/css/tw.css --minify
	go run cmd/generate/main.go
	templ generate
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w" -o bin/main cmd/server/main.go
	scp -r 'content' $(user)@$(ip):/opt/goshipit/
	scp -r 'generated' $(user)@$(ip):/opt/goshipit/
	scp -r 'public' $(user)@$(ip):/opt/goshipit/
	ssh $(user)@$(ip) "sudo service goshipit stop"
	scp 'bin/main' $(user)@$(ip):/opt/goshipit/
	ssh $(user)@$(ip) "sudo service goshipit start"

gen:
	go run cmd/generate/main.go

tw:
	@npx tailwindcss -i input.css -o public/static/css/tw.css --watch

dev: gen
	@templ generate -watch -proxyport=7332 -proxy="http://localhost:8080" -open-browser=false -cmd="go run cmd/server/main.go"
