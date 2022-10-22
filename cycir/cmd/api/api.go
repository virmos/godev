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

	flag.IntVar(&cfg.port, "port", 4002, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production|maintenance}")
	flag.StringVar(&cfg.db.dsn, "dsn", "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5", "DSN")
	flag.StringVar(&cfg.frontend, "frontend", "http://localhost:4000", "url to front end")

	dbHost := flag.String("dbhost", "localhost", "database host")
	dbPort := flag.String("dbport", "5432", "database port")
	dbUser := flag.String("dbuser", "postgres", "database user")
	dbPass := flag.String("dbpass", "qwerqwer", "database password")
	databaseName := flag.String("db", "temp", "database name")
	dbSsl := flag.String("dbssl", "disable", "database ssl setting")
	cfg.db.dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
		*dbHost,
		*dbPort,
		*dbUser,
		*dbPass,
		*databaseName,
		*dbSsl)

	flag.StringVar(&cfg.pusherHost, "pusherHost", "", "pusher host")
	flag.StringVar(&cfg.pusherPort, "pusherPort", "443", "pusher port")
	flag.StringVar(&cfg.pusherApp, "pusherApp", "9", "pusher app id")
	flag.StringVar(&cfg.pusherKey, "pusherKey", "", "pusher key")
	flag.StringVar(&cfg.pusherSecret, "pusherSecret", "", "pusher secret")
	flag.BoolVar(&cfg.pusherSecure, "pusherSecure", false, "pusher server uses SSL (true or false)")
	flag.StringVar(&cfg.Domain, "domain", "localhost", "domain name (e.g. example.com)")
	flag.BoolVar(&cfg.InProduction, "production", false, "application is in production")

	flag.StringVar(&cfg.esAddress, "esAddress", "http://localhost:9200", "elasticsearch address")
	flag.StringVar(&cfg.esUsername, "esUsername", "elastic", "elasticsearch username")
	flag.StringVar(&cfg.esPassword, "esPassword", "EWAq+EaS8dyQV_82TSQd", "elasticsearch password")
	flag.StringVar(&cfg.esIndex, "esIndex", "my-index-000001", "elasticsearch index")

	flag.Parse()

	e, err := os.OpenFile("./foo.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
			fmt.Printf("error opening file: %v", err)
			os.Exit(1)
	}
	infoLog := log.New(e, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(e, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
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

	// Start the email dispatcher
	infoLog.Println("Starting email dispatcher....")
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

	go app.StartMonitoring()
	NewScheduler(app)

	if app.PreferenceMap["monitoring_live"] == "1" {
		app.Scheduler.Start()
	}

	if !cfg.InTest {
		// err = esrepo.CreateIndex(app.config.esIndex)
		// if err != nil {
		// 	errorLog.Println(err)
		// }
	
		err = app.serve()
		if err != nil {
			errorLog.Println(err)
			return err
		}
	}

	return nil
}
