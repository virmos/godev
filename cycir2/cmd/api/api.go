package main

import (
	"cycir/internal/channeldata"
	"cycir/internal/driver"
	"cycir/internal/elastics"
	"cycir/internal/models"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/pusher/pusher-http-go"
	"github.com/robfig/cron/v3"
)

const version = "1.0.0"
const maxWorkerPoolSize = 5
const maxJobMaxWorkers = 5

var preferenceMap map[string]string
var wsClient pusher.Client
var app *application
var cfg config

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	esAddress  string
	esUsername string
	esPassword string
	esIndex    string

	frontend     string
	pusherHost   string
	pusherPort   string
	pusherApp    string
	pusherKey    string
	pusherSecret string
	pusherSecure bool
	Domain       string
	InProduction bool
}

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	version       string
	repo          models.Repository
	esrepo        elastics.Repository
	PusherSecret  string
	MailQueue     chan channeldata.MailJob
	MonitorMap    map[int]cron.EntryID // check server status after 3 minutes
	FunctionMap   map[int]cron.EntryID // send uptime report at 6 am local time
	PreferenceMap map[string]string
	Scheduler     *cron.Cron
	WsClient      pusher.Client
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting Back end server in %s mode on port %d\n", app.config.env, app.config.port)

	return srv.ListenAndServe()
}

func init() {
	gob.Register(models.User{})
	_ = os.Setenv("TZ", "America/Halifax")
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4002, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production|maintenance}")
	flag.StringVar(&cfg.frontend, "frontend", "http://localhost:4000", "url to front end")
	flag.Parse()

	cfg.db.dsn = os.Getenv("DB_DSN")

	cfg.pusherHost = os.Getenv("PUSHER_HOST")
	cfg.pusherPort = os.Getenv("PUSHER_PORT")
	cfg.pusherApp = os.Getenv("PUSHER_APP")
	cfg.pusherKey = os.Getenv("PUSHER_KEY")
	cfg.pusherSecret = os.Getenv("PUSHER_SECRET")
	if (os.Getenv("PUSHER_SECURE") == "disable") {
		cfg.pusherSecure = false
	} else {
		cfg.pusherSecure = true
	}
	if (os.Getenv("IN_PRODUCTION") == "disable") {
		cfg.InProduction = false
	} else {
		cfg.InProduction = true
	}
	cfg.Domain = os.Getenv("DOMAIN")
	cfg.esAddress = os.Getenv("ES_ADDRESS")
	cfg.esUsername = os.Getenv("ES_USERNAME")
	cfg.esPassword = os.Getenv("ES_PASSWORD")
	cfg.esIndex = os.Getenv("ES_INDEX")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := driver.ConnectPostgres(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.SQL.Close()

	// start mail channel
	log.Println("Initializing mail channel and worker pool....")
	mailQueue := make(chan channeldata.MailJob, maxWorkerPoolSize)

	// Start the email dispatcher
	log.Println("Starting email dispatcher....")
	dispatcher := NewDispatcher(mailQueue, maxJobMaxWorkers)
	dispatcher.run()

	// Start elasticsearch
	esCfg := elasticsearch.Config{
		Addresses: []string{
			cfg.esAddress,
		},
		Username: cfg.esUsername,
		Password: cfg.esPassword,
	}
	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		errorLog.Fatal(err)
	}
	esrepo := elastics.NewElasticRepository(es)

	app := &application{
		config:    cfg,
		infoLog:   infoLog,
		errorLog:  errorLog,
		version:   version,
		repo:      models.NewPostgresRepository(db.SQL),
		MailQueue: mailQueue,
		esrepo:    esrepo,
	}

	preferenceMap = make(map[string]string)
	preferences, err := app.repo.AllPreferences()
	if err != nil {
		log.Fatal("Cannot read preferences:", err)
	}

	for _, pref := range preferences {
		preferenceMap[pref.Name] = string(pref.Preference)
	}

	preferenceMap["pusher-host"] = cfg.pusherHost
	preferenceMap["pusher-port"] = cfg.pusherPort
	preferenceMap["pusher-key"] = cfg.pusherKey
	preferenceMap["API"] = cfg.frontend

	app.PreferenceMap = preferenceMap

	// create pusher client
	wsClient = pusher.Client{
		AppID:  cfg.pusherApp,
		Secret: cfg.pusherSecret,
		Key:    cfg.pusherKey,
		Secure: cfg.pusherSecure,
		Host:   fmt.Sprintf("%s:%s", cfg.pusherHost, cfg.pusherPort),
	}

	log.Println("Host", fmt.Sprintf("%s:%s", cfg.pusherHost, cfg.pusherPort))
	log.Println("Secure", cfg.pusherSecure)
	log.Println("Pusher port", cfg.pusherPort)

	app.WsClient = wsClient
	
	monitorMap := make(map[int]cron.EntryID)
	app.MonitorMap = monitorMap
	functionMap := make(map[int]cron.EntryID)
	app.FunctionMap = functionMap

	localZone, _ := time.LoadLocation("Local")
	scheduler := cron.New(cron.WithLocation(localZone), cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))

	app.Scheduler = scheduler

	go app.StartMonitoring()
	NewScheduler(app)

	if app.PreferenceMap["monitoring_live"] == "1" {
		app.Scheduler.Start()
	}
	// err = esrepo.CreateIndex(app.config.esIndex)
	// if err != nil {
	// 	errorLog.Fatal(err)
	// }

	// reports, _ := app.esrepo.GetAllReports("cycir")
	// log.Println(reports)

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
