package teamwork

import (
	"encoding/json"
	"fmt"
)

// People A list of People.
type People []Person

// Person Gets all people
type Person struct {
	AddressCity          string `json:"address-city"`
	AddressCountry       string `json:"address-country"`
	AddressLine1         string `json:"address-line-1"`
	AddressLine2         string `json:"address-line-2"`
	AddressState         string `json:"address-state"`
	AddressZip           string `json:"address-zip"`
	Administrator        bool   `json:"administrator"`
	AvatarURL            string `json:"avatar-url"`
	CompanyID            string `json:"companyId"`
	EmailAddress         string `json:"email-address"`
	FirstName            string `json:"first-name"`
	ID                   string `json:"id"`
	LastName             string `json:"last-name"`
	PhoneNumberFax       string `json:"phone-number-fax"`
	PhoneNumberHome      string `json:"phone-number-home"`
	PhoneNumberMobile    string `json:"phone-number-mobile"`
	PhoneNumberOffice    string `json:"phone-number-office"`
	PhoneNumberOfficeExt string `json:"phone-number-office-ext"`
	TextFormat           string `json:"textFormat"`
	Title                string `json:"title"`
	Twitter              string `json:"twitter"`
	Facebook             string `json:"facebook"`
}

// GetPeopleOps is used to generate the query params for the
// GetPeople API call.
type GetPeopleOps struct {
	// Query people based on these values.

	// Pass this parameter to check if a user exists by email address.
	EmailAddress string `param:"emailaddress"`

	FullProfile *bool `param:"fullprofile"`

	ReturnProjectIds *bool `param:"returnProjectIds"`
}

func (conn *Connection) GetCompanyPeople(id string, ops *GetPeopleOps) (People, error) {
	people := make(People, 0)
	pages := &Pages{}
	params := buildParams(ops)
	method := "GET"
	URL := fmt.Sprintf("%scompanies/%s/people.json%s", conn.Account.URL, id, params)
	reader, headers, err := request(conn.APIToken, method, URL)
	if err != nil {
		return people, err
	}

	getHeaders(headers, pages)
	defer reader.Close()

	err = json.NewDecoder(reader).Decode(&struct {
		*People `json:"people"`
	}{&people})
	if err != nil {
		return people, err
	}

	return people, nil
}

// GetPerson gets a single person based on a person ID.
func (conn *Connection) GetPerson(id string) (Person, error) {
	person := &Person{}
	method := "GET"
	URL := fmt.Sprintf("%speople/%s.json", conn.Account.URL, id)
	reader, _, err := request(conn.APIToken, method, URL)
	if err != nil {
		return *person, err
	}

	defer reader.Close()

	err = json.NewDecoder(reader).Decode(&struct {
		*Person `json:"person"`
	}{person})
	if err != nil {
		return *person, err
	}

	return *person, nil
}
