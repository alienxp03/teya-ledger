package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alienxp03/teya-ledger/api"
)

func main() {
	startServer()
}

func startServer() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		<-ctx.Done()
		stop()
	}()

	addr := flag.String("addr", "localhost:8080", "HTTP network address")
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		logger.Error("Could not listen", "error", err)
		os.Exit(1)
	}

	api := &api.API{
		Logger: logger,
	}

	srv := &http.Server{
		Handler: api,
	}

	go func() {
		<-ctx.Done()
		logger.Info("Shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	// Start the server
	logger.Info("Ready to accept traffic", "address", *addr)
	if err := srv.Serve(lis); err != nil && err != http.ErrServerClosed {
		logger.Error("Could not start server", "error", err)
		os.Exit(1)
	}
}
