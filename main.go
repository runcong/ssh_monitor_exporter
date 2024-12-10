package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var listenAddress = flag.String(
	"listen-address",
	":9188",
	"The address to listen on for HTTP requests.")

var reg = prometheus.NewRegistry()

func registerMetrics() {
	reg.MustRegister(ssh_mon_status)

}

func main() {
	registerMetrics()
	flag.Parse()
	log.Printf("Starting Server: %s", *listenAddress)

	// http.Handle("/metrics", promhttp.Handler())
	http.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		go check_ssh_status()

		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	}))
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
