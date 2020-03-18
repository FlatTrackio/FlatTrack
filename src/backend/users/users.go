// manage user accounts
package users

import (
	"database/sql"

	//"gitlab.com/flattrack/flattrack/src/backend/database"
	"gitlab.com/flattrack/flattrack/src/backend/types"
)

func CreateUser(db *sql.DB, user types.UserSpec) (userInserted types.UserSpec, err error) {
	sqlStatement := `insert into users (names, email, groups, password, phonenumber)
                         values ($1, $2, $3, $4, $5)
                         returning id`
	rows, err := db.Query(sqlStatement, user.Names, user.Email, user.Groups, user.Password, user.PhoneNumber)
	if err != nil {
		return userInserted, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		rows.Scan(&id)
		userInserted = types.UserSpec{
			Id: id,
			Names: user.Names,
			Email: user.Email,
			Groups: user.Groups,
			PhoneNumber: user.PhoneNumber,
		}
	}
	return userInserted, err
}
