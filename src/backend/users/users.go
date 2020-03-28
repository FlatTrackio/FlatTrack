/*
  users
    manage user accounts
*/

package users

import (
	"database/sql"
	"errors"

	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.com/flattrack/flattrack/src/backend/common"
	"gitlab.com/flattrack/flattrack/src/backend/system"
	"gitlab.com/flattrack/flattrack/src/backend/types"
)

// CreateUser
// given a UserSpec, create a user
func CreateUser(db *sql.DB, user types.UserSpec) (userInserted types.UserSpec, err error) {
	if common.RegexMatchName(user.Names) == false || user.Names == "" {
		return userInserted, errors.New("Unable to use the provided name, as it is either empty or not valid")
	}
	if common.RegexMatchEmail(user.Email) == false || user.Email == "" {
		return userInserted, errors.New("Unable to use the provided email")
	}

	// TODO add group validation - requires creating admin and flatmember in migrations
	if common.RegexMatchPassword(user.Password) == false || user.Password == "" {
		return userInserted, errors.New("Unable to use the provided password")
	}
	user.Password = common.HashSHA512(user.Password)
	if user.PhoneNumber != "" && common.RegexMatchPhoneNumber(user.PhoneNumber) == false {
		return userInserted, errors.New("Unable to use the provided phone number")
	}
	localUser, err := GetUserByEmail(db, user.Email)
	if localUser.Email == user.Email || err != nil {
		return userInserted, errors.New("User account already exists")
	}

	sqlStatement := `insert into users (names, email, password, phonenumber)
                         values ($1, $2, $3, $4)
                         returning id`
	rows, err := db.Query(sqlStatement, user.Names, user.Email, user.Password, user.PhoneNumber)
	if err != nil {
		return userInserted, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		rows.Scan(&id)
		userInserted = types.UserSpec{
			Id:          id,
			Names:       user.Names,
			Email:       user.Email,
			Groups:      user.Groups,
			PhoneNumber: user.PhoneNumber,
		}
	}
	return userInserted, err
}

// GetAllUsers
// return all users in the database
func GetAllUsers(db *sql.DB) (users []types.UserSpec, err error) {
	sqlStatement := `select * from users`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var names string
		var email string
		var password string
		var phoneNumber string
		var contractAgreement bool
		var disabled bool
		var hasSetPassword bool
		var taskNotificationFrequency int
		var lastLogin string
		var creationTimestamp int64
		var modificationTimestamp int64
		var deletionTimestamp int64
		rows.Scan(&id, &names, &email, &password, &phoneNumber, &contractAgreement, &disabled, &hasSetPassword, &taskNotificationFrequency, &lastLogin, &creationTimestamp, &modificationTimestamp, &deletionTimestamp)
		users = append(users, types.UserSpec{
			Id:                        id,
			Names:                     names,
			Email:                     email,
			Password:                  password,
			PhoneNumber:               phoneNumber,
			ContractAgreement:         contractAgreement,
			Disabled:                  disabled,
			HasSetPassword:            hasSetPassword,
			TaskNotificationFrequency: taskNotificationFrequency,
			LastLogin:                 lastLogin,
			CreationTimestamp:         creationTimestamp,
			ModificationTimestamp:     modificationTimestamp,
			DeletionTimestamp:         deletionTimestamp,
		})
	}
	return users, err
}

// GetUser
// given a UserSpec and an ID or Email, return a user from the database
func GetUser(db *sql.DB, userSelect types.UserSpec) (user types.UserSpec, err error) {
	if userSelect.Id != "" {
		return GetUserById(db, userSelect.Id)
	}
	if userSelect.Email != "" {
		if common.RegexMatchEmail(userSelect.Email) {
			return user, errors.New("Invalid email address")
		}
		return GetUserByEmail(db, userSelect.Email)
	}
	return user, err
}

// UserObjectFromRows
// construct a UserSpec from database rows
func UserObjectFromRows(rows *sql.Rows) (user types.UserSpec, err error) {
	defer rows.Close()
	for rows.Next() {
		var id string
		var names string
		var email string
		var password string
		var phoneNumber string
		var contractAgreement bool
		var disabled bool
		var hasSetPassword bool
		var taskNotificationFrequency int
		var lastLogin string
		var creationTimestamp int64
		var modificationTimestamp int64
		var deletionTimestamp int64
		rows.Scan(&id, &names, &email, &password, &phoneNumber, &contractAgreement, &disabled, &hasSetPassword, &taskNotificationFrequency, &lastLogin, &creationTimestamp, &modificationTimestamp, &deletionTimestamp)
		user = types.UserSpec{
			Id:                        id,
			Names:                     names,
			Email:                     email,
			Password:                  password,
			PhoneNumber:               phoneNumber,
			ContractAgreement:         contractAgreement,
			Disabled:                  disabled,
			HasSetPassword:            hasSetPassword,
			TaskNotificationFrequency: taskNotificationFrequency,
			LastLogin:                 lastLogin,
			CreationTimestamp:         creationTimestamp,
			ModificationTimestamp:     modificationTimestamp,
			DeletionTimestamp:         deletionTimestamp,
		}
	}
	return user, err
}

// GetUserById
// given an id, return a UserSpec
func GetUserById(db *sql.DB, id string) (user types.UserSpec, err error) {
	sqlStatement := `select * from users where id = $1`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return user, err
	}
	user, err = UserObjectFromRows(rows)
	return user, err
}

// GetUserByEmail
// given a email, return a UserSpec
func GetUserByEmail(db *sql.DB, email string) (user types.UserSpec, err error) {
	sqlStatement := `select * from users where email = $1`
	rows, err := db.Query(sqlStatement, email)
	if err != nil {
		return user, err
	}
	defer rows.Close()
	user, err = UserObjectFromRows(rows)
	return user, err
}

// DeleteUserById
// given an id, delete a user account
func DeleteUserById(db *sql.DB, id string) (err error) {
	sqlStatement := `delete from users where id = $1`
	_, err = db.Query(sqlStatement, id)
	return err
}

// CheckUserPassword
// given an email and password, find the user account with the email, return if the password matches
func CheckUserPassword(db *sql.DB, email string, password string) (matches bool, err error) {
	user, err := GetUserByEmail(db, email)
	if err != nil {
		return matches, err
	}
	passwordHashed := common.HashSHA512(password)
	return user.Password == passwordHashed, err
}

// GenerateJWTauthToken
// given an email, return a usable JWT token
func GenerateJWTauthToken(db *sql.DB, id string) (tokenString string, err error) {
	secret, err := system.GetJWTsecret(db)
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(time.Hour * 24 * 5)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, types.JWTclaim{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	})

	tokenString, err = token.SignedString([]byte(secret))
	return tokenString, err
}

// ValidateJWTauthToken
// given an HTTP request and Authorization header, return if auth is valid
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

	reqClaims := token.Claims.(*types.JWTclaim)
	user, err := GetUserById(db, reqClaims.Id)
	if err != nil || user.Id == "" {
		return false, err
	}

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}
