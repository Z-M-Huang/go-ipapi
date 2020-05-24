package ipapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	countStart    time.Time
	resetDuration time.Duration
	rateLimit     int64
	rateCounter   int64
)

//Response from http://ip-api.com/json/
type Response struct {
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

// Get from http://ip-api.com/json/
// Rate limit: 45 per minute
func Get(host string) (*Response, error) {
	if reachedLimit() {
		return nil, errors.New("Rate limit reached")
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s", host)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	//Per documentation
	if res.StatusCode == http.StatusForbidden ||
		res.StatusCode == http.StatusTooManyRequests {
		rateCounter = rateLimit
		return nil, errors.New("Rate limit reached")
	} else if res.StatusCode != 200 {
		return nil, fmt.Errorf("status: %s", res.Status)
	}

	rateCounter++
	response := &Response{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %s", err.Error())
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %s", err.Error())
	}
	if response.Status != "success" {
		return nil, errors.New(response.Message)
	}
	return response, nil
}

func reachedLimit() bool {
	if countStart.IsZero() {
		countStart = time.Now()
		resetDuration = 1 * time.Minute
		rateLimit = 45
		rateCounter = 0
		return false
	}
	if time.Now().After(countStart.Add(resetDuration)) {
		return rateCounter >= rateLimit
	}
	countStart = time.Now()
	rateCounter = 0
	return false
}
