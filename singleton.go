package boxed

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type singleton struct {
	DbConn      *pgxpool.Pool
	BackendPort int
	FolderPath  string
	JwtSecret   string
}

var (
	instance *singleton
	once     sync.Once
)

// GetInstance method will return the only instance of the singleton struct
func GetInstance() *singleton {
	once.Do(func() {
		// Load needed variables
		if err := godotenv.Load(); err != nil {
			log.Fatal("Couldn't load the .env file. Info:", err)
		}
		dbUrl := os.Getenv("DB_URL")
		if dbUrl == "" {
			log.Fatal("DB_URL is empty")
		}
		folderPath := os.Getenv("FOLDER_PATH")
		if folderPath == "" {
			log.Fatal("FOLDER_PATH is empty")
		}

		backendPortRaw := os.Getenv("BACKEND_PORT")
		if backendPortRaw == "" {
			log.Fatal("BACKEND_PORT is empty")
		}
		jwtSecret := os.Getenv("JWT_SECRET")
		if backendPortRaw == "" {
			log.Fatal("JWT_SECRET is empty")
		}
		backendPort, err := strconv.Atoi(backendPortRaw)
		if err != nil {
			log.Fatal("Error while converting the backendPort to an integer")
		}
		// Make the connection
		config, err := pgxpool.ParseConfig(dbUrl)
		if err != nil {
			log.Fatal("Error while parsing the dbUrl to the pgxpool. Info:", err)
		}
		config.MaxConns = 3
		config.MinConns = 1
		con, err := pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			log.Fatal("Error while connecting to the database. Info:", err)
		}
		instance = &singleton{
			DbConn:      con,
			BackendPort: backendPort,
			FolderPath:  folderPath,
			JwtSecret:   jwtSecret,
		}
	})
	return instance
}
