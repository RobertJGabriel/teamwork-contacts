package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/spf13/viper"
	"github.com/teamwork/teamwork-contacts/core/config"
	"github.com/teamwork/teamwork-contacts/core/teamwork"
)

var (
	conn *teamwork.Connection
)

// InitRoutes - sets up the routes we will respond to
func main() {

	// Load Viper
	config.Load()

	apiToken := viper.GetString("apiKey")

	GetContacts(apiToken)
}

// Get contacts and write them to a file...
func GetContacts(APIToken string) {

	conn, err := teamwork.Connect(APIToken)

	if err != nil {
		fmt.Printf("Error connecting to TeamWork: %s", err.Error())
		return
	}

	// get company people
	peopleOps := &teamwork.GetPeopleOps{}
	people, err := conn.GetCompanyPeople(conn.Account.CompanyID, peopleOps)
	if err != nil {
		fmt.Printf("Error getting Company People: %s", err.Error())
		return
	}

	jsonString, err := json.Marshal(people)

	if err != nil {
		fmt.Printf("Error Converting the string: %s", err.Error())
		return
	}

	// Unmarshal JSON data
	var p []teamwork.Person
	err = json.Unmarshal([]byte(jsonString), &p)

	if err != nil {
		fmt.Printf("Cannot Get user details: %s", err.Error())
		return
	}

	// Get user direction
	user, err := user.Current()

	if err != nil {
		fmt.Printf("Cannot Get user details: %s", err.Error())
		return
	}

	homedir := user.HomeDir

	// Create a csv file
	f, err := os.Create(homedir + "/Desktop/" + "people.csv")

	if err != nil {
		fmt.Printf("Error Creating the file: %s", err.Error())
	}

	defer f.Close()

	// Write Unmarshaled json data to CSV file
	w := csv.NewWriter(f)
	var record []string
	record = append(record, "Given Name")
	record = append(record, "Family Name")
	record = append(record, "Phone 1 - Value")
	record = append(record, "Phone 1 - Type")
	record = append(record, "E-mail 1 - Value")
	w.Write(record)

	for _, obj := range p {
		var record []string
		record = append(record, obj.FirstName)
		record = append(record, obj.LastName)
		record = append(record, obj.PhoneNumberMobile)
		record = append(record, "Mobile")
		record = append(record, obj.EmailAddress)
		w.Write(record)
	}

	w.Flush()
	return
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
