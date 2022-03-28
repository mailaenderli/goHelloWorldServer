package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httptrace "github.com/signalfx/signalfx-go-tracing/contrib/net/http"
	"github.com/signalfx/signalfx-go-tracing/ddtrace/tracer" // global tracer
  	"github.com/signalfx/signalfx-go-tracing/tracing" // helper
)

func TraceIdFromCtx(ctx context.Context) (result string) {
	if span, ok := tracer.SpanFromContext(ctx); ok {
	  result = tracer.TraceIDHex(span.Context())
	}
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	log.Printf("Received request for %s\n", name)
	w.Write([]byte(CreateGreeting(name)))
}

func CreateGreeting(name string) string {
	if name == "" {
		name = "Guest"
	}
	return "Hello, " + name + "\n"
}

func main() {
	// Create Server and Route Handlers

	mux := httptrace.NewServeMux(httptrace.WithServiceName("goHelloWorld"))

	mux.HandleFunc("/", handler)

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
