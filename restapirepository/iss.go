package restapirepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	customErrors "github.com/vidhill/the-starry-night/customerrors"
	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
	"github.com/vidhill/the-starry-night/utils"
)

type ISSLocationRepositoryRest struct {
	Config      domain.ConfigRepository
	Http        domain.HttpRepository
	LocalConfig LocalConfig
	Logger      domain.LoggerRepository
}

type LocalConfig struct {
	url string
}

type ApiResponse struct {
	IssPosition struct {
		Longitude string `json:"longitude"`
		Latitude  string `json:"latitude"`
	} `json:"iss_position"`
	Message string `json:"message"`
}

func (s ISSLocationRepositoryRest) GetCurrentLocation() (model.Coordinates, error) {
	logger := s.Logger
	emptyResult := model.Coordinates{}

	response, err := s.Http.Get(s.LocalConfig.url)

	if err != nil {
		logger.Error("Error fetching", err.Error())
		return emptyResult, err
	}

	if response.StatusCode != http.StatusOK {
		errMessage := fmt.Sprintf("non success response code received from api, received %v", response.StatusCode)
		logger.Error(errMessage)
		return emptyResult, customErrors.NewUnexpectedResponseError(response.StatusCode)
	}

	result := ApiResponse{}
	decodeErr := json.NewDecoder(response.Body).Decode(&result)

	if decodeErr != nil {
		logger.Error("Error decoding", decodeErr.Error())
		return emptyResult, err
	}

	if result.Message != "success" {
		return emptyResult, errors.New("non success response received from api")
	}

	return s.SummarizeResponse(result)
}

//
// Repository 'Constructor' function
//
func NewISSRepositoryRest(
	config domain.ConfigRepository,
	http domain.HttpRepository,
	logger domain.LoggerRepository,
) ISSLocationRepositoryRest {

	localConfig := LocalConfig{
		url: config.GetString("ISS_API_URL"),
	}
	return ISSLocationRepositoryRest{
		Config:      config,
		Http:        http,
		LocalConfig: localConfig,
		Logger:      logger,
	}
}

//
// Helper functions
//

func (s ISSLocationRepositoryRest) SummarizeResponse(a ApiResponse) (model.Coordinates, error) {
	logger := s.Logger
	position := a.IssPosition

	coordinates, err := utils.MakeCoordinatesFromString(position.Latitude, position.Longitude)

	if err != nil {
		logger.Error("Error parsing float string", err.Error())
		return model.Coordinates{}, err
	}

	return coordinates, nil
}

type UnMarshalError struct {
	msg string
	mid error
}

func (e UnMarshalError) Error() string {
	return fmt.Sprintf("%s, %v", e.msg, e.mid)
}

func (e UnMarshalError) Unwrap() error {
	return e.mid
}

func NewUnMarshalError(msg string, err error) UnMarshalError {
	return UnMarshalError{
		msg: msg,
		mid: err,
	}
}
