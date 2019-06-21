package main

import (
	"os"

	log "github.com/go-pkgz/lgr"
	"github.com/zorion79/ksmglog"
)

func main() {
	log.Setup(log.Debug, log.Msec, log.LevelBraces, log.CallerFile, log.CallerFunc) // setup default logger with go-pkgz/lgr

	options := ksmglog.Opts{
		URL:      os.Getenv("EXMPL_KSMG_URL"),
		User:     os.Getenv("EXMPL_KSMG_USER"),
		Password: os.Getenv("EXMPL_KSMG_PASS"),
	}

	service := ksmglog.NewService(options)
	records, err := service.GetLogs()
	if err != nil {
		log.Fatalf("could not get records: %v", err)
	}

	log.Printf("count of record=%d", len(records))
	log.Printf("first record\n%+v", records[0])
}
