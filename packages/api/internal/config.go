package internal

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
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
	ApplicationSecret        string
	ApplicationEnvironment   string
	ApplicationPrefork       bool
	ApplicationName          string
	ApplicationAESPassphrase string
	RequestIDKey             string
	// Tracing config
	TracingEnabled bool
}

var config Configuration

func GetConfig() Configuration {
	if config == (Configuration{}) {
		InitConfig()
	}

	return config
}

func InitConfig() {
	if envPath := os.Getenv("APP_ENV_PATH"); envPath != "" {
		err := godotenv.Load(envPath)
		if err != nil {
			log.Println("error while loading environment :", err)
		}
	} else {
		err := godotenv.Load(fmt.Sprintf(".env%s", getEnvExtension()))
		if err != nil {
			log.Println("error while loading environment :", err)
		}
	}

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "27017"))

	config = Configuration{
		DbHost:                   getEnv("DB_HOST", "localhost"),
		DbPort:                   dbPort,
		DbName:                   getEnv("DB_NAME", "supateam"),
		DbUser:                   getEnv("DB_USER", "supateam"),
		DbPass:                   getEnv("DB_PASS", "supateam"),
		CorsAllowOrigins:         getEnv("CORS_ALLOW_ORIGINS", "*"),
		CorsAllowHeaders:         getEnv("CORS_ALLOW_HEADERS", ""),
		CorsAllowMethods:         getEnv("CORS_ALLOW_METHODS", "GET,POST,HEAD,PUT,DELETE,OPTIONS"),
		ApplicationSecret:        getEnv("APP_SECRET", ""),
		ApplicationEnvironment:   getEnv("APP_ENV", "dev"),
		ApplicationPrefork:       getEnvBool("APP_PREFORK", false),
		ApplicationName:          getEnv("APP_NAME", "SupaTeam"),
		ApplicationAESPassphrase: getEnv("APP_AES_PASSPHRASE", ""),
		RequestIDKey:             "requestid",
		TracingEnabled:           getEnvBool("TRACING_ENABLED", false),
	}
}

func Set(newCfg Configuration) {
	config = newCfg
}

func Clear() {
	config = Configuration{}
}

func getEnvExtension() string {
	extension := ""
	if env := os.Getenv("APP_ENV"); env == "prod" {
		extension = ".prod"
	}

	if env := os.Getenv("APP_ENV"); env == "test" {
		extension = ".test"
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
