/*
  common
    common function calls used through the other packages
*/

// This program is free software: you can redistribute it and/or modify
// it under the terms of the Affero GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the Affero GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package common

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"path"
	"regexp"
	"time"
)

// AppVars ...
// defaults which are overridden with build
var (
	AppBuildVersion     = "0.0.0"
	AppBuildHash        = "???"
	AppBuildDate        = "???"
	AppBuildMode        = "development"
	AppDbMigrationsPath = "/var/run/ko/migrations"
	// #nosec G101
	AppAssetsFolder = "/var/run/ko/web"
)

// GetEnvOrDefault ...
// given the value of an environment variable, return it's data or if not available a default value
func GetEnvOrDefault(envName string, defaultValue string) (output string) {
	output = os.Getenv(envName)
	if output == "" {
		output = defaultValue
	}
	return output
}

func GetDBConnectionString() (output string) {
	return os.ExpandEnv(GetEnvOrDefault("APP_DB_CONNECTION_STRING", ""))
}

// GetDBdatabase ...
// return the database's database to use
func GetDBdatabase() (output string) {
	return GetEnvOrDefault("APP_DB_DATABASE", "flattrack")
}

// GetDBusername ...
// return the database user to use
func GetDBusername() (output string) {
	return GetEnvOrDefault("APP_DB_USERNAME", "flattrack")
}

// GetDBhost ...
// return the database host to use
func GetDBhost() (output string) {
	return GetEnvOrDefault("APP_DB_HOST", "localhost")
}

// GetDBport ...
// return the database port to use
func GetDBport() (output string) {
	return GetEnvOrDefault("APP_DB_PORT", "5432")
}

// GetDBpassword ...
// return the database password to use
func GetDBpassword() (output string) {
	return GetEnvOrDefault("APP_DB_PASSWORD", "flattrack")
}

// GetDBsslMode ...
// return the database sslmode to use
func GetDBsslMode() (output string) {
	return GetEnvOrDefault("APP_DB_SSLMODE", "disable")
}

// GetInstanceURL ...
// return URL of the instance
func GetInstanceURL() (output string) {
	return GetEnvOrDefault("APP_URL", "")
}

// GetSMTPEnabled ...
// return if the instance should send emails
func GetSMTPEnabled() (output string) {
	return GetEnvOrDefault("APP_SMTP_ENABLED", "false")
}

// GetSMTPUsername ...
// return the username to send emails with
func GetSMTPUsername() (output string) {
	return GetEnvOrDefault("APP_SMTP_USERNAME", "")
}

// GetSMTPPassword ...
// return the password to send emails with
func GetSMTPPassword() (output string) {
	return GetEnvOrDefault("APP_SMTP_PASSWORD", "")
}

// GetSMTPHost ...
// return the host to send emails with
func GetSMTPHost() (output string) {
	return GetEnvOrDefault("APP_SMTP_HOST", "")
}

// GetSMTPPort ...
// return the port to send emails with
func GetSMTPPort() (output string) {
	return GetEnvOrDefault("APP_SMTP_PORT", "")
}

// GetMigrationsPath ...
// return the path of the database migrations to use
func GetMigrationsPath() (output string) {
	if envSet := GetEnvOrDefault("APP_DB_MIGRATIONS_PATH", ""); envSet != "" {
		return envSet
	}
	if AppBuildMode == "production" || AppBuildMode == "staging" {
		return AppDbMigrationsPath
	}
	pwd, _ := os.Getwd()
	return fmt.Sprintf("%v/kodata/migrations", pwd)
}

// GetAppPort ...
// return the port which the app should serve HTTP on
func GetAppPort() (output string) {
	return GetEnvOrDefault("APP_PORT", ":8080")
}

// GetAppMetricsPort ...
// return the port which the app should serve metrics on
func GetAppMetricsPort() (output string) {
	return GetEnvOrDefault("APP_PORT_METRICS", ":2112")
}

// GetAppMetricsEnabled ...
// serve metrics endpoint
func GetAppMetricsEnabled() (output bool) {
	return GetEnvOrDefault("APP_METRICS_ENABLED", "true") == "true"
}

// GetAppHealthPort ...
// return the port which the app should serve health on
func GetAppHealthPort() (output string) {
	return GetEnvOrDefault("APP_PORT_HEALTH", ":8081")
}

// GetAppHealthEnabled ...
// serve health endpoint
func GetAppHealthEnabled() (output string) {
	return GetEnvOrDefault("APP_HEALTH_ENABLED", "true")
}

// GetAppEnvFile ...
// location of an env file to load
func GetAppEnvFile() (output string) {
	return GetEnvOrDefault("APP_ENV_FILE", ".env")
}

// GetAppRealIPHeader ...
// the header to use instead of r.RemoteAddr
func GetAppRealIPHeader() (output string) {
	return GetEnvOrDefault("APP_HTTP_REAL_IP_HEADER", "")
}

// GetAppSetupMessage ...
// return a message to display on setup
func GetAppSetupMessage() (output string) {
	return GetEnvOrDefault("APP_SETUP_MESSAGE", "")
}

