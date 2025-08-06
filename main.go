package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
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

func initTracer() (func(context.Context) error, error) {
	ctx := context.Background()

	// OTLP HTTP exporter
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("jaeger.observability.svc.cluster.local:4318"),
		otlptracehttp.WithInsecure(), // no TLS in local cluster
	)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("google-checker"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

func checkGoogle() {
	tracer := otel.Tracer("google-checker")

	for {
		ctx, span := tracer.Start(context.Background(), "CheckGoogle")

		req, err := http.NewRequestWithContext(ctx, "GET", "https://www.google.com", nil)
		if err != nil {
			log.Println("Error creating request:", err)
			googleUp.Set(0)
			span.RecordError(err)
			span.End()
			time.Sleep(10 * time.Second)
			continue
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			log.Println("Google DOWN:", err)
			googleUp.Set(0)
			span.SetAttributes(attribute.Bool("google.up", false))
			if err != nil {
				span.RecordError(err)
			}
		} else {
			log.Println("Google is UP")
			googleUp.Set(1)
			span.SetAttributes(attribute.Bool("google.up", true))
		}

		if resp != nil {
			if err := resp.Body.Close(); err != nil {
				log.Printf("warning: error closing response body: %v", err)
			}
		}
		span.End()
		time.Sleep(10 * time.Second)
	}
}

func main() {
	go checkGoogle()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Serving on :8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
