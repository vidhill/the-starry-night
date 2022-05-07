package domain

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/vidhill/the-starry-night/model"
)

type ISSRepositoryRest struct {
	config      ConfigRepository
	localConfig LocalConfig
	logger      LoggerRepository
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

func (s ISSRepositoryRest) GetCurrentLocation() (model.Coordinates, error) {
	logger := s.logger
	emptyResult := model.Coordinates{}

	response, err := http.Get(s.localConfig.url)

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

func NewISSRepositoryRest(config ConfigRepository, logger LoggerRepository) ISSRepositoryRest {
	localConfig := LocalConfig{
		url: config.GetString("ISS_API_URL"),
	}
	return ISSRepositoryRest{
		config:      config,
		localConfig: localConfig,
		logger:      logger,
	}
}

func (s ISSRepositoryRest) SummarizeResponse(a ApiResponse) (model.Coordinates, error) {
	logger := s.logger
	position := a.IssPosition

	latitude, err := parseFloat(position.Latitude)
	if err != nil {
		logger.Error("Error parsing float string", err.Error())
		return model.Coordinates{}, err
	}

	longitude, err := parseFloat(position.Longitude)
	if err != nil {
		logger.Error("Error parsing float string", err.Error())
		return model.Coordinates{}, err
	}

	result := model.Coordinates{
		Latitude:  latitude,
		Longitude: longitude,
	}

	return result, nil
}

func parseFloat(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}
