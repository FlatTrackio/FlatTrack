/*
	common function calls
*/

package common

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.com/flattrack/flattrack/src/backend/system"
	"gitlab.com/flattrack/flattrack/src/backend/types"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
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

func RegexMatchName(name string) bool {
	matches, _ := regexp.MatchString(`^([ \\u00c0-\\u01ffa-zA-Z'\-])+$`, name)
	return matches
}

func RegexMatchEmail(email string) bool {
	matches, _ := regexp.MatchString(`^\w+@[a-zA-Z_]+?\.[a-zA-Z]{2,3}$`, email)
	return matches
}

func RegexMatchPassword(password string) bool {
	// Minimum eight characters, at least one letter and one number
	matches, _ := regexp.MatchString(`^([a-z]*)([A-Z]*).{10,}$`, password)
	return matches
}

func RegexMatchPhoneNumber(phoneNumber string) bool {
	matches, _ := regexp.MatchString(`^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]\d{3}[\s.-]\d{4}$`, phoneNumber)
	return matches
}

func HashSHA512(input string) (output string) {
	hasher := sha512.New()
	hasher.Write([]byte(input))
	sha512_hash := hex.EncodeToString(hasher.Sum(nil))
	return sha512_hash
}

func GenerateJWTauthToken(db *sql.DB, email string) (tokenString string, err error) {
	secret, err := system.GetJWTsecret(db)
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(time.Hour * 24 * 5)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, types.JWTclaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	})

	tokenString, err = token.SignedString([]byte(secret))
	return tokenString, err
}

func ValidateJWTauthToken(db *sql.DB, r *http.Request) (valid bool, err error) {
	secret, err := system.GetJWTsecret(db)
	if err != nil {
		return false, err
	}
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		return false, nil
	}
	tokenHeaderJWT := strings.Split(tokenHeader, " ")[1]
	claims := &types.JWTclaim{}
	token, err := jwt.ParseWithClaims(tokenHeaderJWT, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}

	return token.Valid, nil
}
