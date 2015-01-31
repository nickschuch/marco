package logging

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

func Info(msg string) {
	currenttime := time.Now().Local()
	log.WithFields(log.Fields{
		"date": currenttime.Format("15:04:05 02-01-2006"),
	}).Info(msg)
}

func ContainerFound(domain string, url string) {
	currenttime := time.Now().Local()
	log.WithFields(log.Fields{
		"date":   currenttime.Format("15:04:05 02-01-2006"),
		"domain": domain,
		"url":    url,
	}).Info("Found a container for the domain.")
}

func ContainerNotFound(domain string) {
	currenttime := time.Now().Local()
	log.WithFields(log.Fields{
		"date":   currenttime.Format("15:04:05 02-01-2006"),
		"domain": domain,
	}).Info("Container not found.")
}
