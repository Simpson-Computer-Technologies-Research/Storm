package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

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
