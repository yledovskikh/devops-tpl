package main

import (
	"context"

	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yledovskikh/devops-tpl/internal/config"
	"github.com/yledovskikh/devops-tpl/internal/db"
	"github.com/yledovskikh/devops-tpl/internal/dumper"
	"github.com/yledovskikh/devops-tpl/internal/handlers"
	"github.com/yledovskikh/devops-tpl/internal/inmemory"
	"github.com/yledovskikh/devops-tpl/internal/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	serverConfig := config.GetServerConfig()
	log.Logger = log.With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	r := chi.NewRouter()

	var s storage.Storage
	var err error
	if serverConfig.DatabaseDSN != "" {
		log.Info().Msg("Use Database Storage")
		s, err = db.New(serverConfig.DatabaseDSN, ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to use data source name")
		}
		//Закрываем коннекты в БД
		defer s.Close()
	} else {
		wg.Add(1)
		s = inmemory.NewMetricStore()
		dumper.Imp(s, serverConfig.StoreFile)
		go dumper.Exec(&wg, ctx, s, serverConfig)
	}
	h := handlers.New(s)
	h.Key = serverConfig.Key

	logger := httplog.NewLogger("server", httplog.Options{
		JSON: true,
	})

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(handlers.CompressResponse)
	r.Use(handlers.DecompressRequest)
	r.Get("/", h.AllMetrics)
	r.Get("/ping", h.Ping)
	r.Post("/update/", h.UpdateJSONMetric)
	r.Post("/updates/", h.UpdatesJSONMetrics)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateURLMetric)
	r.Get("/value/{metricType}/{metricName}", h.GetURLMetric)
	r.Post("/value/", h.GetJSONMetric)

	srv := &http.Server{
		Addr:    serverConfig.ServerAddress,
		Handler: r,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("")
		}
	}()
	log.Info().Msg("Server Started")

	<-done

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("")
	}

	log.Info().Msg("Server Stopped")
	cancel()
	wg.Wait()
	log.Info().Msg("Server Exited Properly")

}
