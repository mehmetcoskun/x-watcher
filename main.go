package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

type WhatsAppMessage struct {
	ChatID  string `json:"chatId"`
	Text    string `json:"text"`
	Session string `json:"session"`
}

func main() {
	godotenv.Load()

	scraper := twitterscraper.New()
	scraper.LoginOpenAccount()

	previousFollowersCount := 0

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		profile, err := scraper.GetProfile(os.Getenv("TWITTER_USERNAME"))
		if err != nil {
			log.Println("Failed to get profile: ", err)
			continue
		}

		followersCount := profile.FollowersCount
		if followersCount > previousFollowersCount {
			message := WhatsAppMessage{
				ChatID:  os.Getenv("WHATSAPP_CHAT_ID"),
				Text:    fmt.Sprintf(os.Getenv("FOLLOWER_INCREASED_TEXT"), followersCount),
				Session: os.Getenv("WHATSAPP_SESSION"),
			}
			sendMessageToAPI(message)
		} else if followersCount < previousFollowersCount {
			message := WhatsAppMessage{
				ChatID:  os.Getenv("WHATSAPP_CHAT_ID"),
				Text:    fmt.Sprintf(os.Getenv("FOLLOWER_DECREASED_TEXT"), followersCount),
				Session: os.Getenv("WHATSAPP_SESSION"),
			}
			sendMessageToAPI(message)
		}

		previousFollowersCount = followersCount
	}
}

func sendMessageToAPI(message WhatsAppMessage) {
	url := os.Getenv("WHATSAPP_API_URL")
	requestBody, err := json.Marshal(message)
	if err != nil {
		log.Println("Failed to marshal message: ", err)
		return
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Failed to send message to API: ", err)
		return
	}
}
