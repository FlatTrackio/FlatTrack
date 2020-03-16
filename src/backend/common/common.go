/*
	common function calls
*/

package common

import (
	"fmt"
	"log"
	"net/http"
	"os"
	// "gitlab.com/flattrack/flattrack/src/backend/types"
)

var (
	APP_BUILD_VERSION      = "0.0.0"
	APP_BUILD_HASH         = "???"
	APP_BUILD_DATE         = "???"
	APP_BUILD_MODE         = "development"
	APP_DB_MIGRATIONS_PATH = "/app/migrations"
)

func GetEnvOrDefault(envName string, defaultValue string) (output string) {
	output = os.Getenv(envName)
	if output == "" {
		output = defaultValue
	}
	return output
}

func GetAppMode() (output string) {
	APP_BUILD_MODE = GetEnvOrDefault("APP_MODE", "development")
	return APP_BUILD_MODE
}

func GetDBdatabase() (output string) {
	return GetEnvOrDefault("APP_DB_DATABASE", "flattrack")
}

func GetDBusername() (output string) {
	return GetEnvOrDefault("APP_DB_USERNAME", "")
}

func GetDBhost() (output string) {
	return GetEnvOrDefault("APP_DB_HOST", "")
}

func GetDBpassword() (output string) {
	return GetEnvOrDefault("APP_DB_PASSWORD", "")
}

func GetMigrationsPath() (output string) {
	envSet := GetEnvOrDefault("APP_DB_MIGRATIONS_PATH", "")
	if envSet != "" {
		return envSet
	}
	if APP_BUILD_MODE == "production" {
		return "/app/migrations"
	}
	pwd, _ := os.Getwd()
	return fmt.Sprintf("%v/migrations", pwd)
}

func GetAppPort() (output string) {
	return GetEnvOrDefault("APP_PORT", ":8080")
}

func GetAppBuildVersion() string {
	return APP_BUILD_VERSION
}

func GetAppBuildHash() string {
	return APP_BUILD_HASH
}

func GetAppBuildDate() string {
	return APP_BUILD_DATE
}

func GetAppBuildMode() string {
	return APP_BUILD_MODE
}

func SetFirstOrSecond(f string, s string) string {
	if f != "" {
		return f
	}
	return s
}

func Logging(next http.Handler) http.Handler {
	// log all requests
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v %v %v %v", r.Method, r.URL, r.Proto, r.Response, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func GetAppDistFolder() string {
	return GetEnvOrDefault("APP_DIST_FOLDER", "./dist")
}

func GetAppVersion() string {
	return APP_BUILD_VERSION
}
