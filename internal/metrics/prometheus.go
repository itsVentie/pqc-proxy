package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ActiveConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "active_connections",
		Help: "Number of active connections",
	})

	BytesTransferred = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "bytes_transferred",
		Help: "Total bytes transferred",
	}, []string{"direction"})

	ErrorsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pqc_proxy_errors_total",
		Help: "Total number of errors",
	})
)

func Init() {
	prometheus.MustRegister(ActiveConnections)
	prometheus.MustRegister(BytesTransferred)
	prometheus.MustRegister(ErrorsTotal)
}

func StartServer(addr string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, nil)
}
