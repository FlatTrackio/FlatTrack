/*
  users
    manage user accounts
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

package users

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/imdario/mergo"

	jwt "github.com/golang-jwt/jwt"
	"gitlab.com/flattrack/flattrack/pkg/common"
	"gitlab.com/flattrack/flattrack/pkg/groups"
	"gitlab.com/flattrack/flattrack/pkg/system"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

var (
	jwtAlg *jwt.SigningMethodHMAC = jwt.SigningMethodHS256
)

// UserManager manages user accounts
type UserManager struct {
	DB *sql.DB
}

// ValidateUser ...
// given a UserSpec, return if it's valid
func ValidateUser(db *sql.DB, user types.UserSpec, allowEmptyPassword bool) (valid bool, err error) {
	if len(user.Names) == 0 || len(user.Names) > 60 || user.Names == "" {
		return false, fmt.Errorf("Unable to use the provided name, as it is either empty or too long or too short")
	}
	if !common.RegexMatchEmail(user.Email) || user.Email == "" {
		return false, fmt.Errorf("Unable to use the provided email, as it is either empty or not valid")
	}

	if len(user.Groups) == 0 {
		return false, fmt.Errorf("No groups provided; please select at least one group")
	}
	groupsIncludeFlatmember := false
	for _, groupItem := range user.Groups {
		if groupItem == "flatmember" {
			groupsIncludeFlatmember = true
		}
		group, err := groups.GetGroupByName(db, groupItem)
		if err != nil || group.ID == "" {
			return false, fmt.Errorf("Unable to use the provided group '%v' as it is invalid", groupItem)
		}
	}
	if !groupsIncludeFlatmember {
		return false, fmt.Errorf("User account must be in the flatmember group")
	}

	if user.Birthday != 0 && !common.ValidateBirthday(user.Birthday) {
		return false, fmt.Errorf("Unable to use the provided birthday, your birthday year must not be within the last 15 years")
	}

	if (!common.RegexMatchPassword(user.Password) || user.Password == "") && !allowEmptyPassword {
		return false, fmt.Errorf("Unable to use the provided password, as it is either empty of invalid")
	}
	if user.PhoneNumber != "" && !common.RegexMatchPhoneNumber(user.PhoneNumber) {
		return false, fmt.Errorf("Unable to use the provided phone number")
	}

	return true, nil
}

// CreateUser ...
// given a UserSpec, create a user
func CreateUser(db *sql.DB, user types.UserSpec, allowEmptyPassword bool) (userInserted types.UserSpec, err error) {
	var userCreationSecretInserted types.UserCreationSecretSpec

	validUser, err := ValidateUser(db, user, allowEmptyPassword)
	if !validUser || err != nil {
		return types.UserSpec{}, err
	}
	localUser, err := GetUserByEmail(db, user.Email, false)
	if err == nil || localUser.ID != "" {
		return types.UserSpec{}, fmt.Errorf("Email address is unable to be used")
	}
	if localUser.Email == user.Email {
		return types.UserSpec{}, fmt.Errorf("Email address is already taken")
	}
	if user.Password != "" {
		user.Password = common.HashSHA512(user.Password)
	}
	if !allowEmptyPassword {
		user.Registered = true
	}

	sqlStatement := `insert into users (names, email, password, phonenumber, birthday, contractAgreement, disabled, registered)
                         values ($1, $2, $3, $4, $5, $6, $7, $8)
                         returning *`
	rows, err := db.Query(sqlStatement, user.Names, user.Email, user.Password, user.PhoneNumber, user.Birthday, user.ContractAgreement, user.Disabled, user.Registered)
	if err != nil {
		return types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	userInserted, err = UserObjectFromRows(rows)
	if err != nil {
		return types.UserSpec{}, err
	}
	updatedGroups, err := groups.UpdateUserGroups(db, userInserted.ID, user.Groups)
	if !updatedGroups || err != nil {
		return types.UserSpec{}, err
	}
	userInserted.Groups = user.Groups
	userInserted.Password = ""
	if allowEmptyPassword {
		userCreationSecretInserted, err = CreateUserCreationSecret(db, userInserted.ID)
		if err != nil {
			return types.UserSpec{}, err
		}
		if userCreationSecretInserted.ID == "" {
			return types.UserSpec{}, fmt.Errorf("Failed to created a user creation secret")
		}
	}

	return userInserted, nil
}

// GetAllUsers ...
// return all users in the database
func GetAllUsers(db *sql.DB, includePassword bool, selectors types.UserSelector) (users []types.UserSpec, err error) {
	sqlStatement := `select * from users where deletionTimestamp = 0 order by names`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return []types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		found := true
		user, err := UserObjectFromRows(rows)
		if err != nil {
			return []types.UserSpec{}, err
		}
		groupsOfUser, err := groups.GetGroupNamesOfUserByID(db, user.ID)
		if err != nil {
			return []types.UserSpec{}, err
		}
		if selectors.Group != "" {
			found, err = groups.CheckUserInGroup(db, user.ID, selectors.Group)
			if err != nil {
				return []types.UserSpec{}, err
			}
		} else if selectors.ID != "" {
			found = selectors.ID == user.ID
		} else if selectors.NotID != "" {
			found = !(selectors.NotID == user.ID)
		}
		user.Groups = groupsOfUser
		if !includePassword {
			user.Password = ""
		}
		if !found {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUser ...
// given a UserSpec and an ID or Email, return a user from the database
func GetUser(db *sql.DB, userSelect types.UserSpec, includePassword bool) (user types.UserSpec, err error) {
	if userSelect.ID != "" {
		return GetUserByID(db, userSelect.ID, includePassword)
	}
	if userSelect.Email != "" {
		if common.RegexMatchEmail(userSelect.Email) {
			return types.UserSpec{}, fmt.Errorf("Invalid email address")
		}
		return GetUserByEmail(db, userSelect.Email, includePassword)
	}
	if !includePassword {
		user.Password = ""
	}
	return user, nil
}

// UserObjectFromRowsRestricted ...
// construct a restricted UserSpec from database rows
func UserObjectFromRowsRestricted(rows *sql.Rows) (user types.UserSpec, err error) {
	if err := rows.Scan(&user.ID, &user.Names, &user.Email, &user.PhoneNumber, &user.Birthday, &user.ContractAgreement, &user.Disabled, &user.Registered, &user.LastLogin, &user.CreationTimestamp, &user.ModificationTimestamp, &user.DeletionTimestamp); err != nil {
		return types.UserSpec{}, err
	}
	if err := rows.Err(); err != nil {
		return types.UserSpec{}, err
	}
	return user, nil
}

// UserObjectFromRows ...
// construct a UserSpec from database rows
func UserObjectFromRows(rows *sql.Rows) (user types.UserSpec, err error) {
	if err := rows.Scan(&user.ID, &user.Names, &user.Email, &user.Password, &user.PhoneNumber, &user.Birthday, &user.ContractAgreement, &user.Disabled, &user.Registered, &user.LastLogin, &user.AuthNonce, &user.CreationTimestamp, &user.ModificationTimestamp, &user.DeletionTimestamp); err != nil {
		return types.UserSpec{}, err
	}
	if err := rows.Err(); err != nil {
		return types.UserSpec{}, err
	}
	return user, nil
}

// GetUserByID ...
// given an id, return a UserSpec
func GetUserByID(db *sql.DB, id string, includePassword bool) (user types.UserSpec, err error) {
	sqlStatement := `select * from users where id = $1`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	user, err = UserObjectFromRows(rows)
	if err != nil {
		return types.UserSpec{}, err
	}
	groups, err := groups.GetGroupNamesOfUserByID(db, user.ID)
	if err != nil {
		return types.UserSpec{}, err
	}
	user.Groups = groups
	if !includePassword {
		user.Password = ""
	}
	return user, nil
}

// GetUserByEmail ...
// given a email, return a UserSpec
func GetUserByEmail(db *sql.DB, email string, includePassword bool) (user types.UserSpec, err error) {
	sqlStatement := `select * from users where email = $1`
	rows, err := db.Query(sqlStatement, email)
	if err != nil {
		return types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	user, err = UserObjectFromRows(rows)
	if err != nil {
		return types.UserSpec{}, err
	}
	if user.ID == "" {
		return types.UserSpec{}, fmt.Errorf("Failed to find user")
	}
	groups, err := groups.GetGroupNamesOfUserByID(db, user.ID)
	if err != nil {
		return types.UserSpec{}, err
	}
	user.Groups = groups
	if !includePassword {
		user.Password = ""
	}
	return user, nil
}

// DeleteUserByID ...
// given an id, remove the user account from all the groups and then delete a user account
func DeleteUserByID(db *sql.DB, id string) (err error) {
	userGroups, err := groups.GetGroupsOfUserByID(db, id)
	if err != nil {
		return err
	}
	for _, groupItem := range userGroups {
		err = groups.RemoveUserFromGroup(db, id, groupItem.ID)
		if err != nil {
			return err
		}
	}
	err = DeleteUserCreationSecretByUserID(db, id)
	if err != nil {
		return err
	}
	sqlStatement := `update users set email = '', password = '', deletionTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	return nil
}

// CheckUserPassword ...
// given an email and password, find the user account with the email, return if the password matches
func CheckUserPassword(db *sql.DB, email string, password string) (matches bool, err error) {
	user, err := GetUserByEmail(db, email, true)
	if err != nil {
		return false, err
	}
	passwordHashed := common.HashSHA512(password)
	return user.Password == passwordHashed, nil
}

// GenerateJWTauthToken ...
// given an email, return a usable JWT token
func (u UserManager) GenerateJWTauthToken(id string, authNonce string, expiresIn time.Duration) (tokenString string, err error) {
	systemManager := &system.Manager{DB: u.DB}
	if expiresIn == 0 {
		expiresIn = 24 * 5
	}
	secret, err := systemManager.GetJWTsecret()
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(time.Hour * expiresIn)
	token := jwt.NewWithClaims(jwtAlg, types.JWTclaim{
		ID:        id,
		AuthNonce: authNonce,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	})

	tokenString, err = token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetAuthTokenFromHeader ...
// given a request, retreive the authoriation value
func GetAuthTokenFromHeader(r *http.Request) (string, error) {
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		return "", fmt.Errorf("Unable to find authorization token (header doesn't exist)")
	}
	authorizationHeader := strings.Split(tokenHeader, " ")
	if authorizationHeader[0] != "bearer" || len(authorizationHeader) <= 1 {
		return "", fmt.Errorf("Unable to find authorization token (must be as bearer)")
	}
	return authorizationHeader[1], nil
}

// ValidateJWTauthToken ...
// given an HTTP request and Authorization header, return if auth is valid
func ValidateJWTauthToken(db *sql.DB, r *http.Request) (valid bool, tokenClaims *types.JWTclaim, err error) {
	systemManager := &system.Manager{DB: db}
	secret, err := systemManager.GetJWTsecret()
	if err != nil {
		return false, &types.JWTclaim{}, fmt.Errorf("Unable to find FlatTrack system auth secret. Please contact system administrators or support")
	}
	tokenHeaderJWT, err := GetAuthTokenFromHeader(r)
	if err != nil {
		return false, &types.JWTclaim{}, err
	}
	claims := &types.JWTclaim{}
	token, err := jwt.ParseWithClaims(tokenHeaderJWT, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return false, &types.JWTclaim{}, err
	}

	if !token.Valid {
		return false, &types.JWTclaim{}, fmt.Errorf("Unable to use existing login token as it is invalid")
	}
	if err := token.Claims.Valid(); err != nil {
		return false, &types.JWTclaim{}, fmt.Errorf("Unable to use existing login token as it is invalid")
	}
	if token.Method.Alg() != jwtAlg.Alg() {
		return false, &types.JWTclaim{}, fmt.Errorf("Unable to use login token provided, please log in again")
	}

	reqClaims, ok := token.Claims.(*types.JWTclaim)
	if !ok {
		return false, &types.JWTclaim{}, fmt.Errorf("Unable to read JWT claims")
	}
	user, err := GetUserByID(db, reqClaims.ID, true)
	if err != nil || user.ID == "" || user.DeletionTimestamp != 0 {
		return false, &types.JWTclaim{}, fmt.Errorf("Unable to find the user account which the authentication token belongs to")
	}

	if reqClaims.AuthNonce != user.AuthNonce {
		return false, &types.JWTclaim{}, fmt.Errorf("Authentication has been invalidated, please log in again")
	}

	if user.Disabled {
		return false, &types.JWTclaim{}, fmt.Errorf("Your user account is disabled")
	}

	return token.Valid, reqClaims, nil
}

// InvalidateAllAuthTokens ...
// updates the authNonce to invalidate auth tokens
func InvalidateAllAuthTokens(db *sql.DB, id string) (err error) {
	sqlStatement := `update users set authNonce = md5(random()::text || clock_timestamp()::text)::uuid where id = $1`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	return nil
}

// GetIDFromJWT ...
// return the userID in a JWT from a header in a HTTP request
func GetIDFromJWT(db *sql.DB, r *http.Request) (id string, err error) {
	systemManager := &system.Manager{DB: db}
	secret, err := systemManager.GetJWTsecret()
	if err != nil {
		return "", fmt.Errorf("Unable to find FlatTrack system auth secret. Please contact system administrators or support")
	}
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		return "", fmt.Errorf("Unable to find authorization token (header doesn't exist)")
	}
	authorizationHeader := strings.Split(tokenHeader, " ")
	if authorizationHeader[0] != "bearer" || len(authorizationHeader) <= 1 {
		return "", fmt.Errorf("Unable to find authorization token (must be as bearer)")
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
	user, err := GetUserByID(db, reqClaims.ID, true)
	if err != nil || user.ID == "" {
		log.Printf("error getting user by ID; %v\n", err)
		return "", fmt.Errorf("Unable to find the user account which the authentication token belongs to")
	}
	return user.ID, nil
}

// GetProfile ...
// return user from ID in JWT from HTTP request
func GetProfile(db *sql.DB, r *http.Request) (user types.UserSpec, err error) {
	id, err := GetIDFromJWT(db, r)
	if err != nil {
		return types.UserSpec{}, err
	}
	user, err = GetUserByID(db, id, false)
	if err != nil {
		return types.UserSpec{}, err
	}
	return user, nil
}

// PatchProfile ...
// patches the profile of a user account
func PatchProfile(db *sql.DB, id string, userAccount types.UserSpec) (userAccountPatched types.UserSpec, err error) {
	existingUserAccount, err := GetUserByID(db, id, true)
	if err != nil || existingUserAccount.ID == "" {
		return types.UserSpec{}, fmt.Errorf("Failed to find user account")
	}
	if userAccount.Email != "" && userAccount.Email != existingUserAccount.Email {
		localUser, err := GetUserByEmail(db, userAccount.Email, false)
		if err == nil || localUser.ID != "" {
			return types.UserSpec{}, fmt.Errorf("Email address is unable to be used")
		}
	}
	err = mergo.Merge(&userAccount, existingUserAccount)
	if err != nil {
		return types.UserSpec{}, fmt.Errorf("Failed to update fields in the user account")
	}
	noUpdatePassword := userAccount.Password == existingUserAccount.Password
	valid, err := ValidateUser(db, userAccount, noUpdatePassword)
	if !valid || err != nil {
		return types.UserSpec{}, err
	}
	passwordHashed := common.HashSHA512(userAccount.Password)
	if noUpdatePassword {
		passwordHashed = userAccount.Password
	}

	sqlStatement := `update users set names = $2, email = $3, password = $4, phoneNumber = $5, birthday = $6, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1
                         returning id, names, email, phoneNumber, birthday, contractAgreement, disabled, registered, lastLogin, creationTimestamp, modificationTimestamp, deletionTimestamp`
	rows, err := db.Query(sqlStatement, id, userAccount.Names, userAccount.Email, passwordHashed, userAccount.PhoneNumber, userAccount.Birthday)
	if err != nil {
		// TODO add roll back, if there's failure
		return types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	userAccountPatched, err = UserObjectFromRowsRestricted(rows)
	if err != nil || userAccountPatched.ID == "" {
		return types.UserSpec{}, fmt.Errorf("Failed to patch user account")
	}

	userAccountPatched.Groups = existingUserAccount.Groups
	userAccountPatched.Password = ""
	return userAccountPatched, nil
}

// PatchProfileAdmin ...
// patches a profile with all fields
func PatchProfileAdmin(db *sql.DB, id string, userAccount types.UserSpec) (userAccountPatched types.UserSpec, err error) {
	existingUserAccount, err := GetUserByID(db, id, true)
	if err != nil || existingUserAccount.ID == "" {
		return types.UserSpec{}, fmt.Errorf("Failed to find user account")
	}
	if userAccount.Email != "" && userAccount.Email != existingUserAccount.Email {
		localUser, err := GetUserByEmail(db, userAccount.Email, false)
		if err == nil || localUser.ID != "" {
			return types.UserSpec{}, fmt.Errorf("Email address is unable to be used")
		}
	}
	err = mergo.Merge(&userAccount, existingUserAccount)
	if err != nil {
		return types.UserSpec{}, fmt.Errorf("Failed to update fields in the user account")
	}
	noUpdatePassword := userAccount.Password == existingUserAccount.Password
	valid, err := ValidateUser(db, userAccount, noUpdatePassword)
	if !valid || err != nil {
		return types.UserSpec{}, err
	}
	passwordHashed := common.HashSHA512(userAccount.Password)
	if noUpdatePassword {
		passwordHashed = userAccount.Password
	}

	sqlStatement := `update users set names = $1, email = $2, password = $3, phoneNumber = $4, birthday = $5, contractAgreement = $6, registered = $7, lastLogin = $8, authNonce = $9, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $10
                         returning *`
	rows, err := db.Query(sqlStatement, userAccount.Names, userAccount.Email, passwordHashed, userAccount.PhoneNumber, userAccount.Birthday, userAccount.ContractAgreement, userAccount.Registered, userAccount.LastLogin, userAccount.AuthNonce, id)
	if err != nil {
		// TODO add roll back, if there's failure
		return types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	userAccountPatched, err = UserObjectFromRows(rows)
	if err != nil || userAccountPatched.ID == "" {
		log.Printf("error getting user object from rows: %v\n", err)
		return types.UserSpec{}, fmt.Errorf("Failed to patch user account")
	}

	updatedGroups, err := groups.UpdateUserGroups(db, id, userAccount.Groups)
	if !updatedGroups || err != nil {
		return types.UserSpec{}, err
	}

	userAccountPatched.Password = ""
	userAccountPatched.Groups = userAccount.Groups
	return userAccountPatched, nil
}

// UpdateProfile ...
// updates the profile of a user account
func UpdateProfile(db *sql.DB, id string, userAccount types.UserSpec) (userAccountUpdated types.UserSpec, err error) {
	valid, err := ValidateUser(db, userAccount, false)
	if !valid || err != nil {
		return types.UserSpec{}, err
	}
	existingUserAccount, err := GetUserByID(db, id, true)
	if err != nil || existingUserAccount.ID == "" {
		return types.UserSpec{}, fmt.Errorf("Failed to find user account")
	}
	if userAccount.Email != existingUserAccount.Email {
		localUser, err := GetUserByEmail(db, userAccount.Email, false)
		if err == nil || localUser.ID != "" {
			return types.UserSpec{}, fmt.Errorf("Email address is unable to be used")
		}
	}
	passwordHashed := common.HashSHA512(userAccount.Password)

	sqlStatement := `update users set names = $2, email = $3, password = $4, phoneNumber = $5, birthday = $6, contractAgreement = $7, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1
                         returning id, names, email, phoneNumber, birthday, contractAgreement, disabled, registered, lastLogin, creationTimestamp, modificationTimestamp, deletionTimestamp`
	rows, err := db.Query(sqlStatement, id, userAccount.Names, userAccount.Email, passwordHashed, userAccount.PhoneNumber, userAccount.Birthday, userAccount.ContractAgreement)
	if err != nil {
		// TODO add roll back, if there's failure
		return types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	userAccountUpdated, err = UserObjectFromRowsRestricted(rows)
	if err != nil || userAccountUpdated.ID == "" {
		log.Printf("error getting user object from rows (restricted): %v\n", err)
		return types.UserSpec{}, fmt.Errorf("Failed to update profile")
	}
	userAccountUpdated.Groups = existingUserAccount.Groups
	return userAccountUpdated, nil
}

// UpdateProfileAdmin ...
// updates all fields of a profile
func UpdateProfileAdmin(db *sql.DB, id string, userAccount types.UserSpec) (userAccountUpdated types.UserSpec, err error) {
	valid, err := ValidateUser(db, userAccount, false)
	if !valid || err != nil {
		return types.UserSpec{}, err
	}
	existingUserAccount, err := GetUserByID(db, id, true)
	if err != nil || existingUserAccount.ID == "" {
		return types.UserSpec{}, fmt.Errorf("Failed to find user account")
	}
	if userAccount.Email != existingUserAccount.Email {
		localUser, err := GetUserByEmail(db, userAccount.Email, false)
		if err == nil || localUser.ID != "" {
			return types.UserSpec{}, fmt.Errorf("Email address is unable to be used")
		}
	}
	passwordHashed := common.HashSHA512(userAccount.Password)

	sqlStatement := `update users set names = $2, email = $3, password = $4, phoneNumber = $5, birthday = $6, contractAgreement = $7, registered = $8, lastLogin = $9, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1
                         returning *`
	rows, err := db.Query(sqlStatement, id, userAccount.Names, userAccount.Email, passwordHashed, userAccount.PhoneNumber, userAccount.Birthday, userAccount.ContractAgreement, userAccount.Registered, userAccount.LastLogin)
	if err != nil {
		// TODO add roll back, if there's failure
		return types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	userAccountUpdated, err = UserObjectFromRows(rows)
	if err != nil || userAccountUpdated.ID == "" {
		return types.UserSpec{}, fmt.Errorf("Failed to update profile")
	}

	updatedGroups, err := groups.UpdateUserGroups(db, id, userAccount.Groups)
	if !updatedGroups || err != nil {
		// TODO this could potentially panic, if err is nil
		return types.UserSpec{}, err
	}
	userAccountUpdated.Groups = userAccount.Groups

	userAccountUpdated.Password = ""
	return userAccountUpdated, nil
}

// UserCreationSecretsFromRows ...
// constructs a UserCreationSecretSpec from rows
func UserCreationSecretsFromRows(rows *sql.Rows) (creationSecret types.UserCreationSecretSpec, err error) {
	if err := rows.Scan(&creationSecret.ID, &creationSecret.UserID, &creationSecret.Secret, &creationSecret.Valid, creationSecret.CreationTimestamp, &creationSecret.ModificationTimestamp, &creationSecret.DeletionTimestamp); err != nil {
		return types.UserCreationSecretSpec{}, err
	}
	if err := rows.Err(); err != nil {
		return types.UserCreationSecretSpec{}, err
	}
	return creationSecret, nil
}

// GetAllUserCreationSecrets ...
// returns all UserCreationSecrets from the database
func GetAllUserCreationSecrets(db *sql.DB, secretsSelector types.UserCreationSecretSelector) (creationSecrets []types.UserCreationSecretSpec, err error) {
	sqlStatement := `select * from user_creation_secret`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return []types.UserCreationSecretSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		creationSecret, err := UserCreationSecretsFromRows(rows)
		if err != nil {
			return []types.UserCreationSecretSpec{}, fmt.Errorf("Failed to list user creation secrets")
		}
		if secretsSelector.UserID != "" {
			userExists, err := UserAccountExists(db, secretsSelector.UserID)
			if err != nil {
				return []types.UserCreationSecretSpec{}, err
			}
			if !userExists {
				return []types.UserCreationSecretSpec{}, fmt.Errorf("Unable to find user account")
			}
			if creationSecret.UserID != secretsSelector.UserID {
				continue
			}
		}
		creationSecrets = append(creationSecrets, creationSecret)
	}
	return creationSecrets, nil
}

// GetUserCreationSecret ...
// returns a UserCreationSecret by it's id from the database
func GetUserCreationSecret(db *sql.DB, id string) (creationSecret types.UserCreationSecretSpec, err error) {
	sqlStatement := `select * from user_creation_secret where id = $1`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return types.UserCreationSecretSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	creationSecret, err = UserCreationSecretsFromRows(rows)
	if err != nil {
		return types.UserCreationSecretSpec{}, err
	}
	return creationSecret, nil
}

// CreateUserCreationSecret ...
// creates a user creation secret for account confirming
func CreateUserCreationSecret(db *sql.DB, userID string) (userCreationSecretInserted types.UserCreationSecretSpec, err error) {
	sqlStatement := `insert into user_creation_secret (userId)
                         values ($1)
                         returning *`
	rows, err := db.Query(sqlStatement, userID)
	if err != nil {
		return types.UserCreationSecretSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	userCreationSecretInserted, err = UserCreationSecretsFromRows(rows)
	if err != nil {
		return types.UserCreationSecretSpec{}, err
	}
	return userCreationSecretInserted, nil
}

// DeleteUserCreationSecret ...
// deletes the acccount creation secret, after it's been used
func DeleteUserCreationSecret(db *sql.DB, id string) (err error) {
	sqlStatement := `delete from user_creation_secret where id = $1`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	return nil
}

// DeleteUserCreationSecretByUserID ...
// deletes the acccount creation secret by userid, after it's been used
func DeleteUserCreationSecretByUserID(db *sql.DB, userID string) (err error) {
	sqlStatement := `delete from user_creation_secret where userId = $1`
	rows, err := db.Query(sqlStatement, userID)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	return nil
}

// ConfirmUserAccount ...
// confirms the user account
func (u UserManager) ConfirmUserAccount(id string, secret string, user types.UserSpec) (tokenString string, err error) {
	if user.Password == "" {
		return "", fmt.Errorf("Unable to confirm account, a password must be provided to complete registration")
	}
	userCreationSecret, err := GetUserCreationSecret(u.DB, id)
	if err != nil {
		return "", err
	}
	if userCreationSecret.ID == "" {
		return "", fmt.Errorf("Unable to find account confirmation secret")
	}
	if secret != userCreationSecret.Secret {
		return "", fmt.Errorf("Unable to confirm account, as the secret doesn't match")
	}
	userInDB, err := GetUserByID(u.DB, userCreationSecret.UserID, false)
	if err != nil {
		return "", err
	}

	userAccountPatch := types.UserSpec{
		Names:       user.Names,
		Email:       user.Email,
		Password:    user.Password,
		Birthday:    user.Birthday,
		PhoneNumber: user.PhoneNumber,
		Registered:  true,
	}
	userConfirmed, err := PatchProfileAdmin(u.DB, userCreationSecret.UserID, userAccountPatch)
	if err != nil {
		return "", err
	}
	if userConfirmed.ID == "" || !userConfirmed.Registered {
		return "", fmt.Errorf("Failed to patch profile")
	}
	err = DeleteUserCreationSecret(u.DB, id)
	if err != nil {
		return "", err
	}
	tokenString, err = u.GenerateJWTauthToken(userInDB.ID, userInDB.AuthNonce, 0)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// UserAccountExists ...
// returns bool if user account exists
func UserAccountExists(db *sql.DB, id string) (exists bool, err error) {
	sqlStatement := `select id from users where id = $1`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	var userIDFromDB string
	if err := rows.Scan(&userIDFromDB); err != nil {
		return false, err
	}
	err = rows.Err()
	if err != nil {
		return false, err
	}
	return userIDFromDB == id, nil
}

// GenerateNewAuthNonce ...
// given a user account id, generates a new auth nonce to reset all logins and invalidate all issued JWTs
func GenerateNewAuthNonce(db *sql.DB, id string) (err error) {
	sqlStatement := `update users set authNonce = md5(random()::text || clock_timestamp()::text)::uuid where id = $1
                         returning *`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	_, err = UserObjectFromRows(rows)
	if err != nil {
		return err
	}
	return nil
}

// PatchUserDisabledAdmin ...
// patches is user account to be disabled
func PatchUserDisabledAdmin(db *sql.DB, id string, disabled bool) (userAccount types.UserSpec, err error) {
	sqlStatement := `update users set disabled = $2 where id = $1
                         returning *`
	rows, err := db.Query(sqlStatement, id, disabled)
	if err != nil {
		return userAccount, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	userAccount, err = UserObjectFromRows(rows)
	if err != nil {
		return types.UserSpec{}, err
	}
	return userAccount, nil
}
