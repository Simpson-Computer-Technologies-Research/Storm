package global

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
)

// Function to get a random profile picture
func GetRandomImage(client *http.Client) io.Reader {
	var req, _ = http.NewRequest("GET", "https://picsum.photos/500/500", nil)
	req.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	var resp, _ = client.Do(req)
	return resp.Body
}

// Function to read a json file
func ReadJsonFile(fileName string) map[string]interface{} {
	var (
		result       map[string]interface{}
		jsonFile, _  = os.Open(fileName)
		byteValue, _ = ioutil.ReadAll(jsonFile)
	)
	defer jsonFile.Close()

	// Marshal to map and return it
	json.Unmarshal([]byte(byteValue), &result)
	return result
}

// Function to generate a random string
func RandomString(length int) string {
	var (
		letterRunes []rune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		res         []rune = make([]rune, length)
	)
	for i := range res {
		res[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(res)
}

// Function to check if a slice contains a string
func SliceContains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// Function to generate random user info
func GenerateFakeInfo(client *http.Client) (string, string) {
	var req, _ = http.NewRequest("GET", "https://api.namefake.com/", nil)
	req.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"},
	}
	var (
		data    map[string]interface{}
		resp, _ = client.Do(req)
	)
	json.NewDecoder(resp.Body).Decode(&data)
	var (
		name  = data["name"].(string)
		email = fmt.Sprintf("%v%s@%v", data["username"], RandomString(7), data["email_d"])
	)
	return name, email
}
