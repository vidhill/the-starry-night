package restapirepository

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
	"github.com/vidhill/the-starry-night/utils"
)

type ISSLocationRepositoryRest struct {
	config      domain.ConfigRepository
	http        domain.HttpRepository
	localConfig LocalConfig
	logger      domain.LoggerRepository
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
	logger := s.logger
	emptyResult := model.Coordinates{}

	response, err := s.http.Get(s.localConfig.url)

	if err != nil {
		logger.Error("Error fetching", err.Error())
		return emptyResult, err
	}

	if response.StatusCode != http.StatusOK {
		errMessage := "non success response code received from api"
		logger.Error(errMessage)
		return emptyResult, errors.New(errMessage)
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
		config:      config,
		http:        http,
		localConfig: localConfig,
		logger:      logger,
	}
}

//
// Helper functions
//

func (s ISSLocationRepositoryRest) SummarizeResponse(a ApiResponse) (model.Coordinates, error) {
	logger := s.logger
	position := a.IssPosition

	coordinates, err := utils.MakeCoordinatesFromString(position.Latitude, position.Longitude)

	if err != nil {
		logger.Error("Error parsing float string", err.Error())
		return model.Coordinates{}, err
	}

	return coordinates, nil
}
