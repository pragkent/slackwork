package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pragkent/slackwork/config"
	"github.com/pragkent/slackwork/server"
)

var (
	flagAddr            = flag.String("addr", ":5555", "listening address")
	flagConfig          = flag.String("c", "slackwork.yaml", "config file")
	flagVersion         = flag.Bool("version", false, "show version")
	flagShutdownTimeout = flag.Duration("shutdown-timeout", 15*time.Second, "timeout for gracefully shutdown")
)

func init() {
	flag.Parse()
}

func run() error {
	if *flagVersion {
		printVersion()
		return nil
	}

	c, err := config.Load(*flagConfig)
	if err != nil {
		log.Printf("Load config error: %v", err)
		return err
	}

	srv, err := server.New(*flagAddr, c)
	if err != nil {
		log.Printf("server.New error: %v", err)
		return err
	}

	go func() {
		sigs := make(chan os.Signal)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		<-sigs
		log.Printf("Shutting down server...")

		if err = srv.Shutdown(*flagShutdownTimeout); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}

		log.Println("Server shut down complete")
	}()

	log.Printf("Listening on %v", srv.Addr())
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
