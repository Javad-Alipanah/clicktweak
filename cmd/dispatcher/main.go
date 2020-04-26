package main

import (
	"clicktweak/internal/app/dispatcher"
	"clicktweak/internal/pkg/config"
	"clicktweak/internal/pkg/db/impl"
	"clicktweak/internal/pkg/model"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// TODO: change main file
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
	if err := config.Unmarshal(&dispatcher.Config); err != nil {
		log.Fatalln(err)
	}

	lvl, err := log.ParseLevel(dispatcher.Config.Log.Level)
	if err != nil {
		log.Warningln(err)
		log.Warningln("logging level unchanged")
	}

	log.SetLevel(lvl)

	// connect to database
	conn, err := gorm.Open(dispatcher.Config.Database.Dialect, dispatcher.Config.Database.ConnectionString)
	if err != nil {
		log.Fatalln(err)
	}

	var db = conn
	if lvl.String() == "debug" || lvl.String() == "trace" {
		db = conn.Debug()
	}

	urlDB, err := impl.NewUrlDB(db)
	if err != nil {
		log.Fatalln(err)
	}

	// run workers
	logs := make(chan model.Log, dispatcher.Config.Forwarder.ChannelSize)
	workers := dispatcher.NewWorkers(dispatcher.Config.Forwarder.Url, logs)
	workers.Run(dispatcher.Config.Forwarder.Workers)
	defer workers.Close()

	e := echo.New()
	e.GET("/:id", dispatcher.Redirect(urlDB, logs))

	e.Logger.Fatal(e.Start(dispatcher.Config.App.Listen))
}
