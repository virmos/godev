package main

import (
	"cycir/internal/elastics"
	"cycir/internal/channeldata"
	"cycir/internal/models"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/pusher/pusher-http-go"
	"github.com/robfig/cron/v3"
)

var testApp *application

func TestMain(m *testing.M) {
	gob.Register(models.User{})
	_ = os.Setenv("TZ", "America/Halifax")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	cfg.InTest = true
	preferenceMap = make(map[string]string)
	wsClient := pusher.Client{
		AppID:  "test",
		Secret: "test",
		Key:    "test",
		Secure: false,
		Host:   "test",
	}

	localTime, _ := time.LoadLocation("Local")
	scheduler := cron.New(cron.WithLocation(localTime), cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))
	
	NewMailTestPath()
	mailQueue := make(chan channeldata.MailJob, maxWorkerPoolSize)
	dispatcher := NewDispatcher(mailQueue, maxJobMaxWorkers)
	dispatcher.run()

	testApp = &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		version:       version,
		WsClient:      wsClient,
		PreferenceMap: preferenceMap,
		Scheduler:     scheduler,
		repo:          models.NewTestRepository(),
		esrepo:        elastics.NewTestElasticRepository(),
		MailQueue: mailQueue,
	}

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
