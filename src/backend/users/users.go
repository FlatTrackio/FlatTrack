// manage user accounts
package users

import (
	"database/sql"

	//"gitlab.com/flattrack/flattrack/src/backend/database"
	"gitlab.com/flattrack/flattrack/src/backend/types"
)

func CreateUser(db *sql.DB, user types.UserSpec) (userInserted types.UserSpec, err error) {
	sqlStatement := `
  insert into users (names, email, groups, password, phonenumber)
  values ($1, $2, $3, $4, $5)
  returning (id, names, email, groups, phonenumber)`
	rows, err := db.Query(sqlStatement, user.Names, user.Email, user.Groups, user.Password, user.PhoneNumber)
	if err != nil {
		return userInserted, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var names string
		var email string
		var groups string
		var phoneNumber string
		rows.Scan(&id, &names, &email, &groups, &phoneNumber)
		userInserted = types.UserSpec{
			Id: id,
			Names: names,
			Email: email,
			Groups: groups,
			PhoneNumber: phoneNumber,
		}
	}
	return userInserted, err
}
