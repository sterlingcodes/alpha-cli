package weather

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewCmd(t *testing.T) {
	cmd := NewCmd()
	if cmd.Use != "weather" {
		t.Errorf("expected Use 'weather', got %q", cmd.Use)
	}

	aliases := cmd.Aliases
	if len(aliases) != 1 || aliases[0] != "wttr" {
		t.Errorf("expected Aliases ['wttr'], got %v", aliases)
	}

	subs := map[string]bool{}
	for _, s := range cmd.Commands() {
		subs[s.Use] = true
	}
	for _, name := range []string{"now [location]", "forecast [location]"} {
		if !subs[name] {
			t.Errorf("missing subcommand %q", name)
		}
	}
}

func TestWmoDescription(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{0, "Clear sky"},
		{3, "Overcast"},
		{61, "Slight rain"},
		{95, "Thunderstorm"},
		{999, "Unknown"},
	}
	for _, tt := range tests {
		got := wmoDescription(tt.code)
		if got != tt.want {
			t.Errorf("wmoDescription(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestDegreesToCardinal(t *testing.T) {
	tests := []struct {
		deg  float64
		want string
	}{
		{0, "N"},
		{90, "E"},
		{180, "S"},
		{270, "W"},
		{45, "NE"},
		{225, "SW"},
	}
	for _, tt := range tests {
		got := degreesToCardinal(tt.deg)
		if got != tt.want {
			t.Errorf("degreesToCardinal(%.0f) = %q, want %q", tt.deg, got, tt.want)
		}
	}
}

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		c    float64
		want int
	}{
		{0, 32},
		{100, 212},
		{20, 68},
		{-40, -40},
	}
	for _, tt := range tests {
		got := celsiusToFahrenheit(tt.c)
		if got != tt.want {
			t.Errorf("celsiusToFahrenheit(%.1f) = %d, want %d", tt.c, got, tt.want)
		}
	}
}

func TestKphToMph(t *testing.T) {
	got := kphToMph(100)
	want := int(math.Round(100 / 1.609))
	if got != want {
		t.Errorf("kphToMph(100) = %d, want %d", got, want)
	}
}

// mockServer creates a test server that handles both geocode and forecast endpoints
func mockServer(geoResponse, forecastResponse any, geoStatus, fcStatus int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/v1/search") {
			w.WriteHeader(geoStatus)
			if geoResponse != nil {
				json.NewEncoder(w).Encode(geoResponse)
			}
			return
		}
		if strings.HasPrefix(r.URL.Path, "/v1/forecast") {
			w.WriteHeader(fcStatus)
			if forecastResponse != nil {
				json.NewEncoder(w).Encode(forecastResponse)
			}
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func sampleGeoResponse() map[string]any {
	return map[string]any{
		"results": []map[string]any{
			{
				"name":      "San Francisco",
				"country":   "United States",
				"latitude":  37.7749,
				"longitude": -122.4194,
			},
		},
	}
}

func sampleForecastResponse() map[string]any {
	return map[string]any{
		"current": map[string]any{
			"temperature_2m":        20.0,
			"apparent_temperature":  19.0,
			"weather_code":          0,
			"relative_humidity_2m":  50,
			"wind_speed_10m":        10.0,
			"wind_direction_10m":    0.0,
			"precipitation":         0.0,
			"visibility":            10000.0,
			"uv_index":              3.0,
		},
		"daily": map[string]any{
			"time":                         []string{"2026-03-08"},
			"weather_code":                 []int{0},
			"temperature_2m_max":           []float64{22.0},
			"temperature_2m_min":           []float64{15.0},
			"precipitation_probability_max": []int{20},
			"relative_humidity_2m_max":     []int{60},
		},
	}
}

func sampleForecastWithDays() map[string]any {
	return map[string]any{
		"current": map[string]any{
			"temperature_2m":        20.0,
			"apparent_temperature":  19.0,
			"weather_code":          2,
			"relative_humidity_2m":  50,
			"wind_speed_10m":        10.0,
			"wind_direction_10m":    180.0,
			"precipitation":         0.0,
			"visibility":            10000.0,
			"uv_index":              3.0,
		},
		"daily": map[string]any{
			"time":                         []string{"2026-03-08", "2026-03-09"},
			"weather_code":                 []int{2, 61},
			"temperature_2m_max":           []float64{22.0, 18.0},
			"temperature_2m_min":           []float64{15.0, 12.0},
			"precipitation_probability_max": []int{20, 70},
			"relative_humidity_2m_max":     []int{60, 80},
		},
	}
}

func swapURLs(srv *httptest.Server) func() {
	oldGeo := geocodeURL
	oldFc := forecastURL
	geocodeURL = srv.URL
	forecastURL = srv.URL
	return func() {
		geocodeURL = oldGeo
		forecastURL = oldFc
	}
}

func TestFetchWeatherNow(t *testing.T) {
	srv := mockServer(sampleGeoResponse(), sampleForecastResponse(), http.StatusOK, http.StatusOK)
	defer srv.Close()
	defer swapURLs(srv)()

	err := fetchWeather("San Francisco", 0)
	if err != nil {
		t.Errorf("fetchWeather failed: %v", err)
	}
}

func TestFetchWeatherForecast(t *testing.T) {
	srv := mockServer(sampleGeoResponse(), sampleForecastWithDays(), http.StatusOK, http.StatusOK)
	defer srv.Close()
	defer swapURLs(srv)()

	err := fetchWeather("London", 1)
	if err != nil {
		t.Errorf("fetchWeather with forecast failed: %v", err)
	}
}

func TestFetchWeatherNotFound(t *testing.T) {
	geoResp := map[string]any{"results": []map[string]any{}}
	srv := mockServer(geoResp, nil, http.StatusOK, http.StatusOK)
	defer srv.Close()
	defer swapURLs(srv)()

	err := fetchWeather("xyznonexistent", 0)
	if err == nil {
		t.Error("expected error for not found location, got nil")
	}
}

func TestFetchWeatherHTTPError(t *testing.T) {
	srv := mockServer(nil, nil, http.StatusInternalServerError, http.StatusOK)
	defer srv.Close()
	defer swapURLs(srv)()

	err := fetchWeather("test", 0)
	if err == nil {
		t.Error("expected error for HTTP 500, got nil")
	}
}

func TestFetchWeatherParseError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer srv.Close()
	defer swapURLs(srv)()

	err := fetchWeather("test", 0)
	if err == nil {
		t.Error("expected parse error, got nil")
	}
}
