package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"math/rand"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
)

var tracer = otel.Tracer("github.com/mailaenderli/goHelloWorldServer")

func httpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	
	query := r.URL.Query()
	name := query.Get("name")
	log.Printf("Received request for %s\n", name)

	greeting := CreateGreeting(name)

	span.SetAttributes(attribute.String("httpHandler.greeting", string(greeting)))
	
	RndDelay(ctx)

	w.Write([]byte(greeting))
}

func RndDelay(ctx context.Context) {
	_, span := tracer.Start(ctx, "RndDelay")
	defer span.End()

	rand.Seed(time.Now().UnixNano())
    sleepTime := rand.Intn(2) // n will be between 0 and 2
    time.Sleep(time.Duration(sleepTime)*time.Second)

	span.SetAttributes(attribute.Int("httpHandler.rndDelay", int(sleepTime)))
}

func CreateGreeting(name string) string {
	if name == "" {
		name = "Guest"
	}

	return "Hello, " + name + "\n"
}

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("otlptrace-example"),
		semconv.ServiceVersionKey.String("0.0.1"),
	)
}

func installExportPipeline(ctx context.Context) func() {
	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatalf("creating OTLP trace exporter: %v", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource()),
	)
	otel.SetTracerProvider(tracerProvider)

	return func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Fatalf("stopping tracer provider: %v", err)
		}
	}
}

func main() {
	ctx := context.Background()
	// Registers a tracer Provider globally.
	cleanup := installExportPipeline(ctx)
	defer cleanup()

	// Create Server and Route Handlers
	handler := http.HandlerFunc(httpHandler)
	wrappedHandler := otelhttp.NewHandler(handler, "httpHandler-instrumented")

	mux := http.NewServeMux()
	mux.Handle("/", wrappedHandler)

	srv := &http.Server{
		Handler:      mux,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
