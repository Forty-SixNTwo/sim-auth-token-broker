package graceful

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var ready atomic.Value

func init() {
	ready.Store(true)
}

func StartServer(server *http.Server, timeout time.Duration, logger *slog.Logger, deferables ...func()) error {
	ready.Store(true)

	server.RegisterOnShutdown(func() {
		ready.Store(false)
	})

	server.BaseContext = func(l net.Listener) context.Context {
		return context.Background()
	}

	origHandler := server.Handler
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if ready.Load().(bool) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})

	mux.Handle("/", origHandler)
	server.Handler = mux

	errCh := make(chan error, 1)
	go func() {
		logger.Info("Server listening", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		logger.Info("Shutdown signal received", "signal", sig)
	case err := <-errCh:
		return err
	}

	for _, fn := range deferables {
		fn()
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	logger.Info("Server shutdown")
	return nil
}
