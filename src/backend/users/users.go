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

// ValidateUser
// given a UserSpec, return if it's valid
func ValidateUser(db *sql.DB, user types.UserSpec) (valid bool, err error) {
	if len(user.Names) == 0 || len(user.Names) > 60 || user.Names == "" {
		return false, errors.New("Unable to use the provided name, as it is either empty or too long or too short")
	}
	if common.RegexMatchEmail(user.Email) == false || user.Email == "" {
		return false, errors.New("Unable to use the provided email, as it is either empty or not valid")
	}

	if len(user.Groups) == 0 {
		return false, errors.New("No groups provided; please select at least one group")
	}
	for _, groupItem := range user.Groups {
		group, err := groups.GetGroupByName(db, groupItem)
		if err != nil || group.Id == "" {
			return false, errors.New(fmt.Sprintf("Unable to use the provide group '%v' as it is invalid", groupItem))
		}
	}

	if common.RegexMatchPassword(user.Password) == false || user.Password == "" {
		return false, errors.New("Unable to use the provided password, as it is either empty of invalid")
	}
	if user.PhoneNumber != "" && common.RegexMatchPhoneNumber(user.PhoneNumber) == false {
		return false, errors.New("Unable to use the provided phone number")
	}
	localUser, err := GetUserByEmail(db, user.Email, false)
	if localUser.Email == user.Email || err != nil {
		return false, errors.New("User account already exists")
	}

	return true, err
}

// CreateUser
// given a UserSpec, create a user
func CreateUser(db *sql.DB, user types.UserSpec) (userInserted types.UserSpec, err error) {
	validUser, err := ValidateUser(db, user)
	if !validUser || err != nil {
		return userInserted, err
	}
	user.Password = common.HashSHA512(user.Password)

	sqlStatement := `insert into users (names, email, password, phonenumber, birthday, contractAgreement, disabled, registered)
                         values ($1, $2, $3, $4, $5, $6, $7, $8)
                         returning *`
	rows, err := db.Query(sqlStatement, user.Names, user.Email, user.Password, user.PhoneNumber, user.Birthday, user.ContractAgreement, user.Disabled, user.Registered)
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
	userInserted.Groups = user.Groups
	userInserted.Password = ""
	return userInserted, err
}

// GetAllUsers
// return all users in the database
func GetAllUsers(db *sql.DB, includePassword bool, selectors types.UserSelector) (users []types.UserSpec, err error) {
	sqlStatement := `select * from users where deletionTimestamp = 0`
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
	rows.Scan(&user.Id, &user.Names, &user.Email, &user.Password, &user.PhoneNumber, &user.Birthday, &user.ContractAgreement, &user.Disabled, &user.Registered, &user.LastLogin, &user.AuthNonce, &user.CreationTimestamp, &user.ModificationTimestamp, &user.DeletionTimestamp)
	err = rows.Err()
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
	sqlStatement := `update users set names = '(Deleted User)', email = '', password = '', deletionTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1`
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
func GenerateJWTauthToken(db *sql.DB, id string, authNonce string) (tokenString string, err error) {
	secret, err := system.GetJWTsecret(db)
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(time.Hour * 24 * 5)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, types.JWTclaim{
		Id: id,
		AuthNonce: authNonce,
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

	if reqClaims.AuthNonce != user.AuthNonce {
		return false, errors.New("Authentication has been invalidated, please log in again")
	}

	return token.Valid, nil
}

// InvalidAllAuthTokens
// updates the authNonce to invalidate auth tokens
func InvalidateAllAuthTokens(db *sql.DB, id string) (err error) {
	sqlStatement := `update users set authNonce = md5(random()::text || clock_timestamp()::text)::uuid where id = $1`
	_, err = db.Query(sqlStatement, id)
	return err
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
