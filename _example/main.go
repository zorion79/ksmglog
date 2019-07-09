package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/go-pkgz/lgr"
	"github.com/zorion79/ksmglog"
)

func main() {
	log.Setup(log.Debug, log.Msec, log.LevelBraces, log.CallerFile, log.CallerFunc) // setup default logger with go-pkgz/lgr

	ksmgUrl := []string{os.Getenv("EXMPL_KSMG_URL")}
	if strings.Contains(ksmgUrl[0], ",") {
		ksmgUrl = strings.Split(ksmgUrl[0], ",")
	}

	options := ksmglog.Opts{
		URL:      ksmgUrl,
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

	ctx, cancel := context.WithCancel(context.Background())
	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Printf("[WARN] interrupt signal")
		cancel()
	}()
	go service.Run(ctx)

	newLogChan := service.Channel()
	count := 0
	for {
		select {
		case record, ok := <-newLogChan:
			if !ok {
				log.Printf("[WARN] channel closed")
				return
			}
			count++
			if record.HashString == "" {
				log.Printf("%d new %+v", count, record)
			}
			log.Printf("%d\t%+v", count, record)
		}
	}
}
