package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	getsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sink_get_total",
			Help: "The total get calls",
		},
		[]string{"status"}, // add label for http status
	)

	postsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sink_post_total",
			Help: "The total post calls",
		},
		[]string{"status"}, // add label for http status
	)
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/{key}", mainHandler)

	log.WithFields(logrus.Fields{
		"port": "9009",
	}).Info("starting http sink")

	srv := &http.Server{
		Handler:      r,
		Addr:         ":9009",
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		get(w, r)
	case "POST":
		post(w, r)
	default:
		log.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("unsupported request")

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("unsupported method: %s\n", r.Method)))
	}

}

func get(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")

	log.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"name":   name,
	}).Info("GET request")

	getsProcessed.WithLabelValues(strconv.Itoa(http.StatusOK)).Inc()

	w.Write([]byte(fmt.Sprintln("[]")))
}

func post(w http.ResponseWriter, r *http.Request) {
	val := rand.Intn(100)

	if val < 20 {
		log.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("return 503")

		postsProcessed.WithLabelValues(strconv.Itoa(http.StatusServiceUnavailable)).Inc()

		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		data, _ := ioutil.ReadAll(r.Body)

		log.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"body":   string(data),
		}).Info("POST request")

		postsProcessed.WithLabelValues(strconv.Itoa(http.StatusCreated)).Inc()

		w.WriteHeader(http.StatusCreated)
	}
}
