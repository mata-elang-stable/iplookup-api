package main

import (
	"context"
	"github.com/fadhilyori/iplookup-go/internal/app"
	"github.com/fadhilyori/iplookup-go/internal/config"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func runApp(cmd *cobra.Command, args []string) {
	conf := config.GetConfig()
	conf.SetupLogging()

	mainContext, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	mainApp := app.NewApp(conf.ListenAddress, conf.ListenPort)
	if conf.MMDBRegionFilePath != "" {
		mainApp.LoadRegionMaxmindDB(conf.MMDBRegionFilePath)
	}
	if conf.MMDBASNFilePath != "" {
		mainApp.LoadASNMaxmindDB(conf.MMDBASNFilePath)
	}

	if conf.RedisURL != "" && conf.EnableCache {
		// convert to slice to support multiple redis from comma
		urls := strings.Split(conf.RedisURL, ",")

		mainApp.EnableCache(urls, time.Duration(conf.CacheTTLSec)*time.Second)
	}

	if err := mainApp.Run(mainContext); err != nil {
		log.Fatalf("Error starting app: %s", err)
	}
}
