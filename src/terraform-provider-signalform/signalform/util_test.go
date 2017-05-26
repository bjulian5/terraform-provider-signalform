package signalform

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSendRequestSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `Test Response`)
	}))
	defer server.Close()

	status_code, body, err := sendRequest("GET", server.URL, "token", nil)
	assert.Equal(t, 200, status_code)
	assert.Equal(t, "Test Response\n", string(body))
	assert.Nil(t, err)
}

func TestSendRequestResponseNotFound(t *testing.T) {
	// Handler returns 404 page not found
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	status_code, body, err := sendRequest("POST", server.URL, "token", nil)
	assert.Equal(t, 404, status_code)
	assert.Contains(t, string(body), "page not found")
	assert.Nil(t, err)
}

func TestSendRequestFail(t *testing.T) {
	// Client will fail to send due to invalid URL
	status_code, body, err := sendRequest("GET", "", "token", nil)
	assert.Equal(t, -1, status_code)
	assert.Nil(t, body)
	assert.Contains(t, err.Error(), "Failed sending GET request")
}

func TestSendRequestTimeout(t *testing.T) {
	timeout := time.Duration(1 * time.Second)
	server := httptest.NewServer(http.TimeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(timeout)
	}), timeout, "Timeout occurred"))
	defer server.Close()

	status_code, body, err := sendRequest("POST", server.URL, "token", nil)
	assert.Equal(t, 503, status_code)
	assert.Equal(t, "Timeout occurred", string(body))
	assert.Nil(t, err)
}

func TestValidateSignalfxRelativeTimeMinutes(t *testing.T) {
	_, errors := validateSignalfxRelativeTime("-5m", "time_range")
	assert.Equal(t, 0, len(errors))
}

func TestValidateSignalfxRelativeTimeHours(t *testing.T) {
	_, errors := validateSignalfxRelativeTime("-5h", "time_range")
	assert.Equal(t, 0, len(errors))
}

func TestValidateSignalfxRelativeTimeDays(t *testing.T) {
	_, errors := validateSignalfxRelativeTime("-5d", "time_range")
	assert.Equal(t, 0, len(errors))
}

func TestValidateSignalfxRelativeTimeWeeks(t *testing.T) {
	_, errors := validateSignalfxRelativeTime("-5w", "time_range")
	assert.Equal(t, 0, len(errors))
}

func TestValidateSignalfxRelativeTimeNotAllowed(t *testing.T) {
	_, errors := validateSignalfxRelativeTime("-5M", "time_range")
	assert.Equal(t, 1, len(errors))
}

func TestConversionSignalfxrealtiveTimeIntoMs(t *testing.T) {
	ms, err := fromRangeToMilliSeconds("-15m")
	assert.Equal(t, 900000, ms)
	assert.Nil(t, err)
}

func TestValidateSortByAscending(t *testing.T) {
	_, errors := validateSortBy("+foo", "sort_by")
	assert.Equal(t, 0, len(errors))
}

func TestValidateSortByDescending(t *testing.T) {
	_, errors := validateSortBy("-foo", "sort_by")
	assert.Equal(t, 0, len(errors))
}

func TestValidateSortByNoDirection(t *testing.T) {
	_, errors := validateSortBy("foo", "sort_by")
	assert.Equal(t, 1, len(errors))
}

func TestSanitizeProgramTextSane(t *testing.T) {
	text := "previous = data('statmonster.inbound_lines',filter('source_region','${var.clusters_no_uswest2[count.index]}')).timeshift('2m').sum()\nsignal = data('statmonster.inbo    und_lines',filter('source_region','${var.clusters_no_uswest2[count.index]}')).sum()\ndetect('Low number of log lines', when(signal < (previous * 0.50), '2m', 0.90))"
	assert.Equal(t, text, sanitizeProgramText(text))
}

func TestSanitizeProgramTextBlankLine(t *testing.T) {
	text := "previous = data('statmonster.inbound_lines',filter('source_region','${var.clusters_no_uswest2[count.index]}')).timeshift('2m').sum()\n\nsignal = data('statmonster.inbo    und_lines',filter('source_region','${var.clusters_no_uswest2[count.index]}')).sum()\n\n\ndetect('Low number of log lines', when(signal < (previous * 0.50), '2m', 0.90))"
	sane_text := "previous = data('statmonster.inbound_lines',filter('source_region','${var.clusters_no_uswest2[count.index]}')).timeshift('2m').sum()\nsignal = data('statmonster.inbo    und_lines',filter('source_region','${var.clusters_no_uswest2[count.index]}')).sum()\ndetect('Low number of log lines', when(signal < (previous * 0.50), '2m', 0.90))"
	assert.Equal(t, sane_text, sanitizeProgramText(text))
}
