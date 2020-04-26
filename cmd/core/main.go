package main

import (
	"clicktweak/internal/pkg/config"
	"math/rand"
	"os"
	"time"

	"clicktweak/internal/app/core"
	"clicktweak/internal/pkg/db/impl"

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
	if err := config.Unmarshal(&core.Config); err != nil {
		log.Fatalln(err)
	}

	lvl, err := log.ParseLevel(core.Config.Log.Level)
	if err != nil {
		log.Warningln(err)
		log.Warningln("logging level unchanged")
	}

	log.SetLevel(lvl)

	// seed random package
	seed := time.Now().UnixNano()
	log.Debugf("seed is: %d", seed)
	rand.Seed(seed)

	conn, err := gorm.Open(core.Config.Database.Dialect, core.Config.Database.ConnectionString)
	if err != nil {
		log.Fatalln(err)
	}

	var db = conn
	if lvl.String() == "debug" || lvl.String() == "trace" {
		db = conn.Debug()
	}

	userDB, err := impl.NewUserDB(db)
	if err != nil {
		log.Fatalln(err)
	}

	urlDB, err := impl.NewUrlDB(db)
	if err != nil {
		log.Fatalln(err)
	}

	e := echo.New()
	e.POST("/shorten", core.Shorten(urlDB), middleware.JWT([]byte(core.Config.App.Secret)))
	e.POST("/login", core.Login(userDB, core.Config.App.Secret))
	e.POST("/signup", core.SignUp(userDB, core.Config.App.Secret))

	e.Logger.Fatal(e.Start(core.Config.App.Listen))
}
