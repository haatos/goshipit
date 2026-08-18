package internal

import (
	"os"
	"strings"
)

var Settings *AppSettings

func NewSettings() *AppSettings {
	settings := AppSettings{
		Title:  "goship.it",
		Domain: getEnvOrDefault("DOMAIN", "localhost"),
		Port:   getEnvOrDefault("PORT", ":8080"),
	}
	if !strings.HasPrefix(settings.Port, ":") {
		settings.Port = ":" + settings.Port
	}
	return &settings
}

func getEnvOrDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

type AppSettings struct {
	Title  string
	Domain string
	Port   string
}
