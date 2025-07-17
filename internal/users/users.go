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
	"crypto/subtle"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/imdario/mergo"

	jwt "github.com/golang-jwt/jwt/v5"
	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/internal/groups"
	"gitlab.com/flattrack/flattrack/internal/system"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

var (
	jwtAlg *jwt.SigningMethodHMAC = jwt.SigningMethodHS256
)

var (
	ErrAuthInvalid                                       = fmt.Errorf("Authentication has been invalidated, please log in again")
	ErrEmailAddressAlreadyUsed                           = fmt.Errorf("Email address is unable to be used")
	ErrFailedToCreateUserCreationSecret                  = fmt.Errorf("Failed to create a user creation secret")
	ErrFailedToFindAccount                               = fmt.Errorf("Failed to find user account")
	ErrFailedToFindUser                                  = fmt.Errorf("Failed to find user")
	ErrFailedToListUserCreationSecrets                   = fmt.Errorf("Failed to list user creation secrets")
	ErrFailedToPatchProfile                              = fmt.Errorf("Failed to patch profile")
	ErrFailedToPatchUserAccount                          = fmt.Errorf("Failed to patch user account")
	ErrFailedToUpdateUserAccount                         = fmt.Errorf("Failed to update user account")
	ErrFailedToUpdateProfile                             = fmt.Errorf("Failed to update profile")
	ErrInvalidEmailAddress                               = fmt.Errorf("Invalid email address")
	ErrNoGroupsProvided                                  = fmt.Errorf("No groups provided; please select at least one group")
	ErrUserAccountConfirmPasswordRequiredForRegistration = fmt.Errorf("Unable to confirm account, a password must be provided to complete registration")
	ErrUserAccountConfirmSecretDoesNotMatch              = fmt.Errorf("Unable to confirm account, as the secret doesn't match")
	ErrFailedToFindSystemAuthSecret                      = fmt.Errorf("Unable to find FlatTrack system auth secret. Please contact system administrators or support")
	ErrFailedToFindAccountConfirmSecret                  = fmt.Errorf("Unable to find account confirmation secret")
	ErrFailedToFindAuthToken                             = fmt.Errorf("Unable to find authorization token")
	ErrFailedToFindAuthTokenAccountID                    = fmt.Errorf("Unable to find the user account which the authentication token belongs to")
	ErrFailedToFindUserAccount                           = fmt.Errorf("Unable to find user account")
	ErrFailedToReadJWTClaims                             = fmt.Errorf("Unable to read JWT claims")
	ErrAuthTokenExpired                                  = fmt.Errorf("Unable to use existing login token as it is invalid")
	ErrAuthTokenFailed                                   = fmt.Errorf("Unable to use login token provided, please log in again")
	ErrUserAccountNotFound                               = fmt.Errorf("Failed to find user account")
	ErrUserAccountInvalidBirthday                        = fmt.Errorf("Unable to use the provided birthday, your birthday year must not be within the last 15 years")
	ErrUserAccountInvalidEmail                           = fmt.Errorf("Unable to use the provided email, as it is either empty or not valid")
	ErrUserAccountInvalidGroup                           = fmt.Errorf("Unable to use the provided group as it is invalid")
	ErrUserAccountInvalidName                            = fmt.Errorf("Unable to use the provided name, as it is either empty or too long or too short")
	ErrUserAccountInvalidPassword                        = fmt.Errorf("Unable to use the provided password, as it is either empty of invalid")
	ErrUserAccountInvalidPhoneNumber                     = fmt.Errorf("Unable to use the provided phone number")
	ErrUserAccountMustBeInFlatmemberGroup                = fmt.Errorf("User account must be in the flatmember group")
	ErrUserAccountIsDisabled                             = fmt.Errorf("Your user account is disabled")
	ErrAuthorizationHeaderNotFound                       = fmt.Errorf("Unable to find authorization token (header doesn't exist)")
	ErrUserAccountCreationSecretNotFound                 = fmt.Errorf("Failed to find user account creation secret")
)

