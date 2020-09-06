package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

type Config struct {
	Port          string
	DBSource      string
	FetchTimeout  time.Duration
	BodySizeLimit string
}

func main() {
	cfg := Config{
		Port:          "8080",
		DBSource:      "postgresql://postgres:postgres@localhost:5432?sslmode=disable",
		FetchTimeout:  5 * time.Second,
		BodySizeLimit: "1MB",
	}

	db := sqlx.MustConnect("postgres", cfg.DBSource)
	client := http.Client{Timeout: cfg.FetchTimeout}
	f := NewFetcher(client, db)
	err := f.StartActive()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.BodyLimit(cfg.BodySizeLimit))

	e.GET("/api/fetcher", selectRequestsHandler(db))
	e.GET("/api/fetcher/:id/history", selectHistoryHandler(db))
	e.POST("/api/fetcher", addRequestHandler(db, f))
	e.DELETE("/api/fetcher/:id", deleteRequestHandler(db, f))

	s := &http.Server{
		Addr:         "localhost:" + cfg.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(e.StartServer(s))

}
