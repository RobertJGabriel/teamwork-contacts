package teamwork

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Connection struct {
	Account struct {
		AvatarURL                  string `json:"avatar-url"`
		CanAddProjects             string `json:"canaddprojects"`
		CanManagePeople            string `json:"canManagePeople"`
		ChatEnabled                bool   `json:"chatEnabled"`
		Code                       string `json:"code"`
		CompanyID                  string `json:"companyid"`
		CompanyName                string `json:"companyname"`
		DateFormat                 string `json:"dateFormat"`
		DateSeperator              string `json:"dateSeperator"`
		DeskEnabled                bool   `json:"deskEnabled"`
		DocumentEditorEnabled      bool   `json:"documentEditorEnabled"`
		FirstName                  string `json:"firstname"`
		ID                         string `json:"id"`
		Lang                       string `json:"lang"`
		LastName                   string `json:"lastname"`
		LikesEnabled               bool   `json:"likesEnabled"`
		Logo                       string `json:"logo"`
		Name                       string `json:"name"`
		PlanID                     string `json:"plan-id"`
		ProjectsEnabled            bool   `json:"projectsEnabled"`
		RequireHTTPS               bool   `json:"requirehttps"`
		SslEnabled                 bool   `json:"ssl-enabled"`
		StartOnSundays             bool   `json:"startonsundays"`
		TagsEnabled                bool   `json:"tagsEnabled"`
		TagsLockedToAdmins         bool   `json:"tagsLockedToAdmins"`
		TimeFormat                 string `json:"timeFormat"`
		TkoEnabled                 bool   `json:"TKOEnabled"`
		URL                        string `json:"URL"`
		UserID                     string `json:"userId"`
		UserIsAdmin                bool   `json:"userIsAdmin"`
		UserIsMemberOfOwnerCompany string `json:"userIsMemberOfOwnerCompany"`
	} `json:"account"`
	APIToken string
}

// Connect is the starting point to using the TeamWork API.
// This function returns a Connection which is used to query
// TeamWork via other functions.
func Connect(APIToken string) (*Connection, error) {
	method := "GET"
	url := "http://authenticate.teamworkpm.net/authenticate.json"
	reader, _, err := request(APIToken, method, url)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	connection := &Connection{
		APIToken: APIToken,
	}
	if err := json.NewDecoder(reader).Decode(connection); err != nil {
		return nil, err
	}
	return connection, nil
}

// request is the base level function for calling the TeamWork API.
func request(token, method, url string) (io.ReadCloser, http.Header, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(token, "notused")

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	return resp.Body, resp.Header, nil
}

func buildParams(ops interface{}) string {
	pairs := make([]string, 0)
	v := reflect.ValueOf(ops).Elem()
	for i := 0; i < v.NumField(); i++ {
		var paramValue string
		paramName := v.Type().Field(i).Tag.Get("param") // get value from struct field tag

		isPointer := false
		var kind reflect.Kind
		// Handle either strings or pointers
		switch {
		case v.Field(i).Kind() == reflect.Ptr:
			kind = v.Field(i).Elem().Kind()
			isPointer = true
		case v.Field(i).Kind() == reflect.String:
			paramValue = v.Field(i).Interface().(string)
		}

		// handle pointers
		switch {
		case isPointer && kind == reflect.String:
			if v.Field(i).Interface() != nil {
				paramValue = *v.Field(i).Interface().(*string)
			}
		case isPointer && kind == reflect.Bool:
			if v.Field(i).Interface() != nil {
				paramValue = strconv.FormatBool(*v.Field(i).Interface().(*bool))
			}
		case isPointer && kind == reflect.Int:
			if v.Field(i).Interface() != nil {
				paramValue = strconv.FormatInt(int64(*v.Field(i).Interface().(*int)), 10)
			}
		case isPointer && kind == reflect.Int8:
			if v.Field(i).Interface() != nil {
				paramValue = strconv.FormatInt(int64(*v.Field(i).Interface().(*int8)), 10)
			}
		case isPointer && kind == reflect.Int16:
			if v.Field(i).Interface() != nil {
				paramValue = strconv.FormatInt(int64(*v.Field(i).Interface().(*int16)), 10)
			}
		case isPointer && kind == reflect.Int32:
			if v.Field(i).Interface() != nil {
				paramValue = strconv.FormatInt(int64(*v.Field(i).Interface().(*int32)), 10)
			}
		case isPointer && kind == reflect.Int64:
			if v.Field(i).Interface() != nil {
				paramValue = strconv.FormatInt(*v.Field(i).Interface().(*int64), 10)
			}
		case isPointer && kind == reflect.Float32:
			if v.Field(i).Interface() != nil {
				paramValue = strconv.FormatFloat(float64(*v.Field(i).Interface().(*float32)), 'f', -1, 64)
			}
		case isPointer && kind == reflect.Float64:
			if v.Field(i).Interface() != nil {
				paramValue = strconv.FormatFloat(*v.Field(i).Interface().(*float64), 'f', -1, 64)
			}
		}
		if paramName != "" && paramValue != "" { // make sure we have what we need to set a param
			pair := fmt.Sprintf("%s=%s", paramName, paramValue)
			pairs = append(pairs, pair) // add to the param pairs array
		}
	}
	if len(pairs) > 0 {
		return fmt.Sprintf("?%s", strings.Join(pairs, "&")) // return the params with the leading '?'
	}
	return "" // nothing to send back

}

// getHeaders takes the response headers and populates
// a struct of data according to the `header:"HeaderName"`.
// Function currently only supports Int and String field types.
func getHeaders(headers http.Header, obj interface{}) {
	v := reflect.ValueOf(obj).Elem()
	if v.Kind() == reflect.Struct { // make sure we have a struct
		for i := 0; i < v.NumField(); i++ { // for all fields
			field := v.Field(i)                    // value field.
			if field.IsValid() && field.CanSet() { // is exported and addressable
				headerName := v.Type().Field(i).Tag.Get("header") // get value from struct field tag
				if headerName != "" {                             // make sure the header is set
					headerVal := headers.Get(headerName)
					if headerVal != "" { // make sure we have a value in the header
						switch {
						case field.Kind() == reflect.Int: // Int struct field type
							hVal, err := strconv.ParseInt(headerVal, 10, 64)
							if err != nil {
								log.Printf("Failed to convert header '%s' to a 64 bit Int. \n%s", headerName, err.Error())
								continue
							}
							if !field.OverflowInt(hVal) {
								field.SetInt(hVal)
							}
						case field.Kind() == reflect.String: // String struct field type
							field.SetString(headerVal)
						}
					}
				}
			}
		}
	}
}
