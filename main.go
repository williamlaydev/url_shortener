package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"urlShortener/components"
	"urlShortener/handlers"
	"urlShortener/pages"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env file: %s", err)
	}

	connPool := createDbConnection()
	defer connPool.Close()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Handle API routes
	mux.HandleFunc("POST /api/url/shorten", handlers.NewPostShortenUrlHandler(connPool).ServeHTTP)
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent) // Respond with no content
	})

	mux.Handle("GET /home", templ.Handler(pages.Home()))
	mux.Handle("GET /inputBox", templ.Handler(components.UrlInputDisplay()))
	// Place this last
	mux.HandleFunc("GET /{shortUrl}", handlers.NewGetShortenUrlHandler(connPool).ServeHTTP)

	go func() {
		log.Print("Server started")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	gracefulShutdown(server)
}

func createDbConnection() *pgxpool.Pool {
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("env var DATABASE_URL is missing.")
	}

	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	if err != nil {
		log.Fatalf("Failed to parse database config: %v", err)
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		log.Fatalf("Database failed to connect")
	}
	return conn
}

func gracefulShutdown(server *http.Server) {
	// Create a channel to listen for OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Wait for a termination signal
	<-signalChan
	log.Println("Shutting down server...")

	// Create a context with a timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shutdown complete")
}