// UserManager manages user accounts
type Manager struct {
	groups *groups.Manager
	system *system.Manager
	db     *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{
		groups: groups.NewManager(db),
		system: system.NewManager(db),
		db:     db,
	}
}

// ValidateUser ...
// given a UserSpec, return if it's valid
func (m *Manager) Validate(user types.UserSpec, allowEmptyPassword bool) (valid bool, err error) {
	if len(user.Names) == 0 || len(user.Names) > 60 || user.Names == "" {
		return false, ErrUserAccountInvalidName
	}
	if !common.RegexMatchEmail(user.Email) || user.Email == "" {
		return false, ErrInvalidEmailAddress
	}

	if len(user.Groups) == 0 {
		return false, ErrNoGroupsProvided
	}
	groupsIncludeFlatmember := false
	for _, groupItem := range user.Groups {
		if groupItem == "flatmember" {
			groupsIncludeFlatmember = true
		}
		group, err := m.groups.GetByName(groupItem)
		if err != nil || group.ID == "" {
			return false, ErrUserAccountInvalidGroup
		}
	}
	if !groupsIncludeFlatmember {
		return false, ErrUserAccountMustBeInFlatmemberGroup
	}

	if user.Birthday != 0 && !common.ValidateBirthday(user.Birthday) {
		return false, ErrUserAccountInvalidBirthday
	}

	if (!common.RegexMatchPassword(user.Password) || user.Password == "") && !allowEmptyPassword {
		return false, ErrUserAccountInvalidPassword
	}
	if user.PhoneNumber != "" && !common.RegexMatchPhoneNumber(user.PhoneNumber) {
		return false, ErrUserAccountInvalidPassword
	}

	return true, nil
}

