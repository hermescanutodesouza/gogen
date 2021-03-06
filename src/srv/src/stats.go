package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tecnickcom/statsd/v2"
)

var stats *statsd.Client

// StatsData store a single stats configuration
type StatsData struct {
	Prefix      string `json:"prefix"`       // StatsD client's string prefix that will be used in every bucket name.
	Network     string `json:"network"`      // Network type used by the StatsD client (i.e. udp or tcp).
	Address     string `json:"address"`      // Network address of the StatsD daemon (ip:port) or just (:port).
	FlushPeriod int    `json:"flush_period"` // How often (in milliseconds) the StatsD client's buffer is flushed.
}

// initStats initialize the StatsD client
func initStats(cfg *StatsData) (err error) {
	stats, err = statsd.New(
		statsd.Prefix(cfg.Prefix),
		statsd.Network(cfg.Network),
		statsd.Address(cfg.Address),
		statsd.FlushPeriod(time.Duration(cfg.FlushPeriod)*time.Millisecond),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error":        err,
			"statsAddress": cfg.Address,
			"statsNetwork": cfg.Network,
		}).Error("Unable to connect to the StatD daemon")
	}
	return err
}

// statsHandler wraps an http.Handler to send StasdMetrics
func statsHandler(path string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, hr *http.Request) {
		stats.Increment(path + ".in")
		defer stats.Increment(path + ".out")
		log.Debug("handler: " + path)
		handler.ServeHTTP(rw, hr)
	})
}
