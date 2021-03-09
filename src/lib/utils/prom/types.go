package prom

import (
	"github.com/bench-routes/bench-routes/src/lib/utils"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Prom ecosystem.
	Namespace          = "bench_routes"
	SubsystemPing      = "ping"
	SubsystemJitter    = "jitter"
	SubsystemFloodPing = "flood_ping"
	SubsystemMonitor   = "monitor"
	// Labels.
	LabelDomain    = "domain_or_ip"
	LabelMethod    = "method"
	LabelURL       = "url"
	LabelPingTypes = "type"
)

// MachineMetrics returns the metrics for a machine after initialization.
func MachineMetrics() *utils.MachineMetrics {
	return &utils.MachineMetrics{
		Ping: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "ping_instant_milliseconds",
				Help:      "Instantaneous ping value of the target (domain or IP)",
			}, []string{LabelDomain, LabelPingTypes},
		),
		PingCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "ping_ops_total",
				Help:      "Total ping operations carried out for the target",
			}, []string{LabelDomain},
		),
		Jitter: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "jitter_instant_milliseconds",
				Help:      "Instantaneous jitter value of the target",
			}, []string{LabelDomain},
		),
		JitterCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "jitter_ops_total",
				Help:      "Total jitter operations carried out for the target",
			}, []string{LabelDomain},
		),
		FPing: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "flood_ping_instant_milliseconds",
				Help:      "Instantaneous ping value of the target",
			}, []string{LabelDomain, LabelPingTypes},
		),
		FPingCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "flood_ping_ops_total",
				Help:      "Total flood-ping operations carried out for the target",
			}, []string{LabelDomain},
		),
	}
}

// EndpointMetrics returns the metrics for API endpoints.
func EndpointMetrics() *utils.EndpointMetrics {
	return &utils.EndpointMetrics{
		ResponseLength: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "response_length",
				Help:      "Total number of characters in http response from the target",
			}, []string{LabelMethod, LabelDomain, LabelURL},
		),
		ResponseDelay: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "response_delay",
				Help:      "Time lapse between sending the http request and receiving the http response",
			}, []string{LabelMethod, LabelDomain, LabelURL},
		),
		StatusCode: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "response_status_code",
				Help:      "Status code of the http response",
			}, []string{LabelMethod, LabelDomain, LabelURL},
		),
		MonitorCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "monitoring_total",
				Help:      "Total number of times http request is sent to the target for monitoring",
			}, []string{LabelMethod, LabelDomain, LabelURL},
		),
	}
}
