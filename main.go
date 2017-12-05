package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

var (
	reqUpdate = make(chan int)
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("env: %s\n", err)
	} else {
		log.Printf("env: loaded\n")
	}

	imapInit()
	mqttInit()
}

func main() {
	// Fetch initial state
	err := updateStatus()
	if err != nil {
		log.Fatalln("failed to fetch initial state")
	}

	// Connect to MQTT
	err = runMqtt()
	if err != nil {
		log.Fatalln("failed to connect to mqtt")
	}

	// Run HTTP server
	runHttp()

	// Start background updates
	go func() {
		for {
			time.Sleep(5 * time.Second)
			reqUpdate <- 0
		}
	}()
	go func() {
		for {
			<-reqUpdate
			updateStatus()
		}
	}()

	// Wait for exit
	ch := make(chan bool)
	<-ch
}

func updateStatus() error {
	s, err := fetchImapStatus()
	if err != nil {
		return err
	}

	status := evaluateStatus(s)
	if setStatus(status) {
		mqttPublish(status)
	}
	return nil
}

func evaluateStatus(s *imapStatus) *ArcStatus {
	switch s.unreadCount {
	case 0:
		// Ring: on
		// Core: blue
		return &ArcStatus{
			Ring:        lsOn,
			CoreRed:     lsOff,
			CoreGreen:   lsOff,
			CoreBlue:    lsOn,
			UnreadCount: s.unreadCount,
			Message:     "No unread messages",
		}

	case 1:
		// Ring: on
		// Core: blue (blink)
		return &ArcStatus{
			Ring:        lsOn,
			CoreRed:     lsOff,
			CoreGreen:   lsOff,
			CoreBlue:    lsBlink,
			UnreadCount: s.unreadCount,
			Message:     "1 unread message",
		}

	case 2:
		// Ring: on
		// Core: yellow (blink)
		return &ArcStatus{
			Ring:        lsOn,
			CoreRed:     lsBlink,
			CoreGreen:   lsBlink,
			CoreBlue:    lsOff,
			UnreadCount: s.unreadCount,
			Message:     "2 unread messages",
		}

	case 3:
		// Ring: on
		// Core: red (blink)
		return &ArcStatus{
			Ring:        lsOn,
			CoreRed:     lsBlink,
			CoreGreen:   lsOff,
			CoreBlue:    lsOff,
			UnreadCount: s.unreadCount,
			Message:     "3 unread messages",
		}

	default:
		// Ring: on
		// Core: yellow-red
		return &ArcStatus{
			Ring:        lsBlink,
			CoreRed:     lsBlink,
			CoreGreen:   lsOff,
			CoreBlue:    lsBlink2,
			UnreadCount: s.unreadCount,
			Message:     fmt.Sprintf("%d unread messages", s.unreadCount),
		}
	}
}
