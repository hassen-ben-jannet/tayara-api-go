package main

import (
	"fmt"
	"log"
	"net/http"

	"go-api/router"
)

func main() {
	r := router.Router()
	// fs := http.FileServer(http.Dir("build"))
	// http.Handle("/", fs)
	fmt.Println("Starting server on the port 8080...")

	log.Fatal(http.ListenAndServe(":8080", r))
}

// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/gorilla/mux"
// )

// func handler(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query()
// 	name := query.Get("name")
// 	if name == "" {
// 		name = "Guest"
// 	}
// 	log.Printf("Received request for %s\n", name)
// 	w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
// }

// func healthHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// }

// func readinessHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// }

// func main() {
// 	// Create Server and Route Handlers
// 	r := mux.NewRouter()

// 	r.HandleFunc("/", handler)
// 	r.HandleFunc("/health", healthHandler)
// 	r.HandleFunc("/readiness", readinessHandler)

// 	srv := &http.Server{
// 		Handler:      r,
// 		Addr:         ":8080",
// 		ReadTimeout:  10 * time.Second,
// 		WriteTimeout: 10 * time.Second,
// 	}

// 	// Start Server
// 	go func() {
// 		log.Println("Starting Server")
// 		if err := srv.ListenAndServe(); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	// Graceful Shutdown
// 	waitForShutdown(srv)
// }

// func waitForShutdown(srv *http.Server) {
// 	interruptChan := make(chan os.Signal, 1)
// 	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

// 	// Block until we receive our signal.
// 	<-interruptChan

// 	// create a deadline to wait for.
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()
// 	srv.Shutdown(ctx)

// 	log.Println("Shutting down")
// 	os.Exit(0)
// }
