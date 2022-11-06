package main

import (
	"cycir/internal/channeldata"
	"cycir/internal/driver"
	"cycir/internal/elastics"
	"cycir/internal/models"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/pusher/pusher-http-go"
	"github.com/robfig/cron/v3"
	"gopkg.in/natefinch/lumberjack.v2"
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
	InTest			 bool
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

func main() {
	err := run()
	if err != nil {
		log.Println(err)
	}
}

func run() error {
	gob.Register(models.User{})
	_ = os.Setenv("TZ", "America/Halifax")

	cfg.Domain = os.Getenv("DOMAIN")
	cfg.port, _ = strconv.Atoi(os.Getenv("PORT"))
	cfg.env = os.Getenv("ENV")
	cfg.frontend = os.Getenv("FRONTEND_URL")

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

	infoLog := log.New(nil, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(nil, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog.SetOutput(&lumberjack.Logger{
    Filename:   "./logs/infoLog.log",
    MaxSize:    1,  // megabytes after which new file is created
    MaxBackups: 3,  // number of backups
    MaxAge:     28, //days
	})

	errorLog.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/errorLog.log",
		MaxSize:    1,  // megabytes after which new file is created
		MaxBackups: 3,  // number of backups
		MaxAge:     28, //days
	})
	
	db, err := driver.ConnectPostgres(cfg.db.dsn)
	if err != nil {
		log.Println(err)
	}
	defer db.SQL.Close()

	// start mail channel
	infoLog.Println("Initializing mail channel and worker pool....")
	mailQueue := make(chan channeldata.MailJob, maxWorkerPoolSize)

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
		log.Println(err)
	}
	esrepo := elastics.NewElasticRepository(es)

	// preference map
	repo := models.NewPostgresRepository(db.SQL)
	preferenceMap = make(map[string]string)
	preferences, err := repo.AllPreferences()
	if err != nil {
		errorLog.Println("Cannot read preferences:", err)
		return err
	}

	for _, pref := range preferences {
		preferenceMap[pref.Name] = string(pref.Preference)
	}

	preferenceMap["pusher-host"] = cfg.pusherHost
	preferenceMap["pusher-port"] = cfg.pusherPort
	preferenceMap["pusher-key"] = cfg.pusherKey
	preferenceMap["API"] = cfg.frontend

	// create pusher client
	wsClient = pusher.Client{
		AppID:  cfg.pusherApp,
		Secret: cfg.pusherSecret,
		Key:    cfg.pusherKey,
		Secure: cfg.pusherSecure,
		Host:   fmt.Sprintf("%s:%s", cfg.pusherHost, cfg.pusherPort),
	}

	infoLog.Println("Host", fmt.Sprintf("%s:%s", cfg.pusherHost, cfg.pusherPort))
	infoLog.Println("Secure", cfg.pusherSecure)
	infoLog.Println("Pusher port", cfg.pusherPort)

	// monitoring
	monitorMap := make(map[int]cron.EntryID)
	functionMap := make(map[int]cron.EntryID)

	localZone, _ := time.LoadLocation("Local")
	scheduler := cron.New(cron.WithLocation(localZone), cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))
	app := &application{
		config:    cfg,
		infoLog:   infoLog,
		errorLog:  errorLog,
		version:   version,
		repo:      repo,
		MailQueue: mailQueue,
		esrepo:    esrepo,
		PreferenceMap: preferenceMap,
		WsClient: wsClient,
		MonitorMap: monitorMap,
		FunctionMap: functionMap,
		Scheduler: scheduler,
	}

	// Start the email dispatcher
	infoLog.Println("Starting email dispatcher....")
	dispatcher := NewDispatcher(mailQueue, maxJobMaxWorkers)
	dispatcher.run()

	// start the scheduler
	if app.PreferenceMap["monitoring_live"] == "1" {
		app.Scheduler.Start()
	}

	NewScheduler(app)

	if !cfg.InTest {
		go app.StartMonitoring()
		err = esrepo.CreateIndex(app.config.esIndex)
		if err != nil {
			errorLog.Println(err)
		}
	
		err = app.serve()
		if err != nil {
			errorLog.Println(err)
			return err
		}
	}

	return nil
}
