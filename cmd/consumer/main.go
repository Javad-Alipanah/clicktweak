package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"clicktweak/internal/app/consumer"
	"clicktweak/internal/pkg/config"
	"clicktweak/internal/pkg/db/impl"
	"clicktweak/internal/pkg/model"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func main() {
	// set logger format
	log.StandardLogger().SetFormatter(&prefixed.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	})

	// check config path received
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <config-path>", os.Args[0])
	}

	if err := config.Init(os.Args[1], func() {
		log.Warningln("restart required due to config change")
	}); err != nil {
		log.Fatalln(err)
	}

	// Unmarshal config
	if err := config.Unmarshal(&consumer.Config); err != nil {
		log.Fatalln(err)
	}

	lvl, err := log.ParseLevel(consumer.Config.Log.Level)
	if err != nil {
		log.Warningln(err)
		log.Warningln("logging level unchanged")
	}

	log.SetLevel(lvl)

	// Initialize log database
	db, err := sqlx.Connect(consumer.Config.LogDB.Dialect, consumer.Config.LogDB.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	logDB, err := impl.NewLogDB(db)
	if err != nil {
		log.Fatal(err)
	}

	// make a new reader that consumes from topic Topic, partition 0, at offset 0
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{consumer.Config.Kafka.Host},
		Topic:     consumer.Config.Kafka.Topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(0)

	for {
		var elem = make([]*model.Log, consumer.Config.APP.Batch)
		begin := time.Now()
		var i int
		for i = 0; i < consumer.Config.APP.Batch && time.Now().Sub(begin) < consumer.Config.APP.Timeout; i++ {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Fatalln(err)
				break
			}
			elem[i] = new(model.Log)
			_ = json.Unmarshal(m.Value, elem[i])
		}

		if err = logDB.Save(elem, i); err != nil {
			log.Error(err)
		}
	}

	r.Close()
}
