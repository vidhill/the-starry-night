package restapirepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
)

var (
	// regex to match date portion from string formatted "2017-08-28 16:45"
	dateRegex = regexp.MustCompile("^[0-9-]+")
)

type WeatherbitRepository struct {
	Config      domain.ConfigRepository
	Logger      domain.LoggerRepository
	Http        domain.HttpRepository
	LocalConfig LocalWeatherConfig
}

type LocalWeatherConfig struct {
	CurrentWeatherUrl string
	ApiKey            string
}

// API response JSON struct
type CurrentWeatherResponse struct {
	Data  []InterestedData `json:"data"`
	Count int              `json:"count"`
}

type InterestedData struct {
	Clouds  int    `json:"clouds"`
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
	ObTime  string `json:"ob_time"`
}

func (s WeatherbitRepository) GetCurrent(location model.Coordinates) (domain.WeatherResult, error) {

	localConfig := s.LocalConfig
	logger := s.Logger

	emptyResult := domain.WeatherResult{}

	url := localConfig.CurrentWeatherUrl + makeQueryParams(location, localConfig.ApiKey)

	response, err := s.Http.Get(url)

	if err != nil {
		return emptyResult, err
	}

	if response.StatusCode != http.StatusOK {
		logger.Error("non success response from api, code returned: ", response.StatusCode)
		switch response.StatusCode {
		case http.StatusForbidden:
			return emptyResult, errors.New("api key may be expired")
		case http.StatusBadRequest:
			return emptyResult, errors.New("api invalid request")
		}

		return emptyResult, fmt.Errorf("non success response from api, %v", response.StatusCode)
	}

	result := CurrentWeatherResponse{}
	decodeErr := json.NewDecoder(response.Body).Decode(&result)

	if decodeErr != nil {
		logger.Error("Error decoding", decodeErr.Error())
		return emptyResult, err
	}

	return s.SummarizeResponse(result)

}

//
// Repository 'Constructor' function
//
func NewWeatherbitRepository(config domain.ConfigRepository, http domain.HttpRepository, logger domain.LoggerRepository) WeatherbitRepository {

	apiKey := config.GetString("WEATHER_BIT_API_KEY")
	baseurl := config.GetString("WEATHER_BIT_API_BASE_URL")

	localConfig := LocalWeatherConfig{
		CurrentWeatherUrl: baseurl + "/current?",
		ApiKey:            apiKey,
	}

	return WeatherbitRepository{
		Config:      config,
		Logger:      logger,
		Http:        http,
		LocalConfig: localConfig,
	}
}

//
// Helpers
//

func (s WeatherbitRepository) SummarizeResponse(res CurrentWeatherResponse) (domain.WeatherResult, error) {
	emptyResult := domain.WeatherResult{}
	if len(res.Data) != 1 {
		return emptyResult, errors.New("expected only one response")
	}

	// we are only requesting one item from the api
	data := res.Data[0]

	observerationTime, sunriseTime, sunsetTime, err := determineTimes(data.ObTime, data.Sunrise, data.Sunset)

	if err != nil {
		return emptyResult, err
	}

	result := domain.WeatherResult{
		CloudCover:        data.Clouds,
		ObserverationTime: observerationTime,
		Sunrise:           sunriseTime,
		Sunset:            sunsetTime,
	}
	return result, nil
}

func makeQueryParams(location model.Coordinates, apiKey string) string {

	v := url.Values{}

	v.Add("lat", fmt.Sprintf("%f", location.Latitude))
	v.Add("lon", fmt.Sprintf("%f", location.Longitude))
	v.Add("key", apiKey)

	return v.Encode()
}

// parse the time strings into go Time.time objects
func determineTimes(observerationTimeSt, sunrise, sunset string) (time.Time, time.Time, time.Time, error) {
	observerationTime, err := time.Parse("2006-01-02 15:04:05", appendZeroSeconds(observerationTimeSt))

	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, err
	}

	datePrefix := extractDateString(observerationTimeSt)

	sunriseTime := getTimeOnDate(datePrefix, sunrise)
	sunsetTime := getTimeOnDate(datePrefix, sunset)

	return observerationTime, sunriseTime, sunsetTime, nil
}

func extractDateString(date string) string {
	dateMatch := dateRegex.FindStringSubmatch(date)

	if len(dateMatch) == 1 {
		return dateMatch[0]
	}
	return ""

}

// get Time.time object for the date at the passed time string
func getTimeOnDate(date, timeStr string) time.Time {
	dateWithTime := fmt.Sprintf(`%s %s`, date, appendZeroSeconds(timeStr))

	t, _ := time.Parse("2006-01-02 15:04:05", dateWithTime)
	return t
}

func appendZeroSeconds(timeS string) string {
	return timeS + ":00"
}