// GetAppLoginMessage ...
// return a message to display on login
func GetAppLoginMessage() (output string) {
	return GetEnvOrDefault("APP_LOGIN_MESSAGE", "")
}

// GetAppEmbeddedHTML ...
// return HTML to inject into index.html
func GetAppEmbeddedHTML() (output string) {
	return GetEnvOrDefault("APP_EMBEDDED_HTML", "")
}

// GetAppMinioAccessKey ...
// return the accessKey for file storage
func GetAppMinioAccessKey() (output string) {
	return GetEnvOrDefault("APP_MINIO_ACCESS_KEY", "")
}

// GetAppMinioSecretKey ...
// return the secretKey for file storage
func GetAppMinioSecretKey() (output string) {
	return GetEnvOrDefault("APP_MINIO_SECRET_KEY", "")
}

// GetAppMinioBucket ...
// return the bucket for file storage
func GetAppMinioBucket() (output string) {
	return GetEnvOrDefault("APP_MINIO_BUCKET", "")
}

// GetAppMinioHost ...
// return the host for file storage
func GetAppMinioHost() (output string) {
	return GetEnvOrDefault("APP_MINIO_HOST", "0.0.0.0:9000")
}

// GetAppMinioUseSSL ...
// return if the file storage should use SSL
func GetAppMinioUseSSL() (output string) {
	return GetEnvOrDefault("APP_MINIO_USE_SSL", "")
}

// GetSchedulerDisableUseEndpoint return whether to disable the scheduler and to use and endpoint instead
func GetSchedulerDisableUseEndpoint() bool {
	return GetEnvOrDefault("APP_SCHEDULER_DISABLE_USE_ENDPOINT", "false") == "true"
}

// GetSchedulerEndpointSecret the shared secret to require for posting to scheduler endpoint
func GetSchedulerEndpointSecret() string {
	return GetEnvOrDefault("APP_SCHEDULER_ENDPOINT_SECRET", "")
}

// GetRegistrationSecret the shared secret to require for setting up an instance
func GetRegistrationSecret() string {
	return GetEnvOrDefault("APP_REGISTRATION_SECRET", "")
}

// GetAppBuildVersion ...
// return the version of the current FlatTrack instance
func GetAppBuildVersion() string {
	return AppBuildVersion
}

// GetAppBuildHash ...
// return the commit which the current FlatTrack binary was built from
func GetAppBuildHash() string {
	return AppBuildHash
}

// GetAppBuildDate ...
// return the build date of FlatTrack
func GetAppBuildDate() string {
	return AppBuildDate
}

// GetAppBuildMode ...
// return the mode that the app is built in
func GetAppBuildMode() string {
	return AppBuildMode
}

// SetFirstOrSecond ...
// given first, return it, else return second
func SetFirstOrSecond(first string, second string) string {
	if first != "" {
		return first
	}
	return second
}

// GetAppDistFolder ...
// return the path to the folder containing the frontend assets
func GetAppDistFolder() string {
	if AppBuildMode == "production" || AppBuildMode == "staging" {
		return AppAssetsFolder
	}
	if envSet := GetEnvOrDefault("APP_DIST_FOLDER", ""); envSet != "" {
		return envSet
	}
	pwd, _ := os.Getwd()
	return path.Join(pwd, "web", "dist")
}

// RegexMatchName ...
// regex check for valid name string
func RegexMatchName(name string) bool {
	matches, _ := regexp.MatchString(`^([ \\u00c0-\\u01ffa-zA-Z'\-])+$`, name)
	return matches
}

// RegexMatchEmail ...
// regex check for valid email address string
// must also be <= 70
func RegexMatchEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email) && len(email) <= 70
}

// RegexMatchPassword ...
// regex check for valid password
// rules:
// - 10 - 70 characters
// - at least one lowercase character
// - at least one uppercase character
func RegexMatchPassword(password string) bool {
	matches, _ := regexp.MatchString(`^([a-zA-Z]*).{10,}$`, password)
	return matches && len(password) <= 70
}

// RegexMatchPhoneNumber ...
// regex check for valid phonenumber
func RegexMatchPhoneNumber(phoneNumber string) bool {
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	return re.MatchString(phoneNumber)
}

// ValidateBirthday ...
// return where a birthday timestamp is valid
// validation requirements is between 100 and 15 years ago
func ValidateBirthday(timestamp int64) bool {
	dateNow := time.Now()
	timestampParsed := time.Unix(timestamp, 0)
	above15yearsAgo := dateNow.Year()-timestampParsed.Year() >= 15
	below100yearsAgo := dateNow.Year()-timestampParsed.Year() <= 100
	return above15yearsAgo && below100yearsAgo
}

// HashSHA512 ...
// given an input string, return a SHA512 hashed representation of it
func HashSHA512(input string) (output string) {
	hasher := sha512.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// StringInStringSlice ...
// given a list of string and an input string, return if the input string is in the list of strings
func StringInStringSlice(input string, list []string) bool {
	for _, item := range list {
		if item == input {
			return true
		}
	}
	return false
}

// RandStringRunes generates a random string with length of n using runes
// nolint:gosec
func RandStringRunes(n int) string {
	letterRunes := []rune("bcdfghjklmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ012358")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
