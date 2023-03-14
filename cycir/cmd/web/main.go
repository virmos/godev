package main

import (
	"cycir/internal/cache"
	"cycir/internal/driver"
	"cycir/internal/models"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"gopkg.in/natefinch/lumberjack.v2"
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
		port     string
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
	InTest       bool
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

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	gob.Register(models.User{})
	_ = os.Setenv("TZ", "America/Halifax")

	cfg.Domain = os.Getenv("DOMAIN")
	cfg.Identifier = os.Getenv("IDENTIFIER")
	cfg.port, _ = strconv.Atoi(os.Getenv("PORT"))
	cfg.env = os.Getenv("ENV")
	cfg.backend = os.Getenv("BACKEND_URL")

	cfg.db.dsn = os.Getenv("DB_DSN")
	cfg.pusherHost = os.Getenv("PUSHER_HOST")
	cfg.pusherPort = os.Getenv("PUSHER_PORT")
	cfg.pusherApp = os.Getenv("PUSHER_APP")
	cfg.pusherKey = os.Getenv("PUSHER_KEY")
	cfg.pusherSecret = os.Getenv("PUSHER_SECRET")
	if os.Getenv("PUSHER_SECURE") == "disable" {
		cfg.pusherSecure = false
	} else {
		cfg.pusherSecure = true
	}
	if os.Getenv("IN_PRODUCTION") == "disable" {
		cfg.InProduction = false
	} else {
		cfg.InProduction = true
	}
	cfg.Domain = os.Getenv("DOMAIN")
	cfg.redis.host = os.Getenv("REDIS_HOST")
	cfg.redis.port = os.Getenv("REDIS_PORT")
	cfg.redis.prefix = os.Getenv("REDIS_PREFIX")

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
	// database
	db, err := driver.ConnectPostgres(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.SQL.Close()

	// session
	infoLog.Printf("Initializing session manager....")

	session = scs.New()
	session.Store = postgresstore.New(db.SQL)
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Name = fmt.Sprintf("gbsession_id_%s", cfg.Identifier)
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = cfg.InProduction

	// preference map
	preferenceMap = make(map[string]string)
	repo := models.NewPostgresRepository(db.SQL)
	preferences, err := repo.AllPreferences()
	if err != nil {
		errorLog.Println("Cannot read preferences:", err)
	}

	for _, pref := range preferences {
		preferenceMap[pref.Name] = string(pref.Preference)
	}

	preferenceMap["pusher-host"] = cfg.pusherHost
	preferenceMap["pusher-port"] = cfg.pusherPort
	preferenceMap["pusher-key"] = cfg.pusherKey
	preferenceMap["API"] = cfg.backend
	preferenceMap["version"] = version

	infoLog.Println("Host", fmt.Sprintf("%s:%s", cfg.pusherHost, cfg.pusherPort))
	infoLog.Println("Secure", cfg.pusherSecure)

	if !cfg.InTest {
		app := &application{
			config:        cfg,
			infoLog:       infoLog,
			errorLog:      errorLog,
			version:       version,
			repo:          repo,
			Session:       session,
			PreferenceMap: preferenceMap,
		}
		NewHelpers(app)

		// redis
		redisCache := app.createClientRedisCache()
		app.Cache = redisCache

		err = app.serve()
		if err != nil {
			errorLog.Fatal(err)
			return err
		}
	}
	return nil
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
				fmt.Sprintf("%s:%s", app.config.redis.host, app.config.redis.port),
				redis.DialPassword(app.config.redis.password))
		},

		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}
