package metrics

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
)

type GlobalRegistry interface {
	APILatencies() instrument.Int64Histogram
	StatusCodes() instrument.Int64Counter
	ActiveLedgers() instrument.Int64UpDownCounter
}

type globalRegistry struct {
	// API Latencies
	apiLatencies  instrument.Int64Histogram
	statusCodes   instrument.Int64Counter
	activeLedgers instrument.Int64UpDownCounter
}

func RegisterGlobalRegistry(meterProvider metric.MeterProvider) (GlobalRegistry, error) {
	meter := meterProvider.Meter("global")

	apiLatencies, err := meter.Int64Histogram(
		"ledger.api.time",
		instrument.WithUnit("ms"),
		instrument.WithDescription("Latency of API calls"),
	)
	if err != nil {
		return nil, err
	}

	statusCodes, err := meter.Int64Counter(
		"ledger.api.status",
		instrument.WithUnit("1"),
		instrument.WithDescription("Status codes of API calls"),
	)
	if err != nil {
		return nil, err
	}

	activeLedgers, err := meter.Int64UpDownCounter(
		"ledger.api.ledgers",
		instrument.WithUnit("1"),
		instrument.WithDescription("Number of active ledgers"),
	)
	if err != nil {
		return nil, err
	}

	return &globalRegistry{
		apiLatencies:  apiLatencies,
		statusCodes:   statusCodes,
		activeLedgers: activeLedgers,
	}, nil
}

func (gm *globalRegistry) APILatencies() instrument.Int64Histogram {
	return gm.apiLatencies
}

func (gm *globalRegistry) StatusCodes() instrument.Int64Counter {
	return gm.statusCodes
}

func (gm *globalRegistry) ActiveLedgers() instrument.Int64UpDownCounter {
	return gm.activeLedgers
}

type PerLedgerRegistry interface {
	CacheNumberEntries() instrument.Int64UpDownCounter

	QueryLatencies() instrument.Int64Histogram
	QueryInboundLogs() instrument.Int64Counter
	QueryPendingMessages() instrument.Int64Counter
	QueryProcessedLogs() instrument.Int64Counter
	LockWaitingTime() instrument.Int64Histogram
}

type perLedgerRegistry struct {
	cacheNumberEntries instrument.Int64UpDownCounter

	queryLatencies       instrument.Int64Histogram
	queryInboundLogs     instrument.Int64Counter
	queryPendingMessages instrument.Int64Counter
	queryProcessedLogs   instrument.Int64Counter
	lockWaitingTime      instrument.Int64Histogram
}

func RegisterPerLedgerMetricsRegistry(ledger string) (PerLedgerRegistry, error) {
	// we can now use the global meter provider to create a meter
	// since it was created by the fx
	meter := global.MeterProvider().Meter(ledger)

	cacheNumberEntries, err := meter.Int64UpDownCounter(
		"ledger.cache.pending_entries",
		instrument.WithUnit("1"),
		instrument.WithDescription("Number of entries in the cache"),
	)
	if err != nil {
		return nil, err
	}

	queryLatencies, err := meter.Int64Histogram(
		"ledger.query.time",
		instrument.WithUnit("ms"),
		instrument.WithDescription("Latency of queries processing logs"),
	)
	if err != nil {
		return nil, err
	}

	queryInboundLogs, err := meter.Int64Counter(
		"ledger.query.inbound_logs",
		instrument.WithUnit("1"),
		instrument.WithDescription("Number of inbound logs in CQRS worker"),
	)
	if err != nil {
		return nil, err
	}

	queryPendingMessages, err := meter.Int64Counter(
		"ledger.query.pending_messages",
		instrument.WithUnit("1"),
		instrument.WithDescription("Number of pending messages in CQRS worker"),
	)
	if err != nil {
		return nil, err
	}

	queryProcessedLogs, err := meter.Int64Counter(
		"ledger.query.processed_logs",
		instrument.WithUnit("1"),
		instrument.WithDescription("Number of processed logs in CQRS worker"),
	)
	if err != nil {
		return nil, err
	}

	lockWaitingTime, err := meter.Int64Histogram(
		"ledger.lock.waiting_time",
		instrument.WithUnit("ms"),
		instrument.WithDescription("Latency of lock"),
	)
	if err != nil {
		return nil, err
	}

	return &perLedgerRegistry{
		cacheNumberEntries:   cacheNumberEntries,
		queryLatencies:       queryLatencies,
		queryInboundLogs:     queryInboundLogs,
		queryPendingMessages: queryPendingMessages,
		queryProcessedLogs:   queryProcessedLogs,
		lockWaitingTime:      lockWaitingTime,
	}, nil
}

func (pm *perLedgerRegistry) LockWaitingTime() instrument.Int64Histogram {
	return pm.lockWaitingTime
}

func (pm *perLedgerRegistry) CacheNumberEntries() instrument.Int64UpDownCounter {
	return pm.cacheNumberEntries
}

func (pm *perLedgerRegistry) QueryLatencies() instrument.Int64Histogram {
	return pm.queryLatencies
}

func (pm *perLedgerRegistry) QueryInboundLogs() instrument.Int64Counter {
	return pm.queryInboundLogs
}

func (pm *perLedgerRegistry) QueryPendingMessages() instrument.Int64Counter {
	return pm.queryPendingMessages
}

func (pm *perLedgerRegistry) QueryProcessedLogs() instrument.Int64Counter {
	return pm.queryProcessedLogs
}

type noOpRegistry struct{}

func NewNoOpRegistry() *noOpRegistry {
	return &noOpRegistry{}
}

func (nm *noOpRegistry) LockWaitingTime() instrument.Int64Histogram {
	histogram, _ := metric.NewNoopMeter().Int64Histogram("waiting_time")
	return histogram
}

func (nm *noOpRegistry) CacheMisses() instrument.Int64Counter {
	counter, _ := metric.NewNoopMeter().Int64Counter("cache_misses")
	return counter
}

func (nm *noOpRegistry) CacheNumberEntries() instrument.Int64UpDownCounter {
	counter, _ := metric.NewNoopMeter().Int64UpDownCounter("cache_number_entries")
	return counter
}

func (nm *noOpRegistry) QueryLatencies() instrument.Int64Histogram {
	histogram, _ := metric.NewNoopMeter().Int64Histogram("query_latencies")
	return histogram
}

func (nm *noOpRegistry) QueryInboundLogs() instrument.Int64Counter {
	counter, _ := metric.NewNoopMeter().Int64Counter("query_inbound_logs")
	return counter
}

func (nm *noOpRegistry) QueryPendingMessages() instrument.Int64Counter {
	counter, _ := metric.NewNoopMeter().Int64Counter("query_pending_messages")
	return counter
}

func (nm *noOpRegistry) QueryProcessedLogs() instrument.Int64Counter {
	counter, _ := metric.NewNoopMeter().Int64Counter("query_processed_logs")
	return counter
}

func (nm *noOpRegistry) APILatencies() instrument.Int64Histogram {
	histogram, _ := metric.NewNoopMeter().Int64Histogram("api_latencies")
	return histogram
}

func (nm *noOpRegistry) StatusCodes() instrument.Int64Counter {
	counter, _ := metric.NewNoopMeter().Int64Counter("status_codes")
	return counter
}

func (nm *noOpRegistry) ActiveLedgers() instrument.Int64UpDownCounter {
	counter, _ := metric.NewNoopMeter().Int64UpDownCounter("active_ledgers")
	return counter
}
