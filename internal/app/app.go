package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/msjai/latency-ping/internal/config"
	"github.com/msjai/latency-ping/internal/controller"
	"github.com/msjai/latency-ping/internal/usecase"
	"github.com/msjai/latency-ping/internal/usecase/pingpong"
	repository "github.com/msjai/latency-ping/internal/usecase/repo"
)

// Run -.
func Run(cfg *config.Config) {
	l := cfg.L
	l.Infow("starting server...")

	repo, err := repository.New(cfg)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - repo.New: %w", err))
	}

	pingPong, err := pingpong.New(cfg)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - pingpong.New: %w", err))
	}

	// Use case
	latencyUseCase := usecase.New(
		repo,
		pingPong,
		cfg,
	)

	// initialize chi Mux object
	handler := chi.NewRouter()
	controller.NewRouter(handler, latencyUseCase, cfg)
	server := &http.Server{
		Addr:              cfg.RunAddress,
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}

	// Graceful server shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		l.Infof("Listening on port %v", server.Addr)
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCancel()

		l.Infow("Shutting down server...")
		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				l.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			l.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		l.Fatalf("listen: %s\n", err)
	}
	l.Infow("Server exiting")

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
