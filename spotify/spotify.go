package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	Global "spotify_follower_booster/global"
	"strings"
)

// Function to follow random artists
func FollowRandomArtists(client *http.Client, bearer *string, amount int) {
	var (
		JsonData        map[string]interface{} = Global.ReadJsonFile("spotify/data.json")
		artistList      []interface{}          = JsonData["artistIds"].([]interface{})
		alreadySelected []string
	)
	for i := 0; i < amount; i++ {
		var artistId string = artistList[rand.Intn(len(artistList))].(string)
		for Global.SliceContains(alreadySelected, artistId) {
			artistId = artistList[rand.Intn(len(artistList))].(string)
		}
		alreadySelected = append(alreadySelected, artistId)
		var req, _ = http.NewRequest("PUT", fmt.Sprintf("https://api.spotify.com/v1/me/following?type=artist&ids=%s", artistId), nil)
		req.Header = http.Header{
			"app-platform":        {"WebPlayer"},
			"authorization":       {*bearer},
			"spotify-app-version": {"1.1.88.584.gb23b6713"},
			"User-Agent":          {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
		}
		var resp, _ = client.Do(req)
		defer resp.Body.Close()
	}
}

// Function to upload profile image to spotify's database
func UploadProfileImage(client *http.Client, bearer *string) *string {
	var (
		image         io.Reader = Global.GetRandomImage(client)
		imageBytes, _           = io.ReadAll(image)
		req, _                  = http.NewRequest("POST", "https://image-upload.spotify.com/v4/user-profile", image)
	)
	req.Header = http.Header{
		"Content-Length": {string(imageBytes)},
		"Authorization":  {*bearer},
		"Content-Type":   {"image/jpeg"},
		"User-Agent":     {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	var (
		resp, _             = client.Do(req)
		bodyBytes, _        = ioutil.ReadAll(resp.Body)
		uploadToken  string = strings.Split(strings.Split(string(bodyBytes), "uploadToken\":\"")[1], "\"")[0]
	)
	defer resp.Body.Close()
	return &uploadToken
}

// Function to Update the users profile image
func UpdateProfileImage(client *http.Client, userId *string, bearer *string) {
	var (
		imageToken *string = UploadProfileImage(client, bearer)
		req, _             = http.NewRequest("POST", fmt.Sprintf("https://spclient.wg.spotify.com/identity/v3/profile-image/%s/%s", *userId, *imageToken), nil)
	)
	req.Header = http.Header{
		"App-Platform":        {"WebPlayer"},
		"Authorization":       {*bearer},
		"Spotify-App-Version": {"1.1.88.584.gb23b6713"},
		"User-Agent":          {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	var resp, _ = client.Do(req)
	defer resp.Body.Close()
}

// Function to generate account create spotify data
func GenerateSpotifyData(username *string, email *string, password *string) *bytes.Buffer {
	var data, _ = json.Marshal(map[string]interface{}{
		"account_details": map[string]interface{}{
			"birthdate":     "1994-04-12",
			"consent_flags": map[string]interface{}{"eula_agreed": true, "send_email": false, "third_party_email": true},
			"display_name":  *username,
			"email_and_password_identifier": map[string]interface{}{
				"email":    *email,
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
	return bytes.NewBuffer(data)
}

// Function to create a new spotify account
func CreateNewAccount(client *http.Client, username *string, email *string, password string) (*string, *string) {
	var req, _ = http.NewRequest("POST",
		"https://spclient.wg.spotify.com/signup/public/v2/account/create",
		GenerateSpotifyData(username, email, &password))
	req.Header = http.Header{
		"Content-Type": {"application/json"},
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	var (
		resp, _             = client.Do(req)
		bodyBytes, _        = io.ReadAll(resp.Body)
		userId       string = string(bodyBytes)[4:32]
		loginToken   string = string(bodyBytes)[34:56]
	)
	defer resp.Body.Close()
	return &userId, &loginToken
}

// Function to get the csrf auth token
func GetCSRFToken(client *http.Client) *string {
	var req, _ = http.NewRequest("GET",
		"https://www.spotify.com/uk/signup/?forward_url=https://accounts.spotify.com/en/status&sp_t_counter=1", nil)
	req.Header = http.Header{
		"Content-Type": {"text/html"},
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	var (
		resp, _             = client.Do(req)
		bodyBytes, _        = io.ReadAll(resp.Body)
		csrfToken    string = strings.Split(strings.Split(string(bodyBytes), "csrfToken\":\"")[1], "\"")[0]
	)
	defer resp.Body.Close()
	return &csrfToken
}

// Function to get the bearere authentication token
func GetBearerToken(client *http.Client) *string {
	var req, _ = http.NewRequest("GET", "https://open.spotify.com/get_access_token", nil)
	req.Header = http.Header{
		"accept":              {"application/json"},
		"spotify-app-version": {"1.1.52.204.ge43bc405"},
		"app-platform":        {"WebPlayer"},
		"User-Agent":          {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	var (
		resp, _             = client.Do(req)
		bodyBytes, _        = io.ReadAll(resp.Body)
		bearerToken  string = "Bearer " + strings.Split(strings.Split(string(bodyBytes), "accessToken\":\"")[1], "\"")[0]
	)
	defer resp.Body.Close()
	return &bearerToken
}

// Function to follow the given user with the newly created account
func FollowUser(client *http.Client, bearerToken *string, userId *string) (*http.Response, error) {
	var req, _ = http.NewRequest("PUT", "https://api.spotify.com/v1/me/following?type=user&ids="+*userId, nil)
	req.Header = http.Header{
		"Authorization":       {*bearerToken},
		"spotify-app-version": {"1.1.87.24.g5db224d0"},
		"app-platform":        {"WebPlayer"},
		"User-Agent":          {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	return client.Do(req)
}

// Function to authenticate the account so it doesn't get banned
func AuthenticateAccount(client *http.Client, csrfToken *string, loginToken *string) (*http.Response, error) {
	var (
		data, _ = json.Marshal(map[string]string{"splot": *loginToken})
		req, _  = http.NewRequest("POST", "https://www.spotify.com/api/signup/authenticate", bytes.NewBuffer(data))
	)
	req.Header = http.Header{
		"Accept":       {"application/json"},
		"Content-Type": {"application/json"},
		"x-csrf-token": {*csrfToken},
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	return client.Do(req)
}
