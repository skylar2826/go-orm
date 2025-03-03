package aop

import (
	"context"
	"geektime-go-orm/orm/sql"
	"log"
)

type MonitorMiddlewareBuilder struct {
	logFunc func(query *sql.Query)
}

func (m *MonitorMiddlewareBuilder) Build() Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, qc *QueryContext) *QueryResult {
			query, err := qc.Builder.Build()
			if err != nil {
				return &QueryResult{
					Err: err,
				}
			}
			m.logFunc(query)
			return next(ctx, qc)
		}
	}
}

type Opt func(m *MonitorMiddlewareBuilder)

func WithLogFunc(logFunc func(query *sql.Query)) Opt {
	return func(m *MonitorMiddlewareBuilder) {
		m.logFunc = logFunc
	}
}

func NewMonitorMiddle(opts ...Opt) *MonitorMiddlewareBuilder {
	m := &MonitorMiddlewareBuilder{
		logFunc: func(query *sql.Query) {
			log.Println(query)
		},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}
