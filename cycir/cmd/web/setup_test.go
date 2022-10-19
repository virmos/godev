package main

import (
	"cycir/internal/cache"
	"cycir/internal/models"
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/justinas/nosurf"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
	"context"
)

var testApp *application
var testSession *scs.SessionManager
var testRedisCache cache.RedisCache

func TestMain(m *testing.M) {
	gob.Register(models.User{})
	_ = os.Setenv("TZ", "America/Halifax")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// session
	testSession = scs.New()
	testSession.Lifetime = 24 * time.Hour
	testSession.Cookie.Persist = true
	testSession.Cookie.Name = fmt.Sprintf("gbsession_id_%s", cfg.Identifier)
	testSession.Cookie.SameSite = http.SameSiteLaxMode
	testSession.Cookie.Secure = cfg.InProduction
	cfg.InTest = true

	testApp = &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		Session:  testSession,
		repo:     models.NewTestRepository(),
	}
	// redis
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	pool := redis.Pool{
		MaxIdle:     50,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", s.Addr())
		},
	}

	testRedisCache.Conn = &pool
	testRedisCache.Prefix = "test-cycir"
	testApp.Cache = &testRedisCache

	defer testRedisCache.Conn.Close()

	NewHelpers(testApp)
	SetViews("./views")

	os.Exit(m.Run())
}
func getRoutes() http.Handler {
	mux := chi.NewRouter()
	// default middleware
	mux.Use(getSessionLoad)
	mux.Use(getNoSurf)

	// login
	mux.Get("/", testApp.LoginScreen)
	mux.Post("/", testApp.Login)

	mux.Get("/user/logout", testApp.Logout)

	// redis cache
	mux.Post("/admin/save-in-cache", testApp.SaveInCache)
	mux.Post("/admin/get-from-cache", testApp.GetFromCache)
	mux.Post("/admin/delete-from-cache", testApp.DeleteFromCache)
	mux.Post("/admin/empty-cache", testApp.EmptyCache)

	// overview
	mux.Get("/admin/overview", testApp.AdminDashboard)

	// events
	mux.Get("/admin/events", testApp.Events)

	// schedule
	mux.Get("/admin/schedule", testApp.ListEntries)

	// settings
	mux.Get("/admin/settings", testApp.Settings)

	// service status pages (all hosts)
	mux.Get("/admin/all-healthy", testApp.AllHealthyServices)
	mux.Get("/admin/all-warning", testApp.AllWarningServices)
	mux.Get("/admin/all-problems", testApp.AllProblemServices)
	mux.Get("/admin/all-pending", testApp.AllPendingServices)

	// users
	mux.Get("/admin/users", testApp.AllUsers)
	mux.Get("/admin/user/{id}", testApp.OneUser)

	// hosts
	mux.Get("/admin/host/all", testApp.AllHosts)
	mux.Get("/admin/host/{id}", testApp.Host)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

func getSessionLoad(next http.Handler) http.Handler {
	return testSession.LoadAndSave(next)
}

func getNoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.ExemptPath("/pusher/auth")
	csrfHandler.ExemptPath("/pusher/hook")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   cfg.InProduction,
		SameSite: http.SameSiteStrictMode,
		Domain:   cfg.Domain,
	})

	return csrfHandler
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = testSession.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

// gets the context
func getCtx(req *http.Request) context.Context {
	ctx, err := testApp.Session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
