package chrome

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func isCondoUrl(url string) bool {
	if strings.HasPrefix(url, "https://condo.singaporeexpats.com/condo") {
		return true
	}

	return false
}

func getDebugURL() string {
	resp, err := http.Get("http://localhost:9222/json/version")
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}
	return result["webSocketDebuggerUrl"].(string)
}
