/*
  common
    common function calls used through the other packages
*/

package common

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
)

var (
	APP_BUILD_VERSION      = "0.0.0"
	APP_BUILD_HASH         = "???"
	APP_BUILD_DATE         = "???"
	APP_BUILD_MODE         = "development"
	APP_DB_MIGRATIONS_PATH = "/app/migrations"
)

// GetEnvOrDefault
// given the value of an environment variable, return it's data or if not available a default value
func GetEnvOrDefault(envName string, defaultValue string) (output string) {
	output = os.Getenv(envName)
	if output == "" {
		output = defaultValue
	}
	return output
}

// GetDBdatabase
// return the database's database to use
func GetDBdatabase() (output string) {
	return GetEnvOrDefault("APP_DB_DATABASE", "flattrack")
}

// GetDBusername
// return the database user to use
func GetDBusername() (output string) {
	return GetEnvOrDefault("APP_DB_USERNAME", "")
}

// GetDBhost
// return the database host to use
func GetDBhost() (output string) {
	return GetEnvOrDefault("APP_DB_HOST", "")
}

// GetDBpassword
// return the database password to use
func GetDBpassword() (output string) {
	return GetEnvOrDefault("APP_DB_PASSWORD", "")
}

// GetMigrationsPath
// return the path of the database migrations to use
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

// GetAppPort
// return the port which the app should serve HTTP on
func GetAppPort() (output string) {
	return GetEnvOrDefault("APP_PORT", ":8080")
}

// GetAppBuildVersion
// return the version of the current FlatTrack instance
func GetAppBuildVersion() string {
	return APP_BUILD_VERSION
}

// GetAppBuildHash
// return the commit which the current FlatTrack binary was built from
func GetAppBuildHash() string {
	return APP_BUILD_HASH
}

// GetAppBuildDate
// return the build date of FlatTrack
func GetAppBuildDate() string {
	return APP_BUILD_DATE
}

// GetAppBuildMode
// return the mode that the app is built in
func GetAppBuildMode() string {
	return APP_BUILD_MODE
}

// SetFirstOrSecond
// given first, return it, else return second
func SetFirstOrSecond(first string, second string) string {
	if first != "" {
		return first
	}
	return second
}

// GetAppDistFolder
// return the path to the folder containing the frontend assets
func GetAppDistFolder() string {
	return GetEnvOrDefault("APP_DIST_FOLDER", "./dist")
}

// RegexMatchName
// regex check for valid name string
func RegexMatchName(name string) bool {
	matches, _ := regexp.MatchString(`^([ \\u00c0-\\u01ffa-zA-Z'\-])+$`, name)
	return matches
}

// RegexMatchEmail
// regex check for valid email address string
func RegexMatchEmail(email string) bool {
	matches, _ := regexp.MatchString(`^\w+@[a-zA-Z_]+?\.[a-zA-Z]{2,3}$`, email)
	return matches
}

// RegexMatchPassword
// regex check for valid password
// rules:
// - 10+ characters
// - at least one lowercase character
// - at least one uppercase character
func RegexMatchPassword(password string) bool {
	matches, _ := regexp.MatchString(`^([a-z]*)([A-Z]*).{10,}$`, password)
	return matches
}

// RegexMatchPhoneNumber
// regex check for valid phonenumber
func RegexMatchPhoneNumber(phoneNumber string) bool {
	matches, _ := regexp.MatchString(`^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]\d{3}[\s.-]\d{4}$`, phoneNumber)
	return matches
}

// HashSHA512
// given an input string, return a SHA512 hashed representation of it
func HashSHA512(input string) (output string) {
	hasher := sha512.New()
	hasher.Write([]byte(input))
	sha512_hash := hex.EncodeToString(hasher.Sum(nil))
	return sha512_hash
}
