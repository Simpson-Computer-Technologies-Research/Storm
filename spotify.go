package main

// Import Packages
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"

	Utils "github.com/realTristan/SpotifyBooster/utils"
)

// Define Global Variables
var (
	// The data inside the spotify/data.json file
	JsonData map[string]interface{} = Utils.ReadJsonFile("spotify/data.json")
	// The Available Artist IDs to Follow
	ArtistList []interface{} = JsonData["artistIds"].([]interface{})
)

// The FollowRandomArtists() function is used to follow random
// artists. This functions primary use is to make the bot account
// look more real. Following other accounts help secure that realty.
func FollowRandomArtists(client *http.Client, bearer string, amount int) {
	// Store already selected artist ids
	var alreadySelected []string

	// Iterate over the amount
	for i := 0; i < amount; i++ {
		// Get a random artist id
		var artistId string = ArtistList[rand.Intn(len(ArtistList))].(string)

		// While the artist is already selected, select a new one
		for Utils.SliceContains(alreadySelected, artistId) {
			artistId = ArtistList[rand.Intn(len(ArtistList))].(string)
		}
		// Add the artist to the already selected slice
		alreadySelected = append(alreadySelected, artistId)

		// Establish the request object being used to follow that artist id
		var req, _ = http.NewRequest("PUT", fmt.Sprintf("https://api.spotify.com/v1/me/following?type=artist&ids=%s", artistId), nil)
		req.Header = http.Header{
			"app-platform":        {"WebPlayer"},
			"authorization":       {bearer},
			"spotify-app-version": {"1.1.88.584.gb23b6713"},
			"User-Agent":          {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
		}

		// Send the request and close the body
		var resp, _ = client.Do(req)
		resp.Body.Close()
	}
}

// The UploadProfileImage() function is used to upload a randomized image
// to the spotify database. This function returns the upload token
// which will be used by the UpdateProfileImage() function in changing
// the fake accounts profile image.
func UploadProfileImage(client *http.Client, bearer string) string {
	// Define Variables
	var (
		// Get a random image from api
		image io.Reader = Utils.GetRandomImage(client)
		// Get the image bytes
		imageBytes, _ = io.ReadAll(image)
	)

	// Establish the request object that will be used
	// for updating the profile image
	var req, _ = http.NewRequest("POST", "https://image-upload.spotify.com/v4/user-profile", image)
	req.Header = http.Header{
		"Content-Length": {fmt.Sprint(len(imageBytes))},
		"Authorization":  {bearer},
		"Content-Type":   {"image/jpeg"},
		"User-Agent":     {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}

	// Define Variables
	var (
		// Send the http request
		resp, _ = client.Do(req)
		// Read the response body bytes
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
		// Get the upload token in the response body
		uploadToken string = strings.Split(strings.Split(string(bodyBytes), "uploadToken\":\"")[1], "\"")[0]
	)
	// Close the body once the upload token is returned
	defer resp.Body.Close()
	return uploadToken
}

// The UploadProfileImage() function is used to change the fake accounts
// profile image. This functions primary use is too make the fake account
// look more real. Having a real profile image secures that realty.
func UpdateProfileImage(client *http.Client, userId string, bearer string) {
	// Define Variables
	var (
		// Get a random image token
		imageToken string = UploadProfileImage(client, bearer)
		// Establish a new request object for updating the profile image
		req, _ = http.NewRequest("POST", fmt.Sprintf("https://spclient.wg.spotify.com/identity/v3/profile-image/%s/%s", userId, imageToken), nil)
	)
	// Update the request object headers
	req.Header = http.Header{
		"App-Platform":        {"WebPlayer"},
		"Authorization":       {bearer},
		"Spotify-App-Version": {"1.1.88.584.gb23b6713"},
		"User-Agent":          {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	// Send the http request and close the
	// response body
	var resp, _ = client.Do(req)
	resp.Body.Close()
}

// The GenerateSpotifyData() function is used to generate the necessary
// data used in creating a new account. This function is called by the
// CreateNewAccount() function and returns a bytes.Buffer object.
func GenerateSpotifyData(username string, email string, password string) *bytes.Buffer {
	// Marshal the map written below
	var data, _ = json.Marshal(map[string]interface{}{
		"account_details": map[string]interface{}{
			"birthdate":     "1994-04-12",
			"consent_flags": map[string]interface{}{"eula_agreed": true, "send_email": false, "third_party_email": true},
			"display_name":  username,
			"email_and_password_identifier": map[string]interface{}{
				"email":    email,
				"password": password,
			}, "gender": 1,
		},
		"callback_uri": "https://www.spotify.com/signup/challenge?forward_url=https%3A%2F%2Fopen.spotify.com%2F__noul__%3Fl2l%3D1%26nd%3D1&locale=uk",
		"client_info": map[string]interface{}{
			"api_key":         "a1e486e2729f46d6bb368d6b2bcda326",
			"app_version":     "v2",
			"installation_id": "5fc5952c-9c03-4ee4-8e0d-d58d52bb2e7b",
			"platform":        "www",
		}, "tracking": map[string]string{"creation_flow": "", "creation_point": "https://www.spotify.com/uk/", "referrer": ""},
	})
	// Return the data as a *bytes.Buffer object
	return bytes.NewBuffer(data)
}

// The CreateNewAccount() function is used to create a new
// spotify account. This account will be used to follow
// the user id provided at the start of the program.
func CreateNewAccount(client *http.Client, username string, email string, password string) (string, string) {
	// Establish a new request object
	var req, _ = http.NewRequest("POST",
		"https://spclient.wg.spotify.com/signup/public/v2/account/create",
		GenerateSpotifyData(username, email, password))

	// Update the request objects headers
	req.Header = http.Header{
		"Content-Type": {"application/json"},
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	// Define Variables
	var (
		// Send the http request
		resp, _ = client.Do(req)
		// Get the response body in a bytes slice
		bodyBytes, _ = io.ReadAll(resp.Body)
		// Index the body bytes to get the userId and the loginToken
		userId     string = string(bodyBytes)[4:32]
		loginToken string = string(bodyBytes)[34:56]
	)

	// Close the response body after returning
	// the userId and loginToken
	defer resp.Body.Close()
	return userId, loginToken
}

// The GetCSRFToken() function is used to get the csrf token
// required for making http requests to the spotify api endpoints
//
// Without this token, no requests to the endpoints can be made.
func GetCSRFToken(client *http.Client) string {
	// Establish a new request object
	var req, _ = http.NewRequest("GET",
		"https://www.spotify.com/uk/signup/?forward_url=https://accounts.spotify.com/en/status&sp_t_counter=1", nil)

	// Set the request object's headers
	req.Header = http.Header{
		"Content-Type": {"text/html"},
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	// Define Variables
	var (
		// Send the http request
		resp, _ = client.Do(req)
		// Get the response body in bytes
		bodyBytes, _ = io.ReadAll(resp.Body)
		// Scrape the csrfToken from the response body
		csrfToken string = strings.Split(strings.Split(string(bodyBytes), "csrfToken\":\"")[1], "\"")[0]
	)

	// Close the response body after returning
	// the csrfToken
	defer resp.Body.Close()
	return csrfToken
}

// The GetBearerToken() function is used to get the bearer
// authentication token required for making changes
// to the fake account. Without this token, the fake
// account can't follow the user id provided at the
// start of the program nor change profile info.
func GetBearerToken(client *http.Client) string {
	// Establish a new request object
	var req, _ = http.NewRequest("GET", "https://open.spotify.com/get_access_token", nil)

	// Set the request object's headers
	req.Header = http.Header{
		"accept":              {"application/json"},
		"spotify-app-version": {"1.1.52.204.ge43bc405"},
		"app-platform":        {"WebPlayer"},
		"User-Agent":          {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	// Define Variables
	var (
		// Send the http request
		resp, _ = client.Do(req)
		// Get the response body in bytes
		bodyBytes, _ = io.ReadAll(resp.Body)
		// Get the bearer token used for making changes to the fake account
		bearerToken string = "Bearer " + strings.Split(strings.Split(string(bodyBytes), "accessToken\":\"")[1], "\"")[0]
	)

	// Close the response body after returning
	// the bearerToken
	defer resp.Body.Close()
	return bearerToken
}

// The AuthenticateAccount() function is used to permanently add
// the newly created account to the spotify database. Without this
// function, the fake account would likely get banned or removed.
func AuthenticateAccount(client *http.Client, csrfToken string, loginToken string) (*http.Response, error) {
	// Define Variables
	var (
		// Request body data
		data, _ = json.Marshal(map[string]string{"splot": loginToken})
		// Establish a new request object
		req, _ = http.NewRequest("POST", "https://www.spotify.com/api/signup/authenticate", bytes.NewBuffer(data))
	)
	// Set the request headers
	req.Header = http.Header{
		"Accept":       {"application/json"},
		"Content-Type": {"application/json"},
		"x-csrf-token": {csrfToken},
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}

	// Send the http request
	return client.Do(req)
}

// The FollowUser() function is used to follow the user id
// provided at the start of the function.
func FollowUser(client *http.Client, bearerToken string, userId string) (*http.Response, error) {
	// Establish a new request object
	var req, _ = http.NewRequest("PUT", "https://api.spotify.com/v1/me/following?type=user&ids="+userId, nil)

	// Set the request object headers
	req.Header = http.Header{
		"Authorization":       {bearerToken},
		"spotify-app-version": {"1.1.87.24.g5db224d0"},
		"app-platform":        {"WebPlayer"},
		"User-Agent":          {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	// Send the http request
	return client.Do(req)
}
