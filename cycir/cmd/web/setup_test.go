package main

import (
	"cycir/internal/cache"
	"cycir/internal/models"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
)
var testApp *application
var testSession *scs.SessionManager
var testRedisCache *cache.RedisCache
	
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

	testApp = &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		Session:  testSession,
		Cache: testRedisCache,
	}
	// redis
	testRedisCache = testApp.createClientRedisCache()

	NewHelpers(testApp)

	os.Exit(m.Run())
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
