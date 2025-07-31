package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	googleUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "google_up",
		Help: "Whether Google.com is reachable (1 = up, 0 = down)",
	})
)

func init() {
	prometheus.MustRegister(googleUp)
}

func checkGoogle() {
	for {
		resp, err := http.Get("https://www.google.com")
		if err != nil || resp.StatusCode != 200 {
			log.Println("Google DOWN:", err)
			googleUp.Set(0)
		} else {
			log.Println("Google is UP")
			googleUp.Set(1)
			resp.Body.Close()
		}
		time.Sleep(10 * time.Second)
	}
}

func main() {
	go checkGoogle()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Serving on :8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
