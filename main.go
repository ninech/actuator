package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ninech/frau-schultz/api"
)

// DefaultPort defines the app's default HTTP port
const DefaultPort = "8080"

var srv *http.Server

func main() {
	fmt.Println("Hallo, ich bin Frau Schultz!")
	serveRequests()
	waitForQuitSignal()
	log.Println("Frau Schultz sagt auf Wiedersehen ...")
	gracefullyStop(5 * time.Second)
}

func serveRequests() {
	router := api.GetMainEngine()

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
	signal.Notify(quit, os.Interrupt)
	<-quit
}

func gracefullyStop(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
