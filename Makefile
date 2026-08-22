.DEFAULT_GOAL := dev

region ?= europe-west4
port ?= 8080

deploy:
    docker build --platform=linux/amd64 -t $(region)-docker.pkg.dev/goshipit/goshipit-repo/goship-it-app:latest .
    docker push $(region)-docker.pkg.dev/goshipit/goshipit-repo/goship-it-app:latest
    gcloud run deploy goshipit-service \
      --image=$(region)-docker.pkg.dev/goshipit/goshipit-repo/goship-it-app:latest \
      --region=$(region) \
      --platform=managed \
      --port=$(port) \
      --allow-unauthenticated

gen:
	go run cmd/generate/main.go

tw:
	@npx @tailwindcss/cli -i input.css -o ./internal/assets/public/static/css/tw.css --watch

dev: gen
	@templ generate -watch -proxyport=7332 -proxy="http://localhost:8080" -open-browser=false -cmd="go run cmd/server/main.go"
