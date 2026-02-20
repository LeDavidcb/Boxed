package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/David/Boxed/internal/common/utils"
)

type EnvStruct struct {
	DBUser     string
	DBPassword string
	DBhost     string
	DBPort     string
	BPort      string
	FolderPath string
	JWTSecret  string
}

func main() {
	env := EnvStruct{}                  // Initialize struct
	reader := bufio.NewReader(os.Stdin) // To handle input with spaces

	fmt.Print("Database user: ")
	input, err := reader.ReadString('\n')
	logerr(err)
	env.DBUser = strings.TrimSpace(input)

	fmt.Print("Password: ")
	input, err = reader.ReadString('\n')
	logerr(err)
	env.DBPassword = strings.TrimSpace(input)

	fmt.Print("database host (Default: localhost): ")
	input, err = reader.ReadString('\n')
	logerr(err)
	env.DBhost = strings.TrimSpace(input)
	if env.DBhost == "" {
		env.DBhost = "localhost"
	}

	fmt.Print("Database Port (Default: 5432): ")
	input, err = reader.ReadString('\n')
	logerr(err)
	env.DBPort = strings.TrimSpace(input)
	if env.DBPort == "" {
		env.DBPort = "5432"
	}

	fmt.Print("Backend Port (Default: 8080): ")
	input, err = reader.ReadString('\n')
	logerr(err)
	env.BPort = strings.TrimSpace(input)
	if env.BPort == "" {
		env.BPort = "8080"
	}

	fmt.Print("Folder Path: ")
	input, err = reader.ReadString('\n')
	logerr(err)
	env.FolderPath = strings.TrimSpace(input)
	if strings.HasPrefix(env.FolderPath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Unable to resolve home directory: %v", err)
		}
		env.FolderPath = strings.Replace(env.FolderPath, "~", homeDir, 1)
	}

	fmt.Print("JWT secret password (Leave blank if you want a random generated secret password): ")
	input, err = reader.ReadString('\n')
	logerr(err)
	env.JWTSecret = strings.TrimSpace(input)
	if env.JWTSecret == "" {
		random, err := utils.GenerateRTHash(32)
		logerr(err)
		env.JWTSecret = random
	}

	// Print struct with field names and values
	s := fmt.Sprintf(`DB_URL=postgresql://%v:%v@%v:%v/
BACKEND_PORT=%v
FOLDER_PATH=%v
JWT_SECRET=%v`, env.DBUser, env.DBPassword, env.DBhost, env.DBPort, env.BPort, env.FolderPath, env.JWTSecret)

	// Create or overwrite the .env file
	file, err := os.Create(".env") // This will truncate the file if it already exists
	if err != nil {
		fmt.Println("Error creating .env file:", err)
		return
	}
	defer file.Close()
	// Write the new content into the .env file
	_, err = file.WriteString(s)
	if err != nil {
		fmt.Println("Error writing to .env file:", err)
		return
	}
	fmt.Println(".env file has been successfully replaced with new content!")
}

func logerr(err error) {
	if err != nil {
		log.Fatalf("Error while taking input: %v", err)
	}
}
