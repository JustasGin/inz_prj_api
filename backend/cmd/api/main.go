package main

import (
	"backend/models"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	jwt struct {
		secret string
	}
}

type AppStatus struct {
	Status      string `json:"status" bson:"status"`
	Environment string `json:"environment" bson:"environment"`
	Version     string `json:"version" bson:"version"`
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func main() {
	var cfg config

	port, err := strconv.Atoi(os.Getenv("port"))
	if err != nil {
		log.Println(err)
		return
	}

	flag.IntVar(&cfg.port, "port", port, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENVIRONMENT"), "application environment")
	flag.StringVar(&cfg.db.dsn, "dns", fmt.Sprintf("postgres://%s@localhost/%s?sslmode=disable", os.Getenv("DB_LOGIN"), os.Getenv("DB_NAME")), "Postgres connection string")
	flag.Parse()

	cfg.jwt.secret = os.Getenv("GO_MOVIES_JWT")

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Println("Starting server on port", cfg.port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
