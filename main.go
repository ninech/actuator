package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/ninech/actuator/actuator"
)

// DefaultPort defines the app's default HTTP port
const DefaultPort = "8080"

var srv *http.Server

func main() {
	fmt.Println("Startup sequence initiated ...")

	if err := actuator.LoadConfiguration(); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	serveRequests()
	waitForQuitSignal()
	log.Println("Device shutting down ...")
	gracefullyStop(5 * time.Second)
}

func serveRequests() {
	engine := actuator.NewWebhookEngine(actuator.DebugMode)
	router := engine.GetRouter()

	srv = &http.Server{
		Addr:    ":" + DefaultPort,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()
}

func waitForQuitSignal() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}

func gracefullyStop(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
