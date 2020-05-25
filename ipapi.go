//Package ipapi is an unofficial library for http://ip-api.com
//Author: Z-M-Huang
//Repository: https://github.com/Z-M-Huang/go-ipapi
package ipapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var (
	countStart  time.Time
	ttl         time.Duration
	requestLeft int64
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
		requestLeft = 0
		return nil, errors.New("Rate limit reached")
	} else if res.StatusCode != 200 {
		return nil, fmt.Errorf("status: %s", res.Status)
	}

	countStart = time.Now()
	requestLeft, _ = strconv.ParseInt(res.Header.Get("X-Rl"), 10, 64)
	sec, _ := strconv.ParseInt(res.Header.Get("X-Ttl"), 10, 64)
	ttl = time.Duration(sec) * time.Second
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
		return false
	}
	if time.Now().Before(countStart.Add(ttl)) {
		return requestLeft <= 0
	}
	return false
}
