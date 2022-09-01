package global

// Import Packages
import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
)

// The GetRandomImage() function is used to get a random
// image from the picsum.photos api. This random image
// is used in making the fake account look more real.
func GetRandomImage(client *http.Client) io.Reader {
	// Create a new request object
	var req, _ = http.NewRequest("GET", "https://picsum.photos/500/500", nil)
	// Set the request headers
	req.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	// Send the http request and return it's
	// response body
	var resp, _ = client.Do(req)
	return resp.Body
}

// The ReadJsonFile() function is used to read
// the json data in the provided file.
func ReadJsonFile(fileName string) map[string]interface{} {
	// Define Variables
	var (
		// The readable golang map
		result map[string]interface{}
		// Read the json file
		jsonFile, _  = os.Open(fileName)
		byteValue, _ = ioutil.ReadAll(jsonFile)
	)
	// Close the jsonFile once the function returns
	defer jsonFile.Close()

	// Marshal the json data to the result
	// map, then return said map
	json.Unmarshal([]byte(byteValue), &result)
	return result
}

// The RandomString() function is used to generate
// a random string with the provided length.
func RandomString(length int) string {
	// Define Variables
	var (
		// Characters to use in the random string
		chars string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		// The result byte sloice
		res []byte = make([]byte, length)
	)
	// For each index in the result byte
	for i := range res {
		// Set it to a random character
		res[i] = chars[rand.Intn(len(chars))]
	}
	return string(res)
}

// The SliceContains() function is useds to return
// whether the provided str parameter exists in the
// provided slice paramater
func SliceContains(slice []string, str string) bool {
	// For each value in slice
	for _, v := range slice {
		// If the value equals the provided string
		if v == str {
			return true
		}
	}
	return false
}

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
