package main

import (
	"clicktweak/internal/pkg/config"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"os"
	"time"

	"clicktweak/internal/app/analyzer"
	"clicktweak/internal/pkg/db/impl"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	if err := config.Unmarshal(&analyzer.Config); err != nil {
		log.Fatalln(err)
	}

	lvl, err := log.ParseLevel(analyzer.Config.Log.Level)
	if err != nil {
		log.Warningln(err)
		log.Warningln("logging level unchanged")
	}

	log.SetLevel(lvl)

	// seed random package
	seed := time.Now().UnixNano()
	log.Debugf("seed is: %d", seed)
	rand.Seed(seed)

	// connect to database
	conn, err := gorm.Open(analyzer.Config.Database.Dialect, analyzer.Config.Database.ConnectionString)
	if err != nil {
		log.Fatalln(err)
	}

	// connect to log database
	logConn, err := sqlx.Connect(analyzer.Config.LogDB.Dialect, analyzer.Config.LogDB.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	var db = conn
	if lvl.String() == "debug" || lvl.String() == "trace" {
		db = conn.Debug()
	}

	// create db instances
	urlDB, err := impl.NewUrlDB(db)
	if err != nil {
		log.Fatalln(err)
	}

	logDB, err := impl.NewLogDB(logConn)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.GET("/stats", analyzer.GetStats(urlDB, logDB), middleware.JWT([]byte(analyzer.Config.App.Secret)))
	e.Logger.Fatal(e.Start(analyzer.Config.App.Listen))
}
