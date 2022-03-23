package utils

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// ENVs
func GetEnvOrPanic(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Missing environment variable: %s", key))
	} else {
		return val
	}
}

// Rand numbers
func RandomInteger(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// Delay execution helpers
func RandomDelay(secs int) time.Duration {
	delay := RandomInteger(1, secs)

	logrus.WithField("delaySeconds", delay).Trace("Random delay")
	return time.Duration(delay) * time.Second
}

func ExecuteDelayed(waitTime time.Duration, f func()) {
	logrus.WithField("waitTime", waitTime).Debug("Scheduling delayed action")
	time.AfterFunc(waitTime, f)
}

func ExecuteDelayedRandom(max int, f func()) {
	waitTime := RandomDelay(max)
	logrus.WithField("waitTime", waitTime).Debug("Scheduling delayed action")
	time.AfterFunc(waitTime, f)
}

// Time formatting
func FormatSecondsToMinutes(secs int) string {
	var minutes = secs / 60
	var seconds = secs - minutes*60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)

}

func IndexOf[T comparable](value T, collection []T) int {
	for idx, element := range collection {
		if element == value {
			return idx
		}
	}
	return -1
}
