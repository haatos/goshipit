# Step 1: Build stage with Go and Node.js
FROM golang:1.26-alpine AS builder

# Install Node.js, npm, and certificates
RUN apk add --no-cache nodejs npm ca-certificates git

# Install templ CLI
RUN go install github.com/a-h/templ/cmd/templ@v0.3.1020

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

RUN go env

# 1. Generate templ Go files from .templ templates
RUN templ generate

# 2. Process and minify Tailwind CSS (Adjust input and output paths as needed)
RUN npm i
RUN npx @tailwindcss/cli -i input.css -o ./internal/assets/public/static/css/tw.css --minify

# 3. Build and execute static asset generator Go application
RUN go run ./cmd/generate/main.go

# 4. Build the main server binary with embedded assets
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server/main.go

# Step 2: Minimal runtime image
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/server .

EXPOSE 8080

ENTRYPOINT ["/app/server"]
