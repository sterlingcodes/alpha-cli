package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"

	"github.com/sterlingcodes/alpha-cli/pkg/output"
)

var (
	httpClient  = &http.Client{Timeout: 30 * time.Second}
	geocodeURL  = "https://geocoding-api.open-meteo.com"
	forecastURL = "https://api.open-meteo.com"
)

// Current is LLM-friendly current weather
type Current struct {
	Location   string `json:"location"`
	Condition  string `json:"condition"`
	TempC      int    `json:"temp_c"`
	TempF      int    `json:"temp_f"`
	FeelsLikeC int    `json:"feels_like_c"`
	FeelsLikeF int    `json:"feels_like_f"`
	Humidity   int    `json:"humidity"`
	WindKph    int    `json:"wind_kph"`
	WindMph    int    `json:"wind_mph"`
	WindDir    string `json:"wind_dir"`
	Visibility int    `json:"visibility_km"`
	UV         int    `json:"uv"`
}

// Forecast is LLM-friendly forecast day
type Forecast struct {
	Date       string `json:"date"`
	Condition  string `json:"condition"`
	MaxC       int    `json:"max_c"`
	MinC       int    `json:"min_c"`
	MaxF       int    `json:"max_f"`
	MinF       int    `json:"min_f"`
	ChanceRain int    `json:"chance_rain"`
	Humidity   int    `json:"humidity"`
}

// Weather is the full response
type Weather struct {
	Current  Current    `json:"current"`
	Forecast []Forecast `json:"forecast,omitempty"`
}

// WMO weather code descriptions
var wmoDescriptions = map[int]string{
	0:  "Clear sky",
	1:  "Mainly clear",
	2:  "Partly cloudy",
	3:  "Overcast",
	45: "Fog",
	48: "Rime fog",
	51: "Light drizzle",
	53: "Moderate drizzle",
	55: "Dense drizzle",
	56: "Light freezing drizzle",
	57: "Dense freezing drizzle",
	61: "Slight rain",
	63: "Moderate rain",
	65: "Heavy rain",
	66: "Light freezing rain",
	67: "Heavy freezing rain",
	71: "Slight snow",
	73: "Moderate snow",
	75: "Heavy snow",
	77: "Snow grains",
	80: "Slight rain showers",
	81: "Moderate rain showers",
	82: "Violent rain showers",
	85: "Slight snow showers",
	86: "Heavy snow showers",
	95: "Thunderstorm",
	96: "Thunderstorm with slight hail",
	99: "Thunderstorm with heavy hail",
}

func wmoDescription(code int) string {
	if desc, ok := wmoDescriptions[code]; ok {
		return desc
	}
	return "Unknown"
}

func degreesToCardinal(deg float64) string {
	directions := []string{
		"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE",
		"S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW",
	}
	idx := int(math.Round(deg/22.5)) % 16
	return directions[idx]
}

func celsiusToFahrenheit(c float64) int {
	return int(math.Round(c*9.0/5.0 + 32))
}

func kphToMph(kph float64) int {
	return int(math.Round(kph / 1.609))
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "weather",
		Aliases: []string{"wttr"},
		Short:   "Weather commands",
	}

	cmd.AddCommand(newNowCmd())
	cmd.AddCommand(newForecastCmd())

	return cmd
}

func newNowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "now [location]",
		Short: "Get current weather",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			location := args[0]
			return fetchWeather(location, 0)
		},
	}

	return cmd
}

func newForecastCmd() *cobra.Command {
	var days int

	cmd := &cobra.Command{
		Use:   "forecast [location]",
		Short: "Get weather forecast",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			location := args[0]
			return fetchWeather(location, days)
		},
	}

	cmd.Flags().IntVarP(&days, "days", "d", 3, "Number of days (1-3)")

	return cmd
}

