package dumper

import (
	"context"
	"encoding/json"
	"github.com/yledovskikh/devops-tpl/internal/config"
	"github.com/yledovskikh/devops-tpl/internal/serializer"
	"github.com/yledovskikh/devops-tpl/internal/storage"
	"io/ioutil"
	"log"
	"time"
)

func Exp(s storage.Storage, fileName string) {

	metrics := s.GetMetricsSerialize()
	file, _ := json.MarshalIndent(metrics, "", " ")

	_ = ioutil.WriteFile(fileName, file, 0644)
	log.Printf("Exp metrics to file - %s", fileName)
}

func imp(s storage.Storage, fileName string) {
	file, _ := ioutil.ReadFile(fileName)
	metrics := serializer.Metrics{}
	_ = json.Unmarshal(file, &metrics)
	s.SetMetricsSerialize(metrics)
}

func Exec(ctx context.Context, storage storage.Storage, serverConfig config.ServerConfig) {
	if serverConfig.Restore {
		imp(storage, serverConfig.StoreFile)
	}

	serverConfig.StoreInterval = 5 * time.Second
	dumpInt := time.NewTicker(serverConfig.StoreInterval)
	for {
		select {
		case <-ctx.Done():
			log.Println("INFO dump file before exit")
			Exp(storage, serverConfig.StoreFile)
			return
		case <-dumpInt.C:
			log.Println("INFO dump file")
			Exp(storage, serverConfig.StoreFile)
		}
	}

}
