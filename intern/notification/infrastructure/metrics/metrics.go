package metrics

import (
	"context"
	"log/slog"
	"net/http"

	"final-project/intern/notification/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func AddMetricsHandler(cfg *config.Config) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("metrics server stopped", "error", err.Error())
		}
	}()

	return server
}

func Shutdown(ctx context.Context, server *http.Server) {
	if server == nil {
		return
	}

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown metrics server", "error", err.Error())
	}
}
