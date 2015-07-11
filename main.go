package main

import (
	"log"
	"os"

	"github.com/bronislav/journal-2-papertrail/journal"
	"github.com/bronislav/journal-2-papertrail/papertrail"
)

func main() {
	socket := os.Getenv("JOURNAL_SOCKET")
	if socket == "" {
		socket = journal.DefaultSocket
	}
	url := os.Getenv("PAPERTRAIL_URL")
	if url == "" {
		log.Fatal("non-empty papertrail url (PAPERTRAIL_URL) is required. See http://help.papertrailapp.com")
	}
	logs, err := journal.Follow(socket)
	if err != nil {
		log.Fatal(err.Error())
	}
	logger, err := papertrail.New(url)
	if err != nil {
		log.Fatal(err.Error())
	}
	for {
		select {
		case data := <-logs:
			if _, err := logger.Write(data); err != nil {
				log.Print(err.Error())
			}
		}
	}
}
