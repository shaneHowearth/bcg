package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	cgrpc "github.com/shanehowearth/bcg/customer/integration/grpc/client/v1"
	ngrpc "github.com/shanehowearth/bcg/notify/integration/grpc/client/v1"
)

// Routes -
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,          // Log API request calls
		middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
	)

	bcgRoutes(router)
	return router
}

// Customer instance
var rc cgrpc.CustomerClient

// Notify instance
var notify ngrpc.NotifyClient

func main() {
	router := Routes()

	// Walk all the routes and log them
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("Method: %s Route: %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	// Customer instance
	customerURI, found := os.LookupEnv("CustomerURI")
	if !found {
		log.Fatal("No CustomerURI set, cannot continue")
	}
	rc = cgrpc.CustomerClient{Address: customerURI}

	// Notify instance
	notifyURI, found := os.LookupEnv("NotifyURI")
	if !found {
		log.Fatal("No NotifyURI set, cannot continue")
	}
	notify = ngrpc.NotifyClient{Address: notifyURI}

	portNum := os.Getenv("PORT_NUM")
	server := &http.Server{Addr: "0.0.0.0:" + portNum, Handler: router}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Panicf("Listen and serve returned error: %v", err)
		}
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown returned error %v", err)
	}
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"message": msg})
}

// respondwithJSON write json response format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		// log the error
		log.Printf("writing response generated error: %v", err)
	}
}
