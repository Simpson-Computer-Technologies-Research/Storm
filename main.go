package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	Global "spotify_follower_booster/global"
	Spotify "spotify_follower_booster/spotify"
	"sync"
	"time"
)

// Function to handle New Follower response
func HandleResponse(resp *http.Response, iteration *int, name *string) {
	if resp.StatusCode == 204 {
		fmt.Printf(" >> Added New Follower [%d] [%s]\n", *iteration+1, *name)
	} else {
		var bodyBytes, _ = io.ReadAll(resp.Body)
		fmt.Println(resp.StatusCode, ":", string(bodyBytes))
	}
}

// Function to get amount of followers to add
func GetAmount(text string) int {
	var i int
	fmt.Print(text)
	fmt.Scanf("%d", &i)
	return i
}

// Initialize random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Get amount of follows and create wait group
	var (
		amount    int            = GetAmount(" >> Amount: ")
		waitGroup sync.WaitGroup = sync.WaitGroup{}
	)
	waitGroup.Add(amount)

	// Start Goroutines
	for i := 0; i < amount; i++ {
		go func(iteration int) {
			// Create a new account
			var (
				jar, _                            = cookiejar.New(nil)
				client               *http.Client = &http.Client{Jar: jar}
				name, email                       = Global.GenerateFakeInfo(client)
				_, newUserLoginToken              = Spotify.CreateNewAccount(client, name, email, "secretpassword!") // Replace _ with newUserId if using line 70
				csrfToken            *string      = Spotify.GetCSRFToken(client)
			)

			// Authenticate new account
			Spotify.AuthenticateAccount(client, csrfToken, newUserLoginToken)

			// Follow the user
			var (
				userId                string  = "User ID"
				bearerToken           *string = Spotify.GetBearerToken(client)
				followUserResponse, _         = Spotify.FollowUser(client, bearerToken, &userId)
			)

			// Make the new account look more legit
			// go Spotify.FollowRandomArtists(client, bearerToken, (rand.Intn(30-1) + 1))
			// go Spotify.UpdateProfileImage(client, newUserId, bearerToken)

			// Handle the follow response
			go HandleResponse(followUserResponse, &iteration, name)
		}(i)
	}
	waitGroup.Wait()
}
