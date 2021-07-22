package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DevDisc/thotbot/src/lib/app"
)

func main() {
	application := app.Cli(&app.CliMethods{
		RunApp: runApp,
	})

	err := application.Run(os.Args)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func runApp(cfg *app.Config) error {

	// Create Server
	server, err := app.NewServer(cfg)
	if err != nil {
		return err
	}

	err = server.Start()
	if err != nil {
		return err
	}

	// Run until killed
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = server.Stop()
	if err != nil {
		return err
	}
	return nil
}
