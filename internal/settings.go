package internal

import (
	"bufio"
	"log"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

var Settings *AppSettings

func NewSettings() *AppSettings {
	settings := AppSettings{
		Title:        "goship.it",
		ContactEmail: os.Getenv("CONTACT_EMAIL"),
		Domain:       getEnvOrDefault("DOMAIN", "localhost"),
		Port:         getEnvOrDefault("PORT", ":8080"),
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
	Title        string
	ContactEmail string
	Domain       string
	Port         string
}

func ReadDotenv() {
	path := "./.env"
	re := regexp.MustCompile(`^[^0-9][A-Z0-9_]+=.+$`)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("err opening dotenv: ", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) > 0 && line[0] != '#' && re.Match(line) {
			split := strings.Split(string(line), "=")
			name := strings.TrimSpace(split[0])
			value := strings.TrimSpace(split[1])
			value = strings.Trim(value, `"`)
			if err := os.Setenv(name, value); err != nil {
				slog.Error("unable to setenv", slog.String(name, value))
			}
		} else {
			slog.Debug("not including invalid or empty line", "line", string(line))
		}
	}
}
