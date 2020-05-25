package ipapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	ipsToTest := []string{"8.8.8.8", "2001:4860:4860:0:0:0:0:6464", "google.com"}
	for _, ip := range ipsToTest {
		resp, err := Get(ip)
		assert.Empty(t, err)
		assert.NotEmpty(t, resp)
		assert.NotEmpty(t, resp.Query)
	}
}

func TestInvalidIPAddress(t *testing.T) {
	ipsToTest := []string{"9999.9999.9999.9999", "2001:4860:4860:a:6464"}
	for _, ip := range ipsToTest {
		resp, err := Get(ip)
		assert.Empty(t, resp)
		assert.NotEmpty(t, err)
		assert.Equal(t, "invalid query", err.Error())
	}
}

func TestRateLimit(t *testing.T) {
	requestLeft = 0
	countStart = time.Now()
	ttl = 1 * time.Minute
	resp, err := Get("")
	assert.Empty(t, resp)
	assert.Equal(t, "Rate limit reached", err.Error())
}
