package main

import (
	"fmt"
	"sync"
)

type ledState string

const (
	lsInitial ledState = ""
	lsOff     ledState = "off"
	lsOn      ledState = "on"
	lsBlink   ledState = "blink"
	lsBlink2  ledState = "blink_2"
)

type ArcStatus struct {
	Ring        ledState `json:"ring"`
	CoreRed     ledState `json:"core_r"`
	CoreGreen   ledState `json:"core_g"`
	CoreBlue    ledState `json:"core_b"`
	UnreadCount uint32   `json:"unread"`
	Message     string   `json:"msg"`
}

var (
	currentStatus     *ArcStatus
	storedStatus      *ArcStatus
	currentStatusLock sync.Mutex
)

func init() {
	currentStatus = &ArcStatus{
		Ring:      lsInitial,
		CoreRed:   lsInitial,
		CoreGreen: lsInitial,
		CoreBlue:  lsInitial,
	}
}

func setStatus(status *ArcStatus) bool {
	currentStatusLock.Lock()
	defer currentStatusLock.Unlock()

	if currentStatus.Ring != status.Ring ||
		currentStatus.CoreRed != status.CoreRed ||
		currentStatus.CoreGreen != status.CoreGreen ||
		currentStatus.CoreBlue != status.CoreBlue {
		if storedStatus == nil {
			fmt.Printf("update: [%s,%s,%s,%s] -> [%s,%s,%s,%s] (ring,r,b,g)\n",
				currentStatus.Ring, currentStatus.CoreRed, currentStatus.CoreGreen, currentStatus.CoreBlue,
				status.Ring, status.CoreRed, status.CoreGreen, status.CoreBlue)
			currentStatus = status
			return true
		}
		storedStatus = status
		return false
	}

	return false
}

func getStatus() *ArcStatus {
	currentStatusLock.Lock()
	defer currentStatusLock.Unlock()

	return currentStatus
}

func setTestStatus(status *ArcStatus) {
	currentStatusLock.Lock()
	defer currentStatusLock.Unlock()

	storedStatus = currentStatus
	currentStatus = status
}

func resetStatus() {
	currentStatusLock.Lock()
	defer currentStatusLock.Unlock()

	currentStatus = storedStatus
	storedStatus = nil
}
