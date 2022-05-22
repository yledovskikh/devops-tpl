package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yledovskikh/devops-tpl/internal/config"
	"github.com/yledovskikh/devops-tpl/internal/dumper"
	"github.com/yledovskikh/devops-tpl/internal/handlers"
	"github.com/yledovskikh/devops-tpl/internal/storage"
)

func main() {
	//ctx, cancel := context.WithCancel(context.Background())
	r := chi.NewRouter()
	s := storage.NewMetricStore()
	h := handlers.New(s)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(handlers.CompressResponse)
	//r.Use(middleware.Compress(5))
	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/html")
		rw.Write([]byte("Metrics Collection Server"))
	})
	r.Post("/update/", h.UpdateJSONMetric)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateURLMetric)
	r.Get("/value/{metricType}/{metricName}", h.GetURLMetric)
	r.Post("/value/", h.GetJSONMetric)

	ctx, cancel := context.WithCancel(context.Background())

	serverConfig := config.GetServerConfig()
	if serverConfig.Restore {
		dumper.Imp(s, serverConfig.StoreFile)
	}
	var wg sync.WaitGroup

	wg.Add(1)
	go dumper.Exec(&wg, ctx, s, serverConfig)

	srv := &http.Server{
		Addr:    serverConfig.ServerAddress,
		Handler: r,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	//defer func() {
	//	//dumper.Exp(s, serverConfig.StoreFile) // extra handling here
	//
	//}()
	//
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	cancel()
	wg.Wait()
	log.Print("Server Exited Properly")

}
