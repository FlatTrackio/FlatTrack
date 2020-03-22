// manage user accounts
package users

import (
	"database/sql"
	"errors"
	"fmt"

	"gitlab.com/flattrack/flattrack/src/backend/common"
	"gitlab.com/flattrack/flattrack/src/backend/types"
)

func CreateUser(db *sql.DB, user types.UserSpec) (userInserted types.UserSpec, err error) {
	fmt.Println(user)
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

func GetUserById(db *sql.DB, id string) (user types.UserSpec, err error) {
	sqlStatement := `select * from users where id = $1;`
	fmt.Println(sqlStatement)
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return user, err
	}
	user, err = UserObjectFromRows(rows)
	return user, err
}

func GetUserByEmail(db *sql.DB, email string) (user types.UserSpec, err error) {
	sqlStatement := `select * from users where email = ?;`
	fmt.Println(sqlStatement)
	rows, err := db.Query(sqlStatement, email)
	if err != nil {
		return user, err
	}
	defer rows.Close()
	user, err = UserObjectFromRows(rows)
	return user, err
}

