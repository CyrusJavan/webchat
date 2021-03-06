package main

import (
	"context"
	"fmt"
	"github.com/CyrusJavan/webchat/chatservice"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.WithError(err).Error("run returned with error")
	}
}

func run() error {
	natsUrl := "nats-service"
	if os.Getenv("DEVELOPMENT_MODE") == "TRUE" {
		natsUrl = nats.DefaultURL
	}

	nc, err := nats.Connect(natsUrl, nil)
	if err != nil {
		return fmt.Errorf("could not connect to nats server: %w", err)
	}

	chatServer := chatservice.NewServer(nc)

	srv := http.Server{
		Addr: ":8080",
		Handler: chatServer,
	}

	errChan := make(chan error, 1)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-interrupt:
		log.Println("Received interrupt, starting graceful shutdown...")
		ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctxShutDown); err != nil {
			return err
		}
		log.Println("Shutdown successful")
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("chatserivce:%w", err)
		}
	}

	return nil
}