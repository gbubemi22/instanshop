package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"instashop/internal/database"
)

type Server struct {
	port        int
	db          database.Service
	cleanupFunc func()
}

func NewServer() *http.Server {
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 8080
	}

	db := database.New()

	newServer := &Server{
		port: port,
		db:   db,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
