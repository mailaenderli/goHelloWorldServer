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
)


func httpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, span := tracer.Start(ctx, "httpHandler")
	defer span.End()



	query := r.URL.Query()
	name := query.Get("name")
	log.Printf("Received request for %s\n", name)

	span.SetAttributes(attribute.String("httpHandler.name", string(name)))

	w.Write([]byte(CreateGreeting(name, ctx)))
}

func CreateGreeting(name string, ctx context.Context) string {
	if name == "" {
		name = "Guest"
	}

	_, span := tracer.Start(ctx, "createGreeting")
	defer span.End()

	rand.Seed(time.Now().UnixNano())
    sleepTime := rand.Intn(2) // n will be between 0 and 2
    time.Sleep(time.Duration(sleepTime)*time.Second)

	span.SetAttributes(attribute.Int("httpHandler.rndDelay", int(sleepTime)))

	return "Hello, " + name + "\n"
}

func main() {
	// Create Server and Route Handlers
	handler := http.HandlerFunc(httpHandler)
	wrappedHandler := otelhttp.NewHandler(handler, "httpHandler-instrumented")

	mux := http.NewServeMux()
	mux.HandleFunc("/", wrappedHandler)

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
