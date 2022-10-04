package main

import (
	"cycir/internal/cache"
	"cycir/internal/driver"
	"cycir/internal/models"
	"encoding/gob"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

const version = "1.0.0"

var preferenceMap map[string]string
var app *application
var session *scs.SessionManager
var redisCache *cache.RedisCache
var cfg config

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	redis struct {
		prefix   string
		host     string
		password string
	}
	backend      string
	pusherHost   string
	pusherPort   string
	pusherApp    string
	pusherKey    string
	pusherSecret string
	pusherSecure bool
	Identifier   string
	Domain       string
	InProduction bool
}

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	version       string
	repo          models.Repository
	Session       *scs.SessionManager
	PreferenceMap map[string]string
	TemplateCache map[string]*template.Template
	Cache         cache.Cache
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

	app.infoLog.Printf("Starting Front end server in %s mode on port %d\n", app.config.env, app.config.port)

	return srv.ListenAndServe()
}

func init() {
	gob.Register(models.User{})
	_ = os.Setenv("TZ", "America/Halifax")
}

func main() {
	dbHost := flag.String("dbhost", "localhost", "database host")
	dbPort := flag.String("dbport", "5432", "database port")
	dbUser := flag.String("dbuser", "postgres", "database user")
	dbPass := flag.String("dbpass", "qwerqwer", "database password")
	databaseName := flag.String("db", "temp", "database name")
	dbSsl := flag.String("dbssl", "disable", "database ssl setting")

	if *dbUser == "" || *dbHost == "" || *dbPort == "" || *databaseName == "" {
		fmt.Println("Missing database required flags.")
		os.Exit(1)
	}

	cfg.db.dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
		*dbHost,
		*dbPort,
		*dbUser,
		*dbPass,
		*databaseName,
		*dbSsl)

	flag.StringVar(&cfg.redis.prefix, "redisPrefix", "cycir", "redis prefix")
	flag.StringVar(&cfg.redis.host, "redisHost", "localhost:6379", "redis host")
	flag.StringVar(&cfg.redis.password, "redisPass", "", "redis password")

	flag.StringVar(&cfg.pusherHost, "pusherHost", "", "pusher host")
	flag.StringVar(&cfg.pusherPort, "pusherPort", "443", "pusher port")
	flag.StringVar(&cfg.pusherApp, "pusherApp", "9", "pusher app id")
	flag.StringVar(&cfg.pusherKey, "pusherKey", "", "pusher key")
	flag.StringVar(&cfg.pusherSecret, "pusherSecret", "", "pusher secret")
	flag.BoolVar(&cfg.pusherSecure, "pusherSecure", false, "pusher server uses SSL (true or false)")

	flag.StringVar(&cfg.Identifier, "identifier", "cycir", "unique identifier")
	flag.StringVar(&cfg.Domain, "domain", "localhost", "domain name (e.g. example.com)")
	flag.BoolVar(&cfg.InProduction, "production", false, "application is in production")

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production|maintenance}")
	flag.StringVar(&cfg.backend, "backend", "http://localhost:4002", "url to back end")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// database
	db, err := driver.ConnectPostgres(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.SQL.Close()

	// session
	log.Printf("Initializing session manager....")

	session = scs.New()
	session.Store = postgresstore.New(db.SQL)
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Name = fmt.Sprintf("gbsession_id_%s", cfg.Identifier)
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = cfg.InProduction

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		repo:     models.NewPostgresRepository(db.SQL),
		Session:  session,
	}

	// redis
	redisCache = app.createClientRedisCache()
	app.Cache = redisCache

	// preference map
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
	preferenceMap["API"] = cfg.backend
	preferenceMap["version"] = app.version

	app.PreferenceMap = preferenceMap
	NewHelpers(app)

	log.Println("Host", fmt.Sprintf("%s:%s", cfg.pusherHost, cfg.pusherPort))
	log.Println("Secure", cfg.pusherSecure)

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) createClientRedisCache() *cache.RedisCache {
	cacheClient := cache.RedisCache{
		Conn:   app.createRedisPool(),
		Prefix: app.config.redis.prefix,
	}
	return &cacheClient
}

func (app *application) createRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",
				app.config.redis.host,
				redis.DialPassword(app.config.redis.password))
		},

		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}
