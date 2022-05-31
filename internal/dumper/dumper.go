package dumper

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/yledovskikh/devops-tpl/internal/config"
	"github.com/yledovskikh/devops-tpl/internal/handlers"
	"github.com/yledovskikh/devops-tpl/internal/serializer"
	"github.com/yledovskikh/devops-tpl/internal/storage"
)

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*producer, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}
func (p *producer) WriteMetric(metric *storage.Metric) error {
	return p.encoder.Encode(&metric)
}
func (p *producer) Close() error {
	return p.file.Close()
}

type consumer struct {
	file *os.File
	// заменяем reader на scanner
	scanner *bufio.Scanner
}

func NewConsumer(filename string) (*consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *consumer) ReadMetric() (storage.Metric, error) {
	data := c.scanner.Bytes()

	log.Println("Read string - ", string(data))

	metric := storage.Metric{}
	err := json.Unmarshal(data, &metric)
	if err != nil {
		return storage.Metric{}, err
	}

	return metric, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}

func Exp(s storage.Storage, fileName string) {

	log.Println("Info start export metrics to file")
	producer, err := NewProducer(fileName)
	if err != nil {
		log.Println("Error Exp - producer, err := NewProducer(fileName)", err.Error())
	}
	defer producer.Close()
	gauges := s.GetAllGauges()
	for mName, mValue := range gauges {
		metric := serializer.SerializeGauge(mName, mValue, "")
		if err := producer.WriteMetric(&metric); err != nil {
			log.Println("Error Exp - metric := serializer.DecodingGauge(mName, mValue)", err.Error())
		}
	}
	counters := s.GetAllCounters()
	for mName, mValue := range counters {
		metric := serializer.SerializeCounter(mName, mValue, "")
		if err := producer.WriteMetric(&metric); err != nil {
			log.Println("Error Exp - metric := serializer.DecodingCounter(mName, mValue)", err.Error())
		}
	}
}

func Imp(s storage.Storage, fileName string) {
	log.Println("Info start import DATA from file")
	consumer, err := NewConsumer(fileName)
	if err != nil {
		log.Println(err)
	}
	defer consumer.Close()

	for consumer.scanner.Scan() {
		metric, err := consumer.ReadMetric()
		if err != nil {
			log.Println(err)
		}
		log.Println(metric)
		err = handlers.SaveStoreMetric(metric, s)
		if err != nil {
			log.Println(err)
		}
	}

}

func Exec(wg *sync.WaitGroup, ctx context.Context, storage storage.Storage, serverConfig config.ServerConfig) {
	defer wg.Done()
	dumpInt := time.NewTicker(serverConfig.StoreInterval)
	log.Println("INFO dump file")
	for {
		select {
		case <-ctx.Done():
			log.Println("INFO dump file before exit")
			Exp(storage, serverConfig.StoreFile)
			return
		case <-dumpInt.C:
			//log.Println("INFO dump file")
			Exp(storage, serverConfig.StoreFile)
		}
	}

}
