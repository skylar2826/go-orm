package aop

import (
	"context"
	"fmt"
	"geektime-go-orm/orm"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	trace2 "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"time"
)

func init() {
	initZipkin("", "server")
}

func initZipkin(collectorUrl string, serviceName string) {
	exporter, err := zipkin.New(collectorUrl, zipkin.WithLogger(log.New(os.Stderr, "zipkin-server", log.Llongfile)))
	if err != nil {
		log.Println(err.Error())
		return
	}
	b := trace2.NewBatchSpanProcessor(exporter)
	tp := trace2.NewTracerProvider(trace2.WithSpanProcessor(b), trace2.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(fmt.Sprintf("zipkin-%s", serviceName)))))
	otel.SetTracerProvider(tp)
}

type ObserverMiddlewareBuilder struct {
	logFunc func(query *orm.Query)
	tracer  trace.Tracer
	vector  *prometheus.SummaryVec
}

func (m *ObserverMiddlewareBuilder) setTracerSpan(ctx context.Context, qc *orm.QueryContext, query *orm.Query) trace.Span {
	if m.tracer == nil {
		return nil
	}

	tableName := qc.Model.TableName

	var span trace.Span
	ctx, span = m.tracer.Start(ctx, qc.Type+"-"+tableName)
	span.SetAttributes(attribute.String("table", tableName))
	span.SetAttributes(attribute.String("sql", query.SQL))

	return span
}

func (m *ObserverMiddlewareBuilder) observeVector(ctx context.Context, qc *orm.QueryContext, startTime time.Time) {
	go func() {
		if m.vector != nil {
			m.vector.WithLabelValues(qc.Type, qc.Model.TableName).Observe(float64(time.Now().Sub(startTime).Milliseconds()))
		}
	}()
}

func (m *ObserverMiddlewareBuilder) Build() Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			startTime := time.Now()
			query, err := qc.Builder.Build()
			span := m.setTracerSpan(ctx, qc, query)
			if err != nil {
				return &orm.QueryResult{
					Err: err,
				}
			}
			m.logFunc(query)
			m.observeVector(ctx, qc, startTime)
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			return next(ctx, qc)
		}
	}
}

type Opt func(m *ObserverMiddlewareBuilder)

func WithLogFunc(logFunc func(query *orm.Query)) Opt {
	return func(m *ObserverMiddlewareBuilder) {
		m.logFunc = logFunc
	}
}

var defaultInstrumentationName = "D:\\workspace\\geektime-go-orm\\orm\\aop\\observer.go"

func WithTracer() Opt {
	return func(m *ObserverMiddlewareBuilder) {
		m.tracer = otel.GetTracerProvider().Tracer(defaultInstrumentationName)
	}
}

func WithVector(Namespace string, Subsystem string, Name string, Help string) Opt {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: Namespace,
		Subsystem: Subsystem,
		Name:      Name,
		Help:      Help,
	}, []string{"type", "table"})

	return func(m *ObserverMiddlewareBuilder) {
		m.vector = vector
		prometheus.MustRegister(vector)
	}
}

func NewObserverMiddleBuilder(opts ...Opt) *ObserverMiddlewareBuilder {
	m := &ObserverMiddlewareBuilder{
		logFunc: func(query *orm.Query) {
			fmt.Println("observer middleware: ", query)
		},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}
