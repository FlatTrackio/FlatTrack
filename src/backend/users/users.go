/*
  users
    manage user accounts
*/

package users

import (
	"database/sql"
	"errors"
	"fmt"

	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.com/flattrack/flattrack/src/backend/common"
	"gitlab.com/flattrack/flattrack/src/backend/groups"
	"gitlab.com/flattrack/flattrack/src/backend/system"
	"gitlab.com/flattrack/flattrack/src/backend/types"
)

// CreateUser
// given a UserSpec, create a user
func CreateUser(db *sql.DB, user types.UserSpec) (userInserted types.UserSpec, err error) {
	if len(user.Names) == 0 || len(user.Names) > 60 || user.Names == "" {
		return userInserted, errors.New("Unable to use the provided name, as it is either empty or too long or too short")
	}
	if common.RegexMatchEmail(user.Email) == false || user.Email == "" {
		return userInserted, errors.New("Unable to use the provided email, as it is either empty or not valid")
	}

	for _, groupItem := range user.Groups {
		group, err := groups.GetGroupByName(db, groupItem)
		if err != nil || group.Id == "" {
			return userInserted, errors.New(fmt.Sprintf("Unable to use the provide group '%v' as it is invalid", groupItem))
		}
	}

	if common.RegexMatchPassword(user.Password) == false || user.Password == "" {
		return userInserted, errors.New("Unable to use the provided password, as it is either empty of invalid")
	}
	user.Password = common.HashSHA512(user.Password)
	if user.PhoneNumber != "" && common.RegexMatchPhoneNumber(user.PhoneNumber) == false {
		return userInserted, errors.New("Unable to use the provided phone number")
	}
	localUser, err := GetUserByEmail(db, user.Email, false)
	if localUser.Email == user.Email || err != nil {
		return userInserted, errors.New("User account already exists")
	}

	sqlStatement := `insert into users (names, email, password, phonenumber, birthday, contractAgreement, disabled, registered, taskNotificationFrequency)
                         values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
                         returning *`
	rows, err := db.Query(sqlStatement, user.Names, user.Email, user.Password, user.PhoneNumber, user.Birthday, user.ContractAgreement, user.Disabled, user.Registered, user.TaskNotificationFrequency)
	if err != nil {
		return userInserted, err
	}
	for rows.Next() {
		userInserted, err = UserObjectFromRows(rows)
		if err != nil {
			return userInserted, err
		}
	}

	for _, groupItem := range user.Groups {
		group, err := groups.GetGroupByName(db, groupItem)
		if err != nil || group.Id == "" {
			return userInserted, errors.New(fmt.Sprintf("Unable to use the provide group '%v' as it is invalid", groupItem))
		}
		err = groups.AddUserToGroup(db, userInserted.Id, group.Id)
		if err != nil {
			return userInserted, errors.New("Unable to add user account to the group")
		}
	}
	return userInserted, err
}

// GetAllUsers
// return all users in the database
func GetAllUsers(db *sql.DB, includePassword bool, selectors types.UserSelector) (users []types.UserSpec, err error) {
	sqlStatement := `select * from users`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		found := true
		user, err := UserObjectFromRows(rows)
		if err != nil {
			return users, err
		}
		groupsOfUser, err := groups.GetGroupNamesOfUserById(db, user.Id)
		if err != nil {
			return users, err
		}
		if selectors.Group != "" {
			found, err = groups.CheckUserInGroup(db, user.Id, selectors.Group)
			if err != nil {
				return users, err
			}
		}
		user.Groups = groupsOfUser
		if includePassword == false {
			user.Password = ""
		}
		if found == false {
			continue
		}
		users = append(users, user)
	}
	return users, err
}

// GetUser
// given a UserSpec and an ID or Email, return a user from the database
func GetUser(db *sql.DB, userSelect types.UserSpec, includePassword bool) (user types.UserSpec, err error) {
	if userSelect.Id != "" {
		return GetUserById(db, userSelect.Id, includePassword)
	}
	if userSelect.Email != "" {
		if common.RegexMatchEmail(userSelect.Email) {
			return user, errors.New("Invalid email address")
		}
		return GetUserByEmail(db, userSelect.Email, includePassword)
	}
	if includePassword == false {
		user.Password = ""
	}
	return user, err
}

