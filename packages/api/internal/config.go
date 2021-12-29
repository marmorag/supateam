package internal

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Configuration struct {
	// DB
	DbHost string
	DbPort int
	DbName string
	DbUser string
	DbPass string
	// Cors policy
	CorsAllowOrigins string
	CorsAllowHeaders string
	CorsAllowMethods string
	// App configuration
	ApplicationSecret  string
	ApplicationPrefork bool
	ApplicationName    string
}

var config Configuration

func GetConfig() Configuration {
	if config == (Configuration{}) {
		InitConfig()
	}

	return config
}

func InitConfig() {
	_ = godotenv.Load(fmt.Sprintf(".env%s", getEnvExtension()))

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "27017"))

	config = Configuration{
		DbHost:             getEnv("DB_HOST", "localhost"),
		DbPort:             dbPort,
		DbName:             getEnv("DB_NAME", "supateam"),
		DbUser:             getEnv("DB_USER", "supateam"),
		DbPass:             getEnv("DB_PASS", "supateam"),
		CorsAllowOrigins:   getEnv("CORS_ALLOW_ORIGINS", "*"),
		CorsAllowHeaders:   getEnv("CORS_ALLOW_HEADERS", ""),
		CorsAllowMethods:   getEnv("CORS_ALLOW_METHODS", "GET,POST,HEAD,PUT,DELETE,OPTIONS"),
		ApplicationSecret:  getEnv("APP_SECRET", ""),
		ApplicationPrefork: getEnvBool("APP_PREFORK", false),
		ApplicationName:    getEnv("APP_NAME", "SupaTeam"),
	}
}

func getEnvExtension() string {
	extension := ""
	if env := os.Getenv("APP_ENV"); env == "prod" {
		extension = ".prod"
	}

	return extension
}

func getEnv(key string, defaultValue string) string {
	v, set := os.LookupEnv(key)

	if set {
		return v
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	v, set := os.LookupEnv(key)

	if set {
		return v == "true"
	}
	return defaultValue
}
