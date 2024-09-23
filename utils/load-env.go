package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from the .env file.
//
// Parameters:
// - level: int, the number of levels of directories to search for the .env file.
// No return type.
func LoadEnv(level int) {
	envs := findEnvFiles(level)

	if len(envs) == 0 {
		log.Fatal("Could not resolve .env file")
	}

	err := godotenv.Overload(envs...)
	if err != nil {
		log.Println("Warning: error loading .env file", err)
	}
}

// findEnvFiles finds .env files recursively up to the specified level.
func findEnvFiles(level int) []string {
	var envs []string

	// Iterate through directories up to the specified level
	for i := level; i >= 0; i-- {
		// Construct the path to the .env file
		envPath := filepath.Join(repeat("../", i), ".env")
		// Check if the .env file exists
		if _, err := os.Stat(envPath); err == nil {
			envs = append(envs, envPath)
		}
	}

	return envs
}

// repeat repeats a string 's' 'count' times.
func repeat(s string, count int) string {
	returnString := ""
	for i := 0; i < count; i++ {
		returnString += s
	}
	return returnString
}
