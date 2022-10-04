package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pi6atv/digital-receivers-status/pkg/receivers"
	"github.com/pi6atv/digital-receivers-status/pkg/receivers/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	configPath          = flag.String("config", "drx.yaml", "config file for the exporter")
	listen              = flag.String("listen", ":9001", "where to listen on")
	requestMetrics      = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "drx_web_requests_total", Help: "web requests to the exporter"}, []string{"path"})
	requestErrorMetrics = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "drx_web_errors_total", Help: "web request errors from the exporter"}, []string{"path"})
)

func main() {
	flag.Parse()

	conf := config.ReceiversConfig{}
	yamlFile, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Failed reading config file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Failed parsing config file: %v", err)
	}

	recs, err := receivers.New(conf)
	if err != nil {
		log.Fatalf("Failed setting up receivers: %v", err)
	}

	log.Printf("Starting scrapers")
	recs.StartScrapers()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		requestMetrics.WithLabelValues("json").Inc()
		w.Header().Set("content-type", "application/json")

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			w.WriteHeader(http.StatusOK)
			gz := gzip.NewWriter(w)
			json.NewEncoder(gz).Encode(recs.Receivers)
			gz.Close()
		} else {
			data, err := json.Marshal(recs.Receivers)
			if err != nil {
				requestErrorMetrics.WithLabelValues("json").Inc()
				log.Printf("ERROR: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-length", fmt.Sprintf("%d", len(data)))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(data)
		}

	})

	log.Fatal(http.ListenAndServe(*listen, nil))
}
