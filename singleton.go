package boxed

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type singleton struct {
	DbConn      *pgx.Conn
	BackendHost string
	BackendPort int
}

var lock = &sync.Mutex{}
var instance *singleton

func GetInstance() *singleton {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		// Load needed variables
		if err := godotenv.Load(); err != nil {
			log.Fatal("Couldn't load the .env file. Info:", err)
		}
		dbUrl := os.Getenv("DB_URL")
		if dbUrl == "" {
			log.Fatal("DB_URL is empty")
		}
		backendHost := os.Getenv("BACKEND_HOST")
		if backendHost == "" {
			log.Fatal("BACKEND_HOST is empty")
		}
		backendPortRaw := os.Getenv("BACKEND_PORT")
		if backendPortRaw == "" {
			log.Fatal("BACKEND_PORT is empty")
		}
		backendPort, err := strconv.Atoi(backendPortRaw)
		if err != nil {
			log.Fatal("Error while converting the backendPort to an integer")
		}
		// Make the connection
		con, err := pgx.Connect(context.TODO(), dbUrl)
		if err != nil {
			log.Fatal("Error while connecting to the database. Info:", err)
		}
		instance = &singleton{
			DbConn:      con,
			BackendHost: backendHost,
			BackendPort: backendPort,
		}
		return instance
	}
	return instance
}
