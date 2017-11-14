package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// Fetch initial state
	err := updateStatus()
	if err != nil {
		log.Fatalln("failed to fetch initial state")
	}

	// Connect to MQTT
	reqUpdate := make(chan int)
	err = runMqtt(reqUpdate)
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

	if setStatus(evaluateStatus(s)) {
		mqttPublish()
	}
	return nil
}

func evaluateStatus(s *imapStatus) *ArcStatus {
	if s.unreadCount == 0 {
		return &ArcStatus{
			Ring:        lsOn,
			CoreRed:     lsOff,
			CoreGreen:   lsOff,
			CoreBlue:    lsOn,
			UnreadCount: s.unreadCount,
			Message:     "No unread messages",
		}
	}

	if s.unreadCount == 1 {
		return &ArcStatus{
			Ring:        lsOn,
			CoreRed:     lsOff,
			CoreGreen:   lsOff,
			CoreBlue:    lsBlink,
			UnreadCount: s.unreadCount,
			Message:     "One unread message",
		}
	}

	if s.unreadCount < 3 {
		return &ArcStatus{
			Ring:        lsOn,
			CoreRed:     lsBlink,
			CoreGreen:   lsOff,
			CoreBlue:    lsOff,
			UnreadCount: s.unreadCount,
			Message:     fmt.Sprintf("%d unread messages", s.unreadCount),
		}
	}

	return &ArcStatus{
		Ring:        lsBlink,
		CoreRed:     lsOn,
		CoreGreen:   lsOff,
		CoreBlue:    lsOff,
		UnreadCount: s.unreadCount,
		Message:     fmt.Sprintf("%d unread messages", s.unreadCount),
	}
}
