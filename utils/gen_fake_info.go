package utils

import (
	"encoding/json"
	"net/http"
)

// The GenerateFakeInfo() function is used to generate a fake
// name and email that will be used in creating a new spotify
// account.
func GenerateFakeInfo(client *http.Client) (string, string) {
	// Establish a new request object
	var req, _ = http.NewRequest("GET", "https://api.namefake.com/", nil)

	// Set the request header
	req.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	// Define Variables
	var (
		// Readable golang map
		data map[string]interface{}
		// Send the http request
		resp, _ = client.Do(req)
	)
	// Decode the response data
	json.NewDecoder(resp.Body).Decode(&data)

	// Return the fake name and email
	var email = data["username"].(string) + RandomString(7) + "@" + data["email_d"].(string)
	return data["name"].(string), email
}