// UserObjectFromRows
// construct a UserSpec from database rows
func UserObjectFromRows(rows *sql.Rows) (user types.UserSpec, err error) {
	var id string
	var names string
	var email string
	var password string
	var phoneNumber string
	var birthday int
	var contractAgreement bool
	var disabled bool
	var registered bool
	var taskNotificationFrequency int
	var lastLogin string
	var creationTimestamp int
	var modificationTimestamp int
	var deletionTimestamp int
	err = rows.Scan(&id, &names, &email, &password, &phoneNumber, &birthday, &contractAgreement, &disabled, &registered, &taskNotificationFrequency, &lastLogin, &creationTimestamp, &modificationTimestamp, &deletionTimestamp)
	if err != nil {
		fmt.Println(err)
		return user, err
	}
	user = types.UserSpec{
		Id:                        id,
		Names:                     names,
		Email:                     email,
		Password:                  password,
		PhoneNumber:               phoneNumber,
		ContractAgreement:         contractAgreement,
		Disabled:                  disabled,
		Registered:                registered,
		TaskNotificationFrequency: taskNotificationFrequency,
		LastLogin:                 lastLogin,
		CreationTimestamp:         creationTimestamp,
		ModificationTimestamp:     modificationTimestamp,
		DeletionTimestamp:         deletionTimestamp,
	}
	return user, err
}

// GetUserById
// given an id, return a UserSpec
func GetUserById(db *sql.DB, id string, includePassword bool) (user types.UserSpec, err error) {
	sqlStatement := `select * from users where id = $1`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return user, err
	}
	defer rows.Close()
	rows.Next()
	user, err = UserObjectFromRows(rows)
	if err != nil {
		return user, err
	}
	groups, err := groups.GetGroupNamesOfUserById(db, user.Id)
	if err != nil {
		return user, err
	}
	user.Groups = groups
	if includePassword == false {
		user.Password = ""
	}
	return user, err
}

// GetUserByEmail
// given a email, return a UserSpec
func GetUserByEmail(db *sql.DB, email string, includePassword bool) (user types.UserSpec, err error) {
	sqlStatement := `select * from users where email = $1`
	rows, err := db.Query(sqlStatement, email)
	if err != nil {
		return user, err
	}
	rows.Next()
	user, _ = UserObjectFromRows(rows)
	groups, err := groups.GetGroupNamesOfUserById(db, user.Id)
	if err != nil {
		return user, err
	}
	user.Groups = groups
	if includePassword == false {
		user.Password = ""
	}
	return user, err
}

// DeleteUserById
// given an id, remove the user account from all the groups and then delete a user account
func DeleteUserById(db *sql.DB, id string) (err error) {
	userGroups, err := groups.GetGroupsOfUserById(db, id)
	if err != nil {
		return err
	}
	for _, groupItem := range userGroups {
		err = groups.RemoveUserFromGroup(db, id, groupItem.Id)
		if err != nil {
			return err
		}
	}
	sqlStatement := `delete from users where id = $1`
	_, err = db.Query(sqlStatement, id)
	return err
}

// CheckUserPassword
// given an email and password, find the user account with the email, return if the password matches
func CheckUserPassword(db *sql.DB, email string, password string) (matches bool, err error) {
	user, err := GetUserByEmail(db, email, true)
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
		return false, errors.New("Unable to find FlatTrack system auth secret. Please contact system administrators or support")
	}
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		return false, errors.New("Unable to find authorization token (header doesn't exist)")
	}
	authorizationHeader := strings.Split(tokenHeader, " ")
	if authorizationHeader[0] != "bearer" || len(authorizationHeader) <= 1 {
		return false, errors.New("Unable to find authorization token (must be as bearer)")
	}
	tokenHeaderJWT := authorizationHeader[1]
	claims := &types.JWTclaim{}
	token, err := jwt.ParseWithClaims(tokenHeaderJWT, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}

	reqClaims := token.Claims.(*types.JWTclaim)
	user, err := GetUserById(db, reqClaims.Id, true)
	if err != nil || user.Id == "" {
		return false, errors.New("Unable to find the user account which the authentication token belongs to")
	}

	return token.Valid, nil
}

// GetIdFromJWT
// return the userId in a JWT from a header in a HTTP request
func GetIdFromJWT(db *sql.DB, r *http.Request) (id string, err error) {
	secret, err := system.GetJWTsecret(db)
	if err != nil {
		return "", errors.New("Unable to find FlatTrack system auth secret. Please contact system administrators or support")
	}
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		return "", errors.New("Unable to find authorization token (header doesn't exist)")
	}
	authorizationHeader := strings.Split(tokenHeader, " ")
	if authorizationHeader[0] != "bearer" || len(authorizationHeader) <= 1 {
		return "", errors.New("Unable to find authorization token (must be as bearer)")
	}
	tokenHeaderJWT := authorizationHeader[1]
	claims := &types.JWTclaim{}
	token, err := jwt.ParseWithClaims(tokenHeaderJWT, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	reqClaims := token.Claims.(*types.JWTclaim)
	user, err := GetUserById(db, reqClaims.Id, true)
	if err != nil || user.Id == "" {
		return "", errors.New("Unable to find the user account which the authentication token belongs to")
	}
	return user.Id, err
}

// GetProfile
// return user from Id in JWT from HTTP request
func GetProfile(db *sql.DB, r *http.Request) (user types.UserSpec, err error) {
	id, err := GetIdFromJWT(db, r)
	if err != nil {
		return user, err
	}

	user, err = GetUserById(db, id, false)
	return user, err
}