type geocodeResult struct {
	Results []struct {
		Name      string  `json:"name"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type forecastResult struct {
	Current struct {
		Temperature      float64 `json:"temperature_2m"`
		ApparentTemp     float64 `json:"apparent_temperature"`
		WeatherCode      int     `json:"weather_code"`
		Humidity         int     `json:"relative_humidity_2m"`
		WindSpeed        float64 `json:"wind_speed_10m"`
		WindDirection    float64 `json:"wind_direction_10m"`
		Precipitation    float64 `json:"precipitation"`
		Visibility       float64 `json:"visibility"`
		UVIndex          float64 `json:"uv_index"`
	} `json:"current"`
	Daily struct {
		Time              []string  `json:"time"`
		WeatherCode       []int     `json:"weather_code"`
		TempMax           []float64 `json:"temperature_2m_max"`
		TempMin           []float64 `json:"temperature_2m_min"`
		PrecipProbMax     []int     `json:"precipitation_probability_max"`
		HumidityMax       []int     `json:"relative_humidity_2m_max"`
	} `json:"daily"`
}

//nolint:gocyclo // complex but clear sequential logic
func fetchWeather(location string, forecastDays int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Step 1: Geocode the location
	geoURL := fmt.Sprintf("%s/v1/search?name=%s&count=1", geocodeURL, url.QueryEscape(location))

	geoReq, err := http.NewRequestWithContext(ctx, "GET", geoURL, http.NoBody)
	if err != nil {
		return output.PrintError("fetch_failed", err.Error(), nil)
	}

	geoResp, err := httpClient.Do(geoReq)
	if err != nil {
		return output.PrintError("fetch_failed", err.Error(), nil)
	}
	defer geoResp.Body.Close()

	if geoResp.StatusCode >= 400 {
		return output.PrintError("fetch_failed", fmt.Sprintf("geocode HTTP %d", geoResp.StatusCode), nil)
	}

	var geo geocodeResult
	if err := json.NewDecoder(geoResp.Body).Decode(&geo); err != nil {
		return output.PrintError("parse_failed", err.Error(), nil)
	}

	if len(geo.Results) == 0 {
		return output.PrintError("not_found", "Location not found: "+location, nil)
	}

	place := geo.Results[0]
	loc := place.Name + ", " + place.Country

	// Step 2: Fetch forecast data
	days := forecastDays
	if days < 1 {
		days = 1
	}
	if days > 3 {
		days = 3
	}

	fcURL := fmt.Sprintf(
		"%s/v1/forecast?latitude=%.4f&longitude=%.4f"+
			"&current=temperature_2m,apparent_temperature,weather_code,relative_humidity_2m,wind_speed_10m,wind_direction_10m,precipitation,visibility,uv_index"+
			"&daily=weather_code,temperature_2m_max,temperature_2m_min,precipitation_probability_max,relative_humidity_2m_max"+
			"&temperature_unit=celsius&wind_speed_unit=kmh&forecast_days=%d&timezone=auto",
		forecastURL, place.Latitude, place.Longitude, days+1, // +1 to include today
	)

	fcReq, err := http.NewRequestWithContext(ctx, "GET", fcURL, http.NoBody)
	if err != nil {
		return output.PrintError("fetch_failed", err.Error(), nil)
	}

	fcResp, err := httpClient.Do(fcReq)
	if err != nil {
		return output.PrintError("fetch_failed", err.Error(), nil)
	}
	defer fcResp.Body.Close()

	if fcResp.StatusCode >= 400 {
		return output.PrintError("fetch_failed", fmt.Sprintf("forecast HTTP %d", fcResp.StatusCode), nil)
	}

	var fc forecastResult
	if err := json.NewDecoder(fcResp.Body).Decode(&fc); err != nil {
		return output.PrintError("parse_failed", err.Error(), nil)
	}

	weather := Weather{
		Current: Current{
			Location:   loc,
			Condition:  wmoDescription(fc.Current.WeatherCode),
			TempC:      int(math.Round(fc.Current.Temperature)),
			TempF:      celsiusToFahrenheit(fc.Current.Temperature),
			FeelsLikeC: int(math.Round(fc.Current.ApparentTemp)),
			FeelsLikeF: celsiusToFahrenheit(fc.Current.ApparentTemp),
			Humidity:   fc.Current.Humidity,
			WindKph:    int(math.Round(fc.Current.WindSpeed)),
			WindMph:    kphToMph(fc.Current.WindSpeed),
			WindDir:    degreesToCardinal(fc.Current.WindDirection),
			Visibility: int(math.Round(fc.Current.Visibility / 1000)),
			UV:         int(math.Round(fc.Current.UVIndex)),
		},
	}

	// Add forecast if requested (skip index 0 = today)
	if forecastDays > 0 && len(fc.Daily.Time) > 1 {
		limit := forecastDays
		available := len(fc.Daily.Time) - 1 // exclude today
		if limit > available {
			limit = available
		}

		for i := 1; i <= limit; i++ {
			weather.Forecast = append(weather.Forecast, Forecast{
				Date:       fc.Daily.Time[i],
				Condition:  wmoDescription(fc.Daily.WeatherCode[i]),
				MaxC:       int(math.Round(fc.Daily.TempMax[i])),
				MinC:       int(math.Round(fc.Daily.TempMin[i])),
				MaxF:       celsiusToFahrenheit(fc.Daily.TempMax[i]),
				MinF:       celsiusToFahrenheit(fc.Daily.TempMin[i]),
				ChanceRain: fc.Daily.PrecipProbMax[i],
				Humidity:   fc.Daily.HumidityMax[i],
			})
		}
	}

	return output.Print(weather)
}
