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

type WeatherbitService struct {
	Config      domain.ConfigProvider
	Logger      domain.LogProvider
	Http        domain.HttpProvider
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

func (s WeatherbitService) GetCurrent(location model.Coordinates) (domain.WeatherResult, error) {

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

// Repository 'Constructor' function
func NewWeatherbitRepository(config domain.ConfigProvider, http domain.HttpProvider, logger domain.LogProvider) WeatherbitService {

	apiKey := config.GetString("WEATHER_BIT_API_KEY")
	baseurl := config.GetString("WEATHER_BIT_API_BASE_URL")

	localConfig := LocalWeatherConfig{
		CurrentWeatherUrl: baseurl + "/current?",
		ApiKey:            apiKey,
	}

	return WeatherbitService{
		Config:      config,
		Logger:      logger,
		Http:        http,
		LocalConfig: localConfig,
	}
}

//
// Helpers
//

func (s WeatherbitService) SummarizeResponse(res CurrentWeatherResponse) (domain.WeatherResult, error) {
	emptyResult := domain.WeatherResult{}
	if len(res.Data) != 1 {
		return emptyResult, errors.New("expected only one response")
	}

	// we are only requesting one item from the api
	data := res.Data[0]

	r, err := determineTimes(data.ObTime, data.Sunrise, data.Sunset)

	if err != nil {
		return emptyResult, err
	}

	result := domain.WeatherResult{
		CloudCover: data.Clouds,
		DaylightTimes: model.DaylightTimes{
			Observation: r.Observation,
			Sunrise:     r.Sunrise,
			Sunset:      r.Sunset,
		},
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
func determineTimes(observationTimeSt, sunrise, sunset string) (model.DaylightTimes, error) {
	observationTime, err := time.Parse("2006-01-02 15:04:05", appendZeroSeconds(observationTimeSt))

	if err != nil {
		return model.DaylightTimes{}, err
	}

	datePrefix := extractDateString(observationTimeSt)

	return model.DaylightTimes{
		Observation: observationTime,
		Sunrise:     getTimeOnDate(datePrefix, sunrise),
		Sunset:      getTimeOnDate(datePrefix, sunset),
	}, nil
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
