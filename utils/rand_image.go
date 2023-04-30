package utils

import (
	"io"
	"net/http"
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