// Create ...
// given a UserSpec, create a user
func (m *Manager) Create(user types.UserSpec, allowEmptyPassword bool) (userInserted types.UserSpec, err error) {
	var userCreationSecretInserted types.UserCreationSecretSpec

	validUser, err := m.Validate(user, allowEmptyPassword)
	if !validUser || err != nil {
		return types.UserSpec{}, err
	}
	localUser, err := m.GetByEmail(user.Email, false)
	if err == nil || localUser.Email == user.Email || localUser.ID != "" {
		return types.UserSpec{}, ErrEmailAddressAlreadyUsed
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
	rows, err := m.db.Query(sqlStatement, user.Names, user.Email, user.Password, user.PhoneNumber, user.Birthday, user.ContractAgreement, user.Disabled, user.Registered)
	if err != nil {
		return types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		userInserted, err = userObjectFromRows(rows)
		if err != nil {
			return types.UserSpec{}, err
		}
	}
	if err := m.groups.UpdateUserGroups(userInserted.ID, user.Groups); err != nil {
		return types.UserSpec{}, err
	}
	userInserted.Groups = user.Groups
	userInserted.Password = ""
	if allowEmptyPassword {
		if userCreationSecretInserted, err = m.UserCreationSecrets().Create(userInserted.ID); err != nil {
			return types.UserSpec{}, err
		}
		if userCreationSecretInserted.ID == "" {
			return types.UserSpec{}, ErrFailedToCreateUserCreationSecret
		}
	}

	return userInserted, nil
}

// List ...
// return all users in the database
func (m *Manager) List(includePassword bool, selectors types.UserSelector) (users []types.UserSpec, err error) {
	sqlStatement := `select * from users where deletionTimestamp = 0 order by names`
	rows, err := m.db.Query(sqlStatement)
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
		user, err := userObjectFromRows(rows)
		if err != nil {
			return []types.UserSpec{}, err
		}
		groupsOfUser, err := m.groups.GetGroupNamesOfUserByID(user.ID)
		if err != nil {
			return []types.UserSpec{}, err
		}
		if selectors.Group != "" {
			found, err = m.groups.CheckUserInGroup(user.ID, selectors.Group)
			if err != nil {
				return []types.UserSpec{}, err
			}
		} else if selectors.ID != "" {
			found = selectors.ID == user.ID
		} else if selectors.NotID != "" {
			found = selectors.NotID != user.ID
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

// Get ...
// given a UserSpec and an ID or Email, return a user from the database
func (m *Manager) Get(userSelect types.UserSpec, includePassword bool) (user types.UserSpec, err error) {
	if userSelect.ID != "" {
		return m.GetByID(userSelect.ID, includePassword)
	}
	if userSelect.Email != "" {
		if common.RegexMatchEmail(userSelect.Email) {
			return types.UserSpec{}, ErrInvalidEmailAddress
		}
		return m.GetByEmail(userSelect.Email, includePassword)
	}
	if !includePassword {
		user.Password = ""
	}
	return user, nil
}

// userObjectFromRowsRestricted ...
// construct a restricted UserSpec from database rows
func userObjectFromRowsRestricted(rows *sql.Rows) (user types.UserSpec, err error) {
	if err := rows.Scan(&user.ID, &user.Names, &user.Email, &user.PhoneNumber, &user.Birthday, &user.ContractAgreement, &user.Disabled, &user.Registered, &user.LastLogin, &user.CreationTimestamp, &user.ModificationTimestamp, &user.DeletionTimestamp); err != nil {
		return types.UserSpec{}, err
	}
	if err := rows.Err(); err != nil {
		return types.UserSpec{}, err
	}
	return user, nil
}

// userObjectFromRows ...
// construct a UserSpec from database rows
func userObjectFromRows(rows *sql.Rows) (user types.UserSpec, err error) {
	if err := rows.Scan(&user.ID, &user.Names, &user.Email, &user.Password, &user.PhoneNumber, &user.Birthday, &user.ContractAgreement, &user.Disabled, &user.Registered, &user.LastLogin, &user.AuthNonce, &user.CreationTimestamp, &user.ModificationTimestamp, &user.DeletionTimestamp); err != nil {
		return types.UserSpec{}, err
	}
	if err := rows.Err(); err != nil {
		return types.UserSpec{}, err
	}
	return user, nil
}

// GetByID ...
// given an id, return a UserSpec
func (m *Manager) GetByID(id string, includePassword bool) (user types.UserSpec, err error) {
	sqlStatement := `select * from users where id = $1`
	rows, err := m.db.Query(sqlStatement, id)
	if err != nil {
		return types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		user, err = userObjectFromRows(rows)
		if err != nil {
			return types.UserSpec{}, err
		}
	}
	if user.ID == "" {
		return types.UserSpec{}, ErrUserAccountNotFound
	}
	groups, err := m.groups.GetGroupNamesOfUserByID(user.ID)
	if err != nil {
		return types.UserSpec{}, err
	}
	user.Groups = groups
	if !includePassword {
		user.Password = ""
	}
	return user, nil
}

// GetByEmail ...
// given a email, return a UserSpec
func (m *Manager) GetByEmail(email string, includePassword bool) (user types.UserSpec, err error) {
	sqlStatement := `select * from users where email = $1`
	rows, err := m.db.Query(sqlStatement, email)
	if err != nil {
		return types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		user, err = userObjectFromRows(rows)
		if err != nil {
			return types.UserSpec{}, err
		}
	}
	if user.ID == "" {
		return types.UserSpec{}, ErrFailedToFindAccount
	}
	groups, err := m.groups.GetGroupNamesOfUserByID(user.ID)
	if err != nil {
		return types.UserSpec{}, err
	}
	user.Groups = groups
	if !includePassword {
		user.Password = ""
	}
	return user, nil
}

// DeleteByID ...
// given an id, remove the user account from all the groups and then delete a user account
func (m *Manager) DeleteByID(id string) (err error) {
	userGroups, err := m.groups.GetGroupsOfUserByID(id)
	if err != nil {
		return err
	}
	for _, groupItem := range userGroups {
		err = m.groups.RemoveUserFromGroup(id, groupItem.ID)
		if err != nil {
			return err
		}
	}
	err = m.UserCreationSecrets().DeleteByUserID(id)
	if err != nil {
		return err
	}
	sqlStatement := `update users set email = '', password = '', deletionTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1`
	rows, err := m.db.Query(sqlStatement, id)
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
func (m *Manager) CheckUserPassword(email string, password string) (matches bool, err error) {
	user, err := m.GetByEmail(email, true)
	if err != nil {
		return false, err
	}
	passwordHashed := common.HashSHA512(password)
	if matches := subtle.ConstantTimeCompare([]byte(user.Password), []byte(passwordHashed)) == 1; matches {
		return true, nil
	}
	return false, nil
}

// GenerateJWTauthToken ...
// given an email, return a usable JWT token
func (m *Manager) GenerateJWTauthToken(id string, authNonce string, expiresIn time.Duration) (tokenString string, err error) {
	if expiresIn == 0 {
		expiresIn = 24 * 5
	}
	secret, err := m.system.GetJWTsecret()
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(time.Hour * expiresIn)
	token := jwt.NewWithClaims(jwtAlg, types.JWTclaim{
		ID:        id,
		AuthNonce: authNonce,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	})

	tokenString, err = token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetAuthTokenFromHeader ...
// given a request, retrieve the authoriation value
func GetAuthTokenFromHeader(r *http.Request) (string, error) {
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		return "", ErrAuthorizationHeaderNotFound
	}
	authorizationHeader := strings.Split(tokenHeader, " ")
	if authorizationHeader[0] != "bearer" || len(authorizationHeader) <= 1 {
		return "", ErrAuthorizationHeaderNotFound
	}
	return authorizationHeader[1], nil
}

// ValidateJWTauthToken ...
// given an HTTP request and Authorization header, return if auth is valid
func (m *Manager) ValidateJWTauthToken(r *http.Request) (valid bool, tokenClaims *types.JWTclaim, err error) {
	secret, err := m.system.GetJWTsecret()
	if err != nil {
		return false, &types.JWTclaim{}, ErrFailedToFindSystemAuthSecret
	}
	tokenHeaderJWT, err := GetAuthTokenFromHeader(r)
	if err != nil {
		return false, &types.JWTclaim{}, err
	}
	claims := &types.JWTclaim{}
	token, err := jwt.ParseWithClaims(tokenHeaderJWT, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		if token.Method.Alg() != jwtAlg.Alg() {
			return nil, ErrAuthTokenFailed
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, &types.JWTclaim{}, err
	}

	if !token.Valid {
		return false, &types.JWTclaim{}, ErrAuthInvalid
	}
	reqClaims, ok := token.Claims.(*types.JWTclaim)
	if !ok {
		return false, &types.JWTclaim{}, ErrFailedToReadJWTClaims
	}
	user, err := m.GetByID(reqClaims.ID, true)
	if err != nil || user.ID == "" || user.DeletionTimestamp != 0 {
		return false, &types.JWTclaim{}, ErrFailedToFindAuthTokenAccountID
	}

	if reqClaims.AuthNonce != user.AuthNonce {
		return false, &types.JWTclaim{}, ErrAuthInvalid
	}

	if user.Disabled {
		return false, &types.JWTclaim{}, ErrUserAccountIsDisabled
	}

	return token.Valid, reqClaims, nil
}

// InvalidateAllAuthTokens ...
// updates the authNonce to invalidate auth tokens
func (m *Manager) InvalidateAllAuthTokens(id string) (err error) {
	sqlStatement := `update users set authNonce = md5(random()::text || clock_timestamp()::text)::uuid where id = $1`
	rows, err := m.db.Query(sqlStatement, id)
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
// TODO move into internal/httpserver/common.go
func (m *Manager) GetIDFromJWT(r *http.Request) (id string, err error) {
	valid, claims, err := m.ValidateJWTauthToken(r)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", ErrAuthInvalid
	}
	return claims.ID, nil
}

// GetProfile ...
// return user from ID in JWT from HTTP request
// TODO move into internal/httpserver/common.go
func (m *Manager) GetProfile(r *http.Request) (user types.UserSpec, err error) {
	id, err := m.GetIDFromJWT(r)
	if err != nil {
		return types.UserSpec{}, err
	}
	user, err = m.GetByID(id, false)
	if err != nil {
		return types.UserSpec{}, err
	}
	return user, nil
}

// Patch ...
// patches the profile of a user account
func (m *Manager) Patch(id string, userAccount types.UserSpec) (userAccountPatched types.UserSpec, err error) {
	existingUserAccount, err := m.GetByID(id, true)
	if err != nil || existingUserAccount.ID == "" {
		return types.UserSpec{}, ErrFailedToFindAccount
	}
	if userAccount.Email != "" && userAccount.Email != existingUserAccount.Email {
		localUser, err := m.GetByEmail(userAccount.Email, false)
		if err == nil || localUser.ID != "" {
			return types.UserSpec{}, ErrEmailAddressAlreadyUsed
		}
	}
	err = mergo.Merge(&userAccount, existingUserAccount)
	if err != nil {
		return types.UserSpec{}, ErrFailedToPatchProfile
	}
	noUpdatePassword := userAccount.Password == existingUserAccount.Password
	valid, err := m.Validate(userAccount, noUpdatePassword)
	if !valid || err != nil {
		return types.UserSpec{}, err
	}
	passwordHashed := common.HashSHA512(userAccount.Password)
	if noUpdatePassword {
		passwordHashed = userAccount.Password
	}

	sqlStatement := `update users set names = $2, email = $3, password = $4, phoneNumber = $5, birthday = $6, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1
                         returning id, names, email, phoneNumber, birthday, contractAgreement, disabled, registered, lastLogin, creationTimestamp, modificationTimestamp, deletionTimestamp`
	rows, err := m.db.Query(sqlStatement, id, userAccount.Names, userAccount.Email, passwordHashed, userAccount.PhoneNumber, userAccount.Birthday)
	if err != nil {
		// TODO add roll back, if there's failure
		return types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		userAccountPatched, err = userObjectFromRowsRestricted(rows)
		if err != nil {
			log.Println("error patching user account:", err)
			return types.UserSpec{}, ErrFailedToPatchProfile
		}
	}
	if userAccountPatched.ID == "" {
		return types.UserSpec{}, ErrFailedToPatchProfile
	}

	userAccountPatched.Groups = existingUserAccount.Groups
	userAccountPatched.Password = ""
	return userAccountPatched, nil
}

// PatchAsAdmin ...
// patches a profile with all fields
func (m *Manager) PatchAsAdmin(id string, userAccount types.UserSpec) (userAccountPatched types.UserSpec, err error) {
	existingUserAccount, err := m.GetByID(id, true)
	if err != nil || existingUserAccount.ID == "" {
		return types.UserSpec{}, ErrFailedToFindUserAccount
	}
	if userAccount.Email != "" && userAccount.Email != existingUserAccount.Email {
		localUser, err := m.GetByEmail(userAccount.Email, false)
		if err == nil || localUser.ID != "" {
			return types.UserSpec{}, ErrEmailAddressAlreadyUsed
		}
	}
	err = mergo.Merge(&userAccount, existingUserAccount)
	if err != nil {
		return types.UserSpec{}, ErrFailedToPatchProfile
	}
	noUpdatePassword := userAccount.Password == existingUserAccount.Password
	valid, err := m.Validate(userAccount, noUpdatePassword)
	if !valid || err != nil {
		return types.UserSpec{}, err
	}
	passwordHashed := common.HashSHA512(userAccount.Password)
	if noUpdatePassword {
		passwordHashed = userAccount.Password
	}

	sqlStatement := `update users set names = $1, email = $2, password = $3, phoneNumber = $4, birthday = $5, contractAgreement = $6, registered = $7, lastLogin = $8, authNonce = $9, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $10
                         returning *`
	rows, err := m.db.Query(sqlStatement, userAccount.Names, userAccount.Email, passwordHashed, userAccount.PhoneNumber, userAccount.Birthday, userAccount.ContractAgreement, userAccount.Registered, userAccount.LastLogin, userAccount.AuthNonce, id)
	if err != nil {
		// TODO add roll back, if there's failure
		return types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		userAccountPatched, err = userObjectFromRows(rows)
		if err != nil {
			log.Printf("error getting user object from rows: %v\n", err)
			return types.UserSpec{}, ErrFailedToPatchUserAccount
		}
	}
	if userAccountPatched.ID == "" {
		return types.UserSpec{}, ErrFailedToPatchUserAccount
	}

	if err := m.groups.UpdateUserGroups(id, userAccount.Groups); err != nil {
		return types.UserSpec{}, err
	}

	userAccountPatched.Password = ""
	userAccountPatched.Groups = userAccount.Groups
	return userAccountPatched, nil
}

// Update ...
// updates the profile of a user account
func (m *Manager) Update(id string, userAccount types.UserSpec) (userAccountUpdated types.UserSpec, err error) {
	valid, err := m.Validate(userAccount, false)
	if !valid || err != nil {
		return types.UserSpec{}, err
	}
	existingUserAccount, err := m.GetByID(id, true)
	if err != nil || existingUserAccount.ID == "" {
		return types.UserSpec{}, ErrFailedToFindAccount
	}
	if userAccount.Email != existingUserAccount.Email {
		localUser, err := m.GetByEmail(userAccount.Email, false)
		if err == nil || localUser.ID != "" {
			return types.UserSpec{}, ErrEmailAddressAlreadyUsed
		}
	}
	passwordHashed := common.HashSHA512(userAccount.Password)

	sqlStatement := `update users set names = $2, email = $3, password = $4, phoneNumber = $5, birthday = $6, contractAgreement = $7, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1
                         returning id, names, email, phoneNumber, birthday, contractAgreement, disabled, registered, lastLogin, creationTimestamp, modificationTimestamp, deletionTimestamp`
	rows, err := m.db.Query(sqlStatement, id, userAccount.Names, userAccount.Email, passwordHashed, userAccount.PhoneNumber, userAccount.Birthday, userAccount.ContractAgreement)
	if err != nil {
		// TODO add roll back, if there's failure
		return types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		userAccountUpdated, err = userObjectFromRowsRestricted(rows)
		if err != nil {
			log.Printf("error getting user object from rows (restricted): %v\n", err)
			return types.UserSpec{}, ErrFailedToUpdateProfile
		}
	}
	if userAccountUpdated.ID == "" {
		return types.UserSpec{}, ErrFailedToUpdateProfile
	}
	userAccountUpdated.Groups = existingUserAccount.Groups
	return userAccountUpdated, nil
}

// UpdateAsAdmin ...
// updates all fields of a profile
func (m *Manager) UpdateAsAdmin(id string, userAccount types.UserSpec) (userAccountUpdated types.UserSpec, err error) {
	valid, err := m.Validate(userAccount, false)
	if !valid || err != nil {
		return types.UserSpec{}, err
	}
	existingUserAccount, err := m.GetByID(id, true)
	if err != nil || existingUserAccount.ID == "" {
		return types.UserSpec{}, ErrFailedToFindAccount
	}
	if userAccount.Email != existingUserAccount.Email {
		localUser, err := m.GetByEmail(userAccount.Email, false)
		if err == nil || localUser.ID != "" {
			return types.UserSpec{}, ErrEmailAddressAlreadyUsed
		}
	}
	passwordHashed := common.HashSHA512(userAccount.Password)

	sqlStatement := `update users set names = $2, email = $3, password = $4, phoneNumber = $5, birthday = $6, contractAgreement = $7, registered = $8, lastLogin = $9, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1
                         returning *`
	rows, err := m.db.Query(sqlStatement, id, userAccount.Names, userAccount.Email, passwordHashed, userAccount.PhoneNumber, userAccount.Birthday, userAccount.ContractAgreement, userAccount.Registered, userAccount.LastLogin)
	if err != nil {
		// TODO add roll back, if there's failure
		return types.UserSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		userAccountUpdated, err = userObjectFromRows(rows)
		if userAccountUpdated.ID == "" {
			return types.UserSpec{}, ErrFailedToUpdateProfile
		}
	}
	if err != nil {
		return types.UserSpec{}, ErrFailedToUpdateProfile
	}

	if err := m.groups.UpdateUserGroups(id, userAccount.Groups); err != nil {
		// TODO this could potentially panic, if err is nil
		return types.UserSpec{}, err
	}
	userAccountUpdated.Groups = userAccount.Groups

	userAccountUpdated.Password = ""
	return userAccountUpdated, nil
}

// userCreationSecretsFromRows ...
// constructs a UserCreationSecretSpec from rows
func userCreationSecretsFromRows(rows *sql.Rows) (creationSecret types.UserCreationSecretSpec, err error) {
	if err := rows.Scan(&creationSecret.ID, &creationSecret.UserID, &creationSecret.Secret, &creationSecret.Valid, &creationSecret.CreationTimestamp, &creationSecret.ModificationTimestamp, &creationSecret.DeletionTimestamp); err != nil {
		return types.UserCreationSecretSpec{}, err
	}
	if err := rows.Err(); err != nil {
		return types.UserCreationSecretSpec{}, err
	}
	return creationSecret, nil
}

type userCreationSecretManager struct {
	db *sql.DB
	m  *Manager
}

func (m *Manager) UserCreationSecrets() *userCreationSecretManager {
	return &userCreationSecretManager{
		db: m.db,
		m:  m,
	}
}

// List ...
// returns all UserCreationSecrets from the database
func (m *userCreationSecretManager) List(secretsSelector types.UserCreationSecretSelector) (creationSecrets []types.UserCreationSecretSpec, err error) {
	sqlStatement := `select * from user_creation_secret`
	rows, err := m.db.Query(sqlStatement)
	if err != nil {
		return []types.UserCreationSecretSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		creationSecret, err := userCreationSecretsFromRows(rows)
		if err != nil {
			return []types.UserCreationSecretSpec{}, ErrFailedToListUserCreationSecrets
		}
		if secretsSelector.UserID != "" {
			userExists, err := m.m.UserAccountExists(secretsSelector.UserID)
			if err != nil {
				return []types.UserCreationSecretSpec{}, err
			}
			if !userExists {
				return []types.UserCreationSecretSpec{}, ErrFailedToFindUserAccount
			}
			if creationSecret.UserID != secretsSelector.UserID {
				continue
			}
		}
		creationSecrets = append(creationSecrets, creationSecret)
	}
	return creationSecrets, nil
}

// Get ...
// returns a UserCreationSecret by it's id from the database
func (m *userCreationSecretManager) Get(id string) (creationSecret types.UserCreationSecretSpec, err error) {
	sqlStatement := `select * from user_creation_secret where id = $1`
	rows, err := m.db.Query(sqlStatement, id)
	if err != nil {
		return types.UserCreationSecretSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		creationSecret, err = userCreationSecretsFromRows(rows)
		if err != nil {
			return types.UserCreationSecretSpec{}, err
		}
	}
	if creationSecret.ID == "" {
		return types.UserCreationSecretSpec{}, ErrUserAccountCreationSecretNotFound
	}
	return creationSecret, nil
}

// Create ...
// creates a user creation secret for account confirming
func (m *userCreationSecretManager) Create(userID string) (userCreationSecretInserted types.UserCreationSecretSpec, err error) {
	sqlStatement := `insert into user_creation_secret (userId)
                         values ($1)
                         returning *`
	rows, err := m.db.Query(sqlStatement, userID)
	if err != nil {
		return types.UserCreationSecretSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		userCreationSecretInserted, err = userCreationSecretsFromRows(rows)
		if err != nil {
			return types.UserCreationSecretSpec{}, err
		}
	}
	return userCreationSecretInserted, nil
}

// Delete ...
// deletes the acccount creation secret, after it's been used
func (m *userCreationSecretManager) Delete(id string) (err error) {
	sqlStatement := `delete from user_creation_secret where id = $1`
	rows, err := m.db.Query(sqlStatement, id)
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

// DeleteByUserID ...
// deletes the acccount creation secret by userid, after it's been used
func (m *userCreationSecretManager) DeleteByUserID(userID string) (err error) {
	sqlStatement := `delete from user_creation_secret where userId = $1`
	rows, err := m.db.Query(sqlStatement, userID)
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
func (m *Manager) ConfirmUserAccount(id string, secret string, user types.UserSpec) (tokenString string, err error) {
	if user.Password == "" {
		return "", ErrUserAccountConfirmPasswordRequiredForRegistration
	}
	userCreationSecret, err := m.UserCreationSecrets().Get(id)
	if err != nil {
		return "", err
	}
	if userCreationSecret.ID == "" {
		return "", ErrFailedToFindAccountConfirmSecret
	}
	if secret != userCreationSecret.Secret {
		return "", ErrUserAccountConfirmSecretDoesNotMatch
	}
	userInDB, err := m.GetByID(userCreationSecret.UserID, false)
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
	userConfirmed, err := m.PatchAsAdmin(userCreationSecret.UserID, userAccountPatch)
	if err != nil {
		return "", err
	}
	if userConfirmed.ID == "" || !userConfirmed.Registered {
		return "", ErrFailedToPatchProfile
	}
	err = m.UserCreationSecrets().Delete(id)
	if err != nil {
		return "", err
	}
	tokenString, err = m.GenerateJWTauthToken(userInDB.ID, userInDB.AuthNonce, 0)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// UserAccountExists ...
// returns bool if user account exists
func (m *Manager) UserAccountExists(id string) (exists bool, err error) {
	sqlStatement := `select id from users where id = $1`
	rows, err := m.db.Query(sqlStatement, id)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	var userIDFromDB string
	for rows.Next() {
		if err := rows.Scan(&userIDFromDB); err != nil {
			return false, err
		}
		err = rows.Err()
		if err != nil {
			return false, err
		}
	}
	return userIDFromDB == id, nil
}

// GenerateNewAuthNonce ...
// given a user account id, generates a new auth nonce to reset all logins and invalidate all issued JWTs
func (m *Manager) GenerateNewAuthNonce(id string) (err error) {
	sqlStatement := `update users set authNonce = md5(random()::text || clock_timestamp()::text)::uuid where id = $1
                         returning *`
	rows, err := m.db.Query(sqlStatement, id)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		_, err = userObjectFromRows(rows)
		if err != nil {
			return err
		}
	}
	return nil
}

// PatchDisabledAsAdmin ...
// patches is user account to be disabled
func (m *Manager) PatchDisabledAsAdmin(id string, disabled bool) (userAccount types.UserSpec, err error) {
	sqlStatement := `update users set disabled = $2 where id = $1
                         returning *`
	rows, err := m.db.Query(sqlStatement, id, disabled)
	if err != nil {
		return userAccount, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		userAccount, err = userObjectFromRows(rows)
		if err != nil {
			return types.UserSpec{}, err
		}
	}
	return userAccount, nil
}
