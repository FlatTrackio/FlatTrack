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

	"github.com/imdario/mergo"

	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.com/flattrack/flattrack/src/backend/common"
	"gitlab.com/flattrack/flattrack/src/backend/groups"
	"gitlab.com/flattrack/flattrack/src/backend/system"
	"gitlab.com/flattrack/flattrack/src/backend/types"
)

// ValidateUser
// given a UserSpec, return if it's valid
func ValidateUser(db *sql.DB, user types.UserSpec, allowEmptyPassword bool) (valid bool, err error) {
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

	if user.Birthday != 0 && (common.ValidateBirthday(int64(user.Birthday)*1000) == false) {
		return false, errors.New("Unable to use the provided birthday, your birthday year must not be within the last 15 years")
	}

	if (common.RegexMatchPassword(user.Password) == false || user.Password == "") && allowEmptyPassword == false {
		return false, errors.New("Unable to use the provided password, as it is either empty of invalid")
	}
	if user.PhoneNumber != "" && common.RegexMatchPhoneNumber(user.PhoneNumber) == false {
		return false, errors.New("Unable to use the provided phone number")
	}

	return true, err
}

// CreateUser
// given a UserSpec, create a user
func CreateUser(db *sql.DB, user types.UserSpec, allowEmptyPassword bool) (userInserted types.UserSpec, err error) {
	var userCreationSecretInserted types.UserCreationSecretSpec

	validUser, err := ValidateUser(db, user, allowEmptyPassword)
	if !validUser || err != nil {
		return userInserted, err
	}
	localUser, err := GetUserByEmail(db, user.Email, false)
	if err == nil || localUser.Id != "" {
		return userInserted, err
	}
	if localUser.Email == user.Email {
		return userInserted, errors.New("Email address is already taken")
	}
	if user.Password != "" {
		user.Password = common.HashSHA512(user.Password)
	}
	if allowEmptyPassword == false {
		user.Registered = true
	}

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
	if allowEmptyPassword == true {
		userCreationSecretInserted, err = CreateUserCreationSecret(db, userInserted.Id)
		if err != nil {
			return userInserted, err
		}
		if userCreationSecretInserted.Id == "" {
			return userInserted, errors.New("Failed to created a user creation secret")
		}
	}

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
		} else if selectors.Id != "" {
			found = selectors.Id == user.Id
		} else if selectors.NotId != "" {
			found = !(selectors.NotId == user.Id)
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
	if user.Id == "" {
		return user, errors.New("Failed to find user")
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

	err = DeleteUserCreationSecretByUserId(db, id)
	if err != nil {
		return err
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
func GenerateJWTauthToken(db *sql.DB, id string, authNonce string, expiresIn time.Duration) (tokenString string, err error) {
	if expiresIn == 0 {
		expiresIn = 24 * 5
	}
	secret, err := system.GetJWTsecret(db)
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(time.Hour * expiresIn)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, types.JWTclaim{
		Id:        id,
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

// PatchProfile
// patches the profile of a user account
func PatchProfile(db *sql.DB, id string, userAccount types.UserSpec) (userAccountPatched types.UserSpec, err error) {
	// TODO ensure that an account can't change their address to someone else's
	existingUserAccount, err := GetUserById(db, id, true)
	if err != nil || existingUserAccount.Id == "" {
		return userAccountPatched, errors.New("Failed to find user account")
	}
	err = mergo.Merge(&userAccount, existingUserAccount)
	if err != nil {
		return userAccountPatched, errors.New("Failed to update fields in the user account")
	}
	noUpdatePassword := userAccount.Password == existingUserAccount.Password
	valid, err := ValidateUser(db, userAccount, noUpdatePassword)
	if !valid || err != nil {
		return existingUserAccount, err
	}
	passwordHashed := common.HashSHA512(userAccount.Password)
	if noUpdatePassword == true {
		passwordHashed = userAccount.Password
	}

	sqlStatement := `update users set names = $1, email = $2, password = $3, phoneNumber = $4, birthday = $5, contractAgreement = $6, disabled = $7, registered = $8, lastLogin = $9, authNonce = $10, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $11
                         returning *`
	rows, err := db.Query(sqlStatement, userAccount.Names, userAccount.Email, passwordHashed, userAccount.PhoneNumber, userAccount.Birthday, userAccount.ContractAgreement, userAccount.Disabled, userAccount.Registered, userAccount.LastLogin, userAccount.AuthNonce, id)
	if err != nil {
		// TODO add roll back, if there's failure
		return userAccountPatched, err
	}
	rows.Next()
	userAccountPatched, err = UserObjectFromRows(rows)
	if err != nil || userAccountPatched.Id == "" {
		return userAccountPatched, errors.New("Failed to create shopping list")
	}

	updatedGroups, err := groups.UpdateUserGroups(db, id, userAccount.Groups)
	if updatedGroups == false || err != nil {
		return existingUserAccount, err
	}

	userAccountPatched.Groups = userAccount.Groups
	userAccountPatched.Password = ""
	return userAccountPatched, err
}

// UserCreationSecretsFromRows
// constructs a UserCreationSecretSpec from rows
func UserCreationSecretsFromRows(rows *sql.Rows) (creationSecret types.UserCreationSecretSpec, err error) {
	rows.Scan(&creationSecret.Id, &creationSecret.UserId, &creationSecret.Secret, &creationSecret.Valid, creationSecret.CreationTimestamp, &creationSecret.ModificationTimestamp, &creationSecret.DeletionTimestamp)
	err = rows.Err()
	return creationSecret, err
}

// GetAllUserCreationSecrets
// returns all UserCreationSecrets from the database
func GetAllUserCreationSecrets(db *sql.DB, secretsSelector types.UserCreationSecretSelector) (creationSecrets []types.UserCreationSecretSpec, err error) {
	sqlStatement := `select * from user_creation_secret`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return creationSecrets, err
	}
	for rows.Next() {
		creationSecret, err := UserCreationSecretsFromRows(rows)
		if err != nil {
			return creationSecrets, errors.New("Failed to list user creation secrets")
		}
		if secretsSelector.UserId != "" {
			userExists, err := UserAccountExists(db, secretsSelector.UserId)
			if err != nil {
				return creationSecrets, err
			}
			if userExists == false {
				return creationSecrets, errors.New("Unable to find user account")
			}
			if creationSecret.UserId != secretsSelector.UserId {
				continue
			}
		}
		creationSecrets = append(creationSecrets, creationSecret)
	}
	return creationSecrets, err
}

// GetUserCreationSecret
// returns a UserCreationSecret by it's id from the database
func GetUserCreationSecret(db *sql.DB, id string) (creationSecret types.UserCreationSecretSpec, err error) {
	sqlStatement := `select * from user_creation_secret where id = $1`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return creationSecret, err
	}
	defer rows.Close()
	rows.Next()
	creationSecret, err = UserCreationSecretsFromRows(rows)
	return creationSecret, err
}

// CreateUserCreationSecret
// creates a user creation secret for account confirming
func CreateUserCreationSecret(db *sql.DB, userId string) (userCreationSecretInserted types.UserCreationSecretSpec, err error) {
	sqlStatement := `insert into user_creation_secret (userId)
                         values ($1)
                         returning *`
	rows, err := db.Query(sqlStatement, userId)
	if err != nil {
		return userCreationSecretInserted, err
	}
	defer rows.Close()
	rows.Next()
	userCreationSecretInserted, err = UserCreationSecretsFromRows(rows)
	return userCreationSecretInserted, err
}

// DeleteUserCreationSecret
// deletes the acccount creation secret, after it's been used
func DeleteUserCreationSecret(db *sql.DB, id string) (err error) {
	sqlStatement := `delete from user_creation_secret where id = $1`
	_, err = db.Query(sqlStatement, id)
	return err
}

// DeleteUserCreationSecretByUserId
// deletes the acccount creation secret by userid, after it's been used
func DeleteUserCreationSecretByUserId(db *sql.DB, userId string) (err error) {
	sqlStatement := `delete from user_creation_secret where userId = $1`
	_, err = db.Query(sqlStatement, userId)
	return err
}

// ConfirmUserAccount
// confirms the user account
func ConfirmUserAccount(db *sql.DB, id string, secret string, user types.UserSpec) (tokenString string, err error) {
	if user.Password == "" {
		return tokenString, errors.New("Unable to confirm account, a password must be provided to complete registration")
	}
	userCreationSecret, err := GetUserCreationSecret(db, id)
	if err != nil {
		return tokenString, err
	}
	if userCreationSecret.Id == "" {
		return tokenString, errors.New("Unable to find account confirmation secret")
	}
	if secret != userCreationSecret.Secret {
		return tokenString, errors.New("Unable to confirm account, as the secret doesn't match")
	}
	userInDB, err := GetUserById(db, userCreationSecret.UserId, false)
	if err != nil {
		return tokenString, err
	}

	userAccountPatch := types.UserSpec{
		Names:       user.Names,
		Email:       user.Email,
		Password:    user.Password,
		Birthday:    user.Birthday,
		PhoneNumber: user.PhoneNumber,
		Registered:  true,
	}
	userConfirmed, err := PatchProfile(db, userCreationSecret.UserId, userAccountPatch)
	if err != nil {
		return tokenString, err
	}
	if userConfirmed.Id == "" || userConfirmed.Registered == false {
		return tokenString, errors.New("Failed to patch profile")
	}
	err = DeleteUserCreationSecret(db, id)

	return GenerateJWTauthToken(db, userInDB.Id, userInDB.AuthNonce, 0)
}

// UserAccountExists
// returns bool if user account exists
func UserAccountExists(db *sql.DB, id string) (exists bool, err error) {
	sqlStatement := `select id from users where id = $1`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	rows.Next()
	var userIdFromDB string
	rows.Scan(&userIdFromDB)
	err = rows.Err()
	return userIdFromDB == id, err
}
