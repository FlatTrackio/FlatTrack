package main

import (
	"log"
	"os"
	"strings"

	"gitlab.com/flattrack/flattrack/pkg/common"
	"gitlab.com/flattrack/flattrack/pkg/database"
	"gitlab.com/flattrack/flattrack/pkg/migrations"
	"gitlab.com/flattrack/flattrack/pkg/types"
	"gitlab.com/flattrack/flattrack/pkg/users"
)

var (
	defaultUserAccountGroups = []string{"flatmember", "admin"}
	defaultUserPassword      = "P@ssw0rd123!"
)

func guestsUserAccountsFromFlatString(input string, hostname string) (guests []types.UserSpec) {
	usernames := strings.Split(input, " ")
	for _, username := range usernames {
		guests = append(guests, types.UserSpec{
			Names: username,
			Email: username + "@" + hostname,
		})
	}
	return guests
}

func main() {
	log.Println("Preparing to add new guests")
	flattrackDevHostname := os.Getenv("FLATTRACK_DEV_HOSTNAME")
	userAccountsNamesFlat := os.Getenv("SHARINGIO_PAIR_GUEST_NAMES")
	userAccountFromFlatNames := guestsUserAccountsFromFlatString(userAccountsNamesFlat, flattrackDevHostname)

	var userAccounts []types.UserSpec
	if len(userAccountFromFlatNames) > 0 {
		userAccounts = userAccountFromFlatNames
	} else {
		log.Fatalln("No guests declared to create")
		return
	}

	log.Println("Connecting to local Postgres instance")
	dbUsername := common.GetDBusername()
	dbPassword := common.GetDBpassword()
	dbHostname := common.GetDBhost()
	dbDatabase := common.GetDBdatabase()
	db, err := database.DB(dbUsername, dbPassword, dbHostname, dbDatabase)
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = migrations.Migrate(db)
	if err != nil {
		log.Fatalln("migrations:", err)
		return
	}

	if len(userAccounts) == 0 {
		log.Fatalln("No user accounts found in environment")
		return
	}
	for _, user := range userAccounts {
		user.Groups = defaultUserAccountGroups
		user.Password = defaultUserPassword
		user.Registered = true
		log.Printf("Inserting '%v'/'%v'\n", user.Names, user.Email)
		_, err = users.CreateUser(db, user, false)
		if err != nil {
			log.Fatalln(err)
			return
		}
		log.Printf("Inserted '%v'/'%v'\n", user.Names, user.Email)
		log.Println("Completed processing account")
	}

	log.Println("It's collaboration time!")
}
