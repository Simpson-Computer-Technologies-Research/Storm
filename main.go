package main

// Import Packages
import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	Utils "github.com/realTristan/SpotifyBooster/utils"
)

// The HandleResponse() function is used to handle the
// new follow response.
func HandleResponse(resp *http.Response, iteration int, name string) {
	// If the response was a success
	if resp.StatusCode == 204 {
		// Print the success messagef
		fmt.Printf(" >> Added New Follower [%d] [%s]\n", iteration+1, name)
	} else

	// Print the response body and error status code
	{
		var bodyBytes, _ = io.ReadAll(resp.Body)
		fmt.Println(resp.StatusCode, ":", string(bodyBytes))
	}
}

// Main function
func main() {
	// Initialize a random seed
	rand.Seed(time.Now().UnixNano())

	// Define Variables
	var (
		// Amount of followers to add
		amount int = Utils.InputInt(" >> Amount: ")
		// Wait group for goroutines
		waitGroup sync.WaitGroup = sync.WaitGroup{}
	)

	// Add the amount of followers to the wait group
	waitGroup.Add(amount)

	// Start Goroutines
	for i := 0; i < amount; i++ {
		go func(iteration int) {
			// Define Variables
			var (
				// Create a new request client with cookies
				jar, _              = cookiejar.New(nil)
				client *http.Client = &http.Client{Jar: jar}

				// Create a new account (all variables below)
				name, email = Utils.GenerateFakeInfo(client)

				// Replace _ with newUserId if using line 70 (below variable)
				_, newUserLoginToken = CreateNewAccount(client, name, email, "secretpassword!")

				// Get a new csrfToken
				csrfToken string = GetCSRFToken(client)
			)

			// Authenticate new account (spotify.go)
			AuthenticateAccount(client, csrfToken, newUserLoginToken)

			// Define Variables
			var (
				// The user who you want to boost their followers
				userId string = "User ID"

				// Follow the above user (all variables below)
				bearerToken           string = GetBearerToken(client)
				followUserResponse, _        = FollowUser(client, bearerToken, userId)
			)

			// Make the new account look more legit
			// go FollowRandomArtists(client, bearerToken, (rand.Intn(30-1) + 1))
			// go UpdateProfileImage(client, newUserId, bearerToken)

			// Handle the follow response
			HandleResponse(followUserResponse, iteration, name)
		}(i)
	}
	waitGroup.Wait()
}
